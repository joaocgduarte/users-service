package handler

import (
	"context"

	users "github.com/plagioriginal/users-service-grpc/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Logout handler.
// Deletes a refresh token from the database
func (srv UserGRPCHandler) Logout(ctx context.Context, in *users.RefreshRequest) (*users.TokenResponse, error) {
	if in == nil || len(in.RefreshToken) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid input")
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
