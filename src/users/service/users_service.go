package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
	"golang.org/x/crypto/bcrypt"
)

type DefaultUserService struct {
	Logger         *log.Logger
	UserRepo       domain.UserRepository
	RoleRepo       domain.RoleRepository
	ContextTimeout time.Duration
}

// Constructor
func New(
	logger *log.Logger,
	userRepo domain.UserRepository,
	roleRepo domain.RoleRepository,
	contextTimeout time.Duration,
) domain.UserService {
	return DefaultUserService{
		logger,
		userRepo,
		roleRepo,
		contextTimeout,
	}
}

// Stores q new user based on username and password
func (s DefaultUserService) Store(ctx context.Context, request domain.StoreUserRequest) (*domain.User, error) {
	_, cancel := context.WithTimeout(ctx, s.ContextTimeout)
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

// Gets a user by the username and password. With role attached
func (s DefaultUserService) GetUserByLogin(ctx context.Context, request domain.GetUserRequest) (*domain.User, error) {
	_, cancel := context.WithTimeout(ctx, s.ContextTimeout)
	defer cancel()

	if len(request.Username) == 0 || len(request.Password) == 0 {
		return nil, domain.ErrBadParamInput
	}

	user, err := s.UserRepo.GetByUsername(ctx, request.Username)
	if err != nil || user.ID == uuid.Nil {
		s.Logger.Println("error getting user by username: " + err.Error())
		return nil, domain.ErrNotFound
	}

	userRole, err := s.RoleRepo.GetByUUID(ctx, user.RoleId)
	if err != nil {
		s.Logger.Println("error fetching role of user: " + err.Error())
		return nil, domain.ErrNotFound
	}

	user.Role = &userRole
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))

	if err != nil {
		return nil, domain.ErrNotFound
	}
	return user, nil
}

func (s DefaultUserService) GetUserByUUID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	_, cancel := context.WithTimeout(ctx, s.ContextTimeout)
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
