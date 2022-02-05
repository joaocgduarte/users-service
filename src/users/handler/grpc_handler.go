package handler

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
	"github.com/plagioriginal/user-microservice/users/tokens"
	users "github.com/plagioriginal/users-service-grpc/users"
)

type UserGRPCHandler struct {
	users.UnimplementedUsersServer
	l            *log.Logger
	tokenManager tokens.TokenManager
	userService  domain.UserService
}

func NewUserGRPCHandler(
	l *log.Logger,
	tokenManager tokens.TokenManager,
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
		return nil, domain.ErrInvalidToken
	}

	token, err := srv.tokenManager.ParseJWT(in.AccessToken)

	if err != nil {
		return nil, err
	}

	if !srv.tokenManager.IsJWTokenValid(token) {
		return nil, domain.ErrInvalidToken
	}

	tokenRole, _ := srv.tokenManager.GetUserRoleFromToken(token)

	if tokenRole != domain.DEFAULT_ROLE_ADMIN.RoleSlug {
		return nil, domain.ErrNotAllowed
	}

	if len(in.Username) == 0 || len(in.Password) == 0 || len(in.Role) == 0 {
		return nil, domain.ErrBadParamInput
	}

	user, err := srv.userService.Store(ctx, in.Username, in.Password, in.Role)

	if err != nil {
		return nil, err
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
		return nil, domain.ErrNotFound
	}

	user, err := srv.userService.GetUserByLogin(ctx, loginRequest.Username, loginRequest.Password)

	if err != nil {
		return nil, err
	}

	// Generates the tokens of said user.
	token, err := srv.tokenManager.GenerateTokens(ctx, user)

	if err != nil {
		return nil, err
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
		return nil, domain.ErrNotFound
	}

	oldRefreshToken, _ := uuid.Parse(in.RefreshToken)
	tokens, err := srv.tokenManager.RefreshAllTokens(ctx, oldRefreshToken)

	if err != nil {
		return nil, err
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
		return nil, domain.ErrNotFound
	}

	isDeleted := srv.tokenManager.DeleteRefreshToken(ctx, in.RefreshToken)

	if !isDeleted {
		srv.l.Println("Error deleting token: " + in.RefreshToken)
	}

	return &users.TokenResponse{
		AccessToken:  "",
		RefreshToken: "",
	}, nil
}
