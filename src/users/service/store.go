package service

import (
	"context"

	"github.com/plagioriginal/user-microservice/domain"
	"golang.org/x/crypto/bcrypt"
)

// Stores q new user based on username and password
func (s DefaultUserService) Store(ctx context.Context, request domain.StoreUserRequest) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ContextTimeout)
	defer cancel()

	role, err := s.RoleRepo.GetBySlug(ctx, request.RoleSlug)
	if err != nil {
		return nil, err
	}

	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(request.Password), s.BcryptHashingCost)
	if err != nil {
		return nil, err
	}
	password := string(passwordBytes[:])

	user, err := s.UserRepo.Store(ctx, domain.User{
		Username: request.Username,
		Password: password,
		RoleId:   role.ID,
	})
	if err != nil {
		return nil, err
	}

	user.Role = &role
	return user, nil
}
