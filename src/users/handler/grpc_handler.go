package handler

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
	users "github.com/plagioriginal/users-service-grpc/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGRPCHandler struct {
	users.UnimplementedUsersServer
	l            *log.Logger
	tokenManager domain.AccessTokenHandler
	userService  domain.UserService
}

func NewUserGRPCHandler(
	l *log.Logger,
	tokenManager domain.AccessTokenHandler,
	userService domain.UserService,
) users.UsersServer {
	return UserGRPCHandler{
		l:            l,
		tokenManager: tokenManager,
		userService:  userService,
	}
}

// Adds a new user.
func (srv UserGRPCHandler) AddUser(ctx context.Context, in *users.NewUserRequest) (*users.UserResponse, error) {
	if len(in.AccessToken) == 0 {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	token, err := srv.tokenManager.ParseJWT(in.AccessToken)

	if err != nil {
		srv.l.Println("error parsing jwt token in add-user: " + err.Error())
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}

	if !srv.tokenManager.IsJWTokenValid(token) {
		srv.l.Printf("invalid token %v\n", token)
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	tokenRole, _ := srv.tokenManager.GetUserRoleFromToken(token)

	if tokenRole != domain.DEFAULT_ROLE_ADMIN.RoleSlug {
		return nil, status.Error(codes.Unauthenticated, "incorrect permissions")
	}

	if len(in.Username) == 0 || len(in.Password) == 0 || len(in.Role) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	user, err := srv.userService.Store(ctx, domain.StoreUserRequest{
		Username: in.GetUsername(),
		Password: in.GetPassword(),
		RoleSlug: in.GetRole(),
	})

	if err != nil {
		srv.l.Printf("error storing a user: %v\n", err)
		return nil, status.Error(codes.Unknown, "error storing user")
	}

	return &users.UserResponse{
		Id:        user.ID.String(),
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role: &users.UserResponse_RoleResponse{
			Id:        user.RoleId.String(),
			RoleLabel: user.Role.RoleLabel,
			RoleSlug:  user.Role.RoleSlug,
		},
	}, nil
}

// Gets the tokens for the login
func (srv UserGRPCHandler) Login(ctx context.Context, loginRequest *users.LoginRequest) (*users.TokenResponse, error) {
	if len(loginRequest.Username) == 0 || len(loginRequest.Password) == 0 {
		return nil, status.Error(codes.NotFound, "resource found")
	}

	user, err := srv.userService.GetUserByLogin(ctx, domain.GetUserRequest{
		Username: loginRequest.GetUsername(),
		Password: loginRequest.GetPassword(),
	})

	if err != nil {
		srv.l.Printf("error getting the user by login: %v\n", err)
		return nil, status.Error(codes.NotFound, "resource found")
	}

	// Generates the tokens of said user.
	token, err := srv.tokenManager.GenerateTokens(ctx, user)

	if err != nil {
		srv.l.Printf("error generating tokens on login: %v\n", err)
		return nil, status.Error(codes.Unknown, "error generating tokens")
	}

	result := &users.TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User: &users.UserResponse{
			Id:        user.ID.String(),
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role: &users.UserResponse_RoleResponse{
				Id:        user.RoleId.String(),
				RoleLabel: user.Role.RoleLabel,
				RoleSlug:  user.Role.RoleSlug,
			},
		},
	}

	return result, nil
}

// Refreshes the tokens.
func (srv UserGRPCHandler) Refresh(ctx context.Context, in *users.RefreshRequest) (*users.TokenResponse, error) {
	if len(in.RefreshToken) == 0 {
		return nil, status.Error(codes.NotFound, "resource found")
	}

	oldRefreshToken, _ := uuid.Parse(in.RefreshToken)
	tokens, err := srv.tokenManager.RefreshAllTokens(ctx, oldRefreshToken)

	if err != nil {
		srv.l.Printf("error generating tokens on refresh: %v\n", err)
		return nil, status.Error(codes.Unknown, "error generating tokens")
	}

	token, _ := srv.tokenManager.ParseJWT(tokens.AccessToken)
	userId, _ := srv.tokenManager.GetUserIDFromToken(token)
	user, _ := srv.userService.GetUserByUUID(ctx, userId)

	result := &users.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User: &users.UserResponse{
			Id:        user.ID.String(),
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role: &users.UserResponse_RoleResponse{
				Id:        user.RoleId.String(),
				RoleLabel: user.Role.RoleLabel,
				RoleSlug:  user.Role.RoleSlug,
			},
		},
	}

	return result, nil
}

// Logout handler.
// Deletes a refresh token from the database
func (srv UserGRPCHandler) Logout(ctx context.Context, in *users.RefreshRequest) (*users.TokenResponse, error) {
	if len(in.RefreshToken) == 0 {
		return nil, status.Error(codes.NotFound, "resource found")
	}

	isDeleted := srv.tokenManager.DeleteRefreshToken(ctx, in.RefreshToken)

	if !isDeleted {
		srv.l.Println("error deleting token: " + in.RefreshToken)
	}

	return &users.TokenResponse{
		AccessToken:  "",
		RefreshToken: "",
	}, nil
}
