package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
	"golang.org/x/crypto/bcrypt"
)

// Gets a user by the username and password. With role attached
func (s DefaultUserService) GetUserByLogin(ctx context.Context, request domain.GetUserRequest) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ContextTimeout)
	defer cancel()

	if len(request.Username) == 0 || len(request.Password) == 0 {
		return nil, domain.ErrBadParamInput
	}

	user, err := s.UserRepo.GetByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}
	if user.ID == uuid.Nil {
		return nil, domain.ErrNotFound
	}

	userRole, err := s.RoleRepo.GetByUUID(ctx, user.RoleId)
	if err != nil {
		return nil, err
	}

	user.Role = &userRole
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))

	if err != nil {
		return nil, err
	}
	return user, nil
}
