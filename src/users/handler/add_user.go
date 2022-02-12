package handler

import (
	"context"

	"github.com/plagioriginal/user-microservice/domain"
	users "github.com/plagioriginal/users-service-grpc/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Adds a new user.
func (srv UserGRPCHandler) AddUser(ctx context.Context, in *users.NewUserRequest) (*users.UserResponse, error) {
	if in == nil || len(in.AccessToken) == 0 {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	token, err := srv.tokenManager.ParseJWT(in.AccessToken)
	if err != nil {
		srv.l.Println("error parsing jwt token in add-user: " + err.Error())
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	if !srv.tokenManager.IsJWTokenValid(token) {
		srv.l.Printf("invalid token %v\n", token)
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	tokenRole, err := srv.tokenManager.GetUserRoleFromToken(token)
	if err != nil {
		srv.l.Printf("error getting role from token: %v\n", err)
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}

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
		return nil, status.Error(codes.Internal, "error storing user")
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
