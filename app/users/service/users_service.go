package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
	"github.com/plagioriginal/user-microservice/users/tokens"
	"golang.org/x/crypto/bcrypt"
)

type DefaultUserService struct {
	UserRepo       domain.UserRepository
	RoleRepo       domain.RoleRepository
	TokenManager   tokens.TokenManager
	ContextTimeout time.Duration
}

// Constructor
func New(
	userRepo domain.UserRepository,
	roleRepo domain.RoleRepository,
	tokenManager tokens.TokenManager,
	contextTimeout time.Duration,
) domain.UserService {
	return DefaultUserService{
		userRepo,
		roleRepo,
		tokenManager,
		contextTimeout,
	}
}

// Stores q new user based on username and password
func (s DefaultUserService) Store(ctx context.Context, username string, password string, roleSlug string) (*domain.User, error) {
	_, cancel := context.WithTimeout(ctx, s.ContextTimeout)
	defer cancel()

	role, err := s.RoleRepo.GetBySlug(ctx, roleSlug)

	if err != nil {
		return &domain.User{}, errors.New("role doesn't exist")
	}

	passwordBytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	password = string(passwordBytes[:])

	userToAdd := domain.User{
		Username: username,
		Password: password,
		RoleId:   role.ID,
	}

	user, err := s.UserRepo.Store(ctx, userToAdd)

	if err != nil {
		return &domain.User{}, err
	}

	user.Role = &role
	return user, nil
}

// Gets the user login JWT
func (s DefaultUserService) GetLoginJWT(ctx context.Context, username string, password string) (string, error) {
	_, cancel := context.WithTimeout(ctx, s.ContextTimeout)
	defer cancel()

	if len(username) == 0 || len(password) == 0 {
		return "", domain.ErrBadParamInput
	}

	user, err := s.UserRepo.GetByUsername(ctx, username)

	if err != nil || user.ID == uuid.Nil {
		return "", domain.ErrNotFound
	}

	userRole, err := s.RoleRepo.GetByUUID(ctx, user.RoleId)

	if err != nil {
		return "", domain.ErrNotFound
	}

	user.Role = &userRole

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", domain.ErrNotFound
	}

	return s.TokenManager.GenerateJWT(user)
}

func (s DefaultUserService) GetUserByUUID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	_, cancel := context.WithTimeout(ctx, s.ContextTimeout)
	defer cancel()

	if id == uuid.Nil {
		return &domain.User{}, domain.ErrBadParamInput
	}

	user, err := s.UserRepo.GetByUUID(ctx, id)

	if err != nil {
		return &domain.User{}, domain.ErrNotFound
	}

	userRole, err := s.RoleRepo.GetByUUID(ctx, user.RoleId)

	if err != nil {
		return &domain.User{}, domain.ErrBadParamInput
	}

	user.Role = &userRole
	return user, nil
}

func (s DefaultUserService) GetRefreshToken(ctx context.Context, user *domain.User) error {
	panic("cenas")
}
