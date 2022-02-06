package service

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
	"github.com/plagioriginal/user-microservice/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func newService(userRepo domain.UserRepository, roleRepo domain.RoleRepository) domain.UserService {
	return New(
		log.New(ioutil.Discard, "tests: ", log.Flags()),
		userRepo,
		roleRepo,
		time.Duration(5*time.Second),
	)
}

func Test_Store_FailIfInvalidRole(t *testing.T) {
	roleRepo := new(mocks.RoleRepository)
	roleRepo.On("GetBySlug", mock.Anything, "inexistant").
		Once().Return(domain.Role{}, errors.New("boom"))

	service := newService(nil, roleRepo)

	user, err := service.Store(context.TODO(), domain.StoreUserRequest{
		Username: "",
		Password: "",
		RoleSlug: "inexistant",
	})

	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error fetching role")
}

func Test_Store_FailIfErrorStoringUser(t *testing.T) {
	uuid := uuid.New()
	roleRepo := new(mocks.RoleRepository)
	roleRepo.On("GetBySlug", mock.Anything, "admin").
		Once().Return(domain.Role{ID: uuid}, nil)

	userRepo := new(mocks.UserRepository)
	userRepo.On("Store", mock.Anything, mock.MatchedBy(func(user domain.User) bool {
		return user.Username == "username" && user.RoleId.String() == uuid.String() && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password")) == nil
	})).Once().Return(nil, errors.New("boom"))

	service := newService(userRepo, roleRepo)

	user, err := service.Store(context.TODO(), domain.StoreUserRequest{
		Username: "username",
		Password: "password",
		RoleSlug: "admin",
	})

	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error storing user")
}

func Test_Store_Success(t *testing.T) {
	uuid := uuid.New()
	role := domain.Role{ID: uuid}
	roleRepo := new(mocks.RoleRepository)
	roleRepo.On("GetBySlug", mock.Anything, "admin").
		Once().Return(role, nil)

	userRepo := new(mocks.UserRepository)
	userRepo.On("Store", mock.Anything, mock.MatchedBy(func(user domain.User) bool {
		return user.Username == "username" && user.RoleId.String() == uuid.String() && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password")) == nil
	})).Once().Return(&domain.User{
		Username: "username",
		Password: "password",
		RoleId:   uuid,
	}, nil)

	service := newService(userRepo, roleRepo)

	user, err := service.Store(context.TODO(), domain.StoreUserRequest{
		Username: "username",
		Password: "password",
		RoleSlug: "admin",
	})

	assert.NotNil(t, user)
	assert.Equal(t, user, &domain.User{
		Username: "username",
		Password: "password",
		RoleId:   uuid,
		Role:     &role,
	})
	assert.Nil(t, err)
}

func Test_GetUserByLogin_FailIfInvalidInput(t *testing.T) {
	requests := []domain.GetUserRequest{
		{
			Username: "",
			Password: "",
		},
		{
			Username: "cenas",
			Password: "",
		},
		{
			Username: "",
			Password: "cenas",
		},
	}

	for _, req := range requests {
		service := newService(nil, nil)
		user, err := service.GetUserByLogin(context.TODO(), req)
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid parameter")
	}
}

func Test_GetUserByLogin_FailIfGetByUsernameError(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	userRepo.On("GetByUsername", mock.Anything, "username").Once().
		Return(nil, errors.New("boom"))

	service := newService(userRepo, nil)
	user, err := service.GetUserByLogin(context.TODO(), domain.GetUserRequest{Username: "username", Password: "casd"})
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "resource not found")
}

func Test_GetUserByLogin_FailIfRoleFetchingError(t *testing.T) {
	roleUuid := uuid.New()

	userRepo := new(mocks.UserRepository)
	userRepo.On("GetByUsername", mock.Anything, "username").Once().
		Return(&domain.User{
			ID:     uuid.New(),
			RoleId: roleUuid,
		}, nil)

	roleRepo := new(mocks.RoleRepository)
	roleRepo.On("GetByUUID", mock.Anything, roleUuid).
		Once().Return(domain.Role{}, errors.New("boom"))

	service := newService(userRepo, roleRepo)
	user, err := service.GetUserByLogin(context.TODO(), domain.GetUserRequest{Username: "username", Password: "casd"})
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "resource not found")
}

func Test_GetUserByLogin_FailIfPasswordsDontMatchError(t *testing.T) {
	roleUuid := uuid.New()

	userRepo := new(mocks.UserRepository)
	userRepo.On("GetByUsername", mock.Anything, "username").Once().
		Return(&domain.User{
			ID:       uuid.New(),
			RoleId:   roleUuid,
			Password: "dont match password",
		}, nil)

	role := domain.Role{ID: roleUuid}
	roleRepo := new(mocks.RoleRepository)
	roleRepo.On("GetByUUID", mock.Anything, roleUuid).
		Once().Return(role, nil)

	service := newService(userRepo, roleRepo)
	user, err := service.GetUserByLogin(context.TODO(), domain.GetUserRequest{Username: "username", Password: "casd"})
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "resource not found")
}

func Test_GetUserByLogin_Success(t *testing.T) {
	roleUuid := uuid.New()
	userUuid := uuid.New()

	passwordBytes, err := bcrypt.GenerateFromPassword([]byte("password"), 14)
	assert.Nil(t, err)
	password := string(passwordBytes[:])

	userRepo := new(mocks.UserRepository)
	userRepo.On("GetByUsername", mock.Anything, "username").Once().
		Return(&domain.User{
			ID:       userUuid,
			RoleId:   roleUuid,
			Password: password,
		}, nil)

	role := domain.Role{ID: roleUuid}
	roleRepo := new(mocks.RoleRepository)
	roleRepo.On("GetByUUID", mock.Anything, roleUuid).
		Once().Return(role, nil)

	service := newService(userRepo, roleRepo)
	user, err := service.GetUserByLogin(context.TODO(), domain.GetUserRequest{Username: "username", Password: "password"})
	assert.Equal(t, user, &domain.User{
		ID:       userUuid,
		RoleId:   roleUuid,
		Password: password,
		Role:     &role,
	})
	assert.Nil(t, err)
}

func Test_GetUserUUID_FailIfGetByUsernameError(t *testing.T) {
	uuid := uuid.New()
	userRepo := new(mocks.UserRepository)
	userRepo.On("GetByUUID", mock.Anything, uuid).Once().
		Return(nil, errors.New("boom"))

	service := newService(userRepo, nil)
	user, err := service.GetUserByUUID(context.TODO(), uuid)
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "resource not found")
}

func Test_GetUserUUID_FailIfRoleFetchingError(t *testing.T) {
	userUuid := uuid.New()
	roleUuid := uuid.New()

	userRepo := new(mocks.UserRepository)
	userRepo.On("GetByUUID", mock.Anything, userUuid).Once().
		Return(&domain.User{
			ID:     userUuid,
			RoleId: roleUuid,
		}, nil)

	roleRepo := new(mocks.RoleRepository)
	roleRepo.On("GetByUUID", mock.Anything, roleUuid).
		Once().Return(domain.Role{}, errors.New("boom"))

	service := newService(userRepo, roleRepo)
	user, err := service.GetUserByUUID(context.TODO(), userUuid)
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "resource not found")
}

func Test_GetUserUUID_Success(t *testing.T) {
	userUuid := uuid.New()
	roleUuid := uuid.New()

	userRepo := new(mocks.UserRepository)
	userRepo.On("GetByUUID", mock.Anything, userUuid).Once().
		Return(&domain.User{
			ID:     userUuid,
			RoleId: roleUuid,
		}, nil)

	role := domain.Role{ID: roleUuid}
	roleRepo := new(mocks.RoleRepository)
	roleRepo.On("GetByUUID", mock.Anything, roleUuid).
		Once().Return(role, nil)

	service := newService(userRepo, roleRepo)
	user, err := service.GetUserByUUID(context.TODO(), userUuid)
	assert.Equal(t, user, &domain.User{
		ID:     userUuid,
		RoleId: roleUuid,
		Role:   &role,
	})
	assert.Nil(t, err)
}
