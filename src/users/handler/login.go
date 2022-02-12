package handler

import (
	"context"

	"github.com/plagioriginal/user-microservice/domain"
	users "github.com/plagioriginal/users-service-grpc/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Gets the tokens for the login
func (srv UserGRPCHandler) Login(ctx context.Context, in *users.LoginRequest) (*users.TokenResponse, error) {
	if in == nil || len(in.Username) == 0 || len(in.Password) == 0 {
		return nil, status.Error(codes.NotFound, "resource found")
	}

	user, err := srv.userService.GetUserByLogin(ctx, domain.GetUserRequest{
		Username: in.GetUsername(),
		Password: in.GetPassword(),
	})

	if err != nil {
		srv.l.Printf("error getting the user by login: %v\n", err)
		return nil, status.Error(codes.NotFound, "resource found")
	}

	// Generates the tokens of said user.
	token, err := srv.tokenManager.GenerateTokens(ctx, user)
	if err != nil {
		srv.l.Printf("error generating tokens on login: %v\n", err)
		return nil, status.Error(codes.Internal, "error generating tokens")
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
