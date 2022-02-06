package service

import (
	"context"
	"errors"

	"github.com/plagioriginal/user-microservice/domain"
	"golang.org/x/crypto/bcrypt"
)

// Stores q new user based on username and password
func (s DefaultUserService) Store(ctx context.Context, request domain.StoreUserRequest) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ContextTimeout)
	defer cancel()

	role, err := s.RoleRepo.GetBySlug(ctx, request.RoleSlug)
	if err != nil {
		s.Logger.Println("error fetching role: " + err.Error())
		return nil, errors.New("error fetching role")
	}

	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
	if err != nil {
		s.Logger.Println("error encrypting password: " + err.Error())
		return nil, errors.New("error encrypting password")
	}
	password := string(passwordBytes[:])

	userToAdd := domain.User{
		Username: request.Username,
		Password: password,
		RoleId:   role.ID,
	}

	user, err := s.UserRepo.Store(ctx, userToAdd)

	if err != nil {
		s.Logger.Println("error storing user: " + err.Error())
		return nil, errors.New("error storing user")
	}

	user.Role = &role
	return user, nil
}
