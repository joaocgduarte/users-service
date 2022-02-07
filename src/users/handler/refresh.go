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
