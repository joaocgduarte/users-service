package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
)

func (s DefaultUserService) GetUserByUUID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ContextTimeout)
	defer cancel()

	if id == uuid.Nil {
		return nil, domain.ErrBadParamInput
	}

	user, err := s.UserRepo.GetByUUID(ctx, id)
	if err != nil {
		s.Logger.Println("error fetching user by uuid: " + err.Error())
		return nil, domain.ErrNotFound
	}

	userRole, err := s.RoleRepo.GetByUUID(ctx, user.RoleId)
	if err != nil {
		s.Logger.Println("error fetching role by uuid: " + err.Error())
		return nil, domain.ErrNotFound
	}
	user.Role = &userRole
	return user, nil
}
