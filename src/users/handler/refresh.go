package handler

import (
	"context"

	"github.com/google/uuid"
	users "github.com/plagioriginal/users-service-grpc/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Refreshes the tokens.
func (srv UserGRPCHandler) Refresh(ctx context.Context, in *users.RefreshRequest) (*users.TokenResponse, error) {
	if in == nil || len(in.RefreshToken) == 0 {
		return nil, status.Error(codes.NotFound, "resource found")
	}

	oldRefreshToken, err := uuid.Parse(in.RefreshToken)
	if err != nil {
		srv.l.Printf("error parsing token on refresh handler: %v\n", err)
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}

	tokens, err := srv.tokenManager.RefreshAllTokens(ctx, oldRefreshToken)
	if err != nil {
		srv.l.Printf("error generating tokens on refresh: %v\n", err)
		return nil, status.Error(codes.Internal, "error generating tokens")
	}

	token, err := srv.tokenManager.ParseJWT(tokens.AccessToken)
	if err != nil {
		srv.l.Printf("error parsing jwt token on refresh handler: %v\n", err)
		return nil, status.Error(codes.Internal, "error generating tokens")
	}

	userId, err := srv.tokenManager.GetUserIDFromToken(token)
	if err != nil {
		srv.l.Printf("error getting user id on refresh handler: %v\n", err)
		return nil, status.Error(codes.Internal, "error generating tokens")
	}

	user, err := srv.userService.GetUserByUUID(ctx, userId)
	if err != nil {
		srv.l.Printf("error getting user by id on refresh handler: %v\n", err)
		return nil, status.Error(codes.Internal, "error generating tokens")
	}

	return &users.TokenResponse{
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
	}, nil
}
