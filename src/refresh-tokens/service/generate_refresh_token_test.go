package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
	"github.com/plagioriginal/user-microservice/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerateRefreshToken_InvalidInput(t *testing.T) {
	res, err := newService(nil, nil).
		GenerateRefreshToken(context.TODO(), nil)
	assert.Empty(t, res)
	assert.Error(t, err)
	assert.Equal(t, err, domain.ErrBadParamInput)
}

func TestGenerateRegreshToken_UserDoesntHaveToken_ErrorOnStoringToken(t *testing.T) {
	user := &domain.User{}
	tokenRepo := new(mocks.RefreshTokenRepository)
	tokenRepo.On("Store", mock.Anything, mock.MatchedBy(func(rt domain.RefreshToken) bool {
		return rt.ValidUntil.After(time.Now()) && rt.Token != uuid.Nil
	})).Once().Return(domain.RefreshToken{}, errors.New("boom"))

	res, err := newService(tokenRepo, nil).
		GenerateRefreshToken(context.TODO(), user)
	assert.Empty(t, res)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
	tokenRepo.AssertExpectations(t)
}

func TestGenerateRegreshToken_UserDoesntHaveToken_ErrorSavingTokenReferenceAndDeleting(t *testing.T) {
	user := &domain.User{}
	tokenRepo := new(mocks.RefreshTokenRepository)
	userRepo := new(mocks.UserRepository)
	rt := domain.RefreshToken{
		Id:         uuid.New(),
		Token:      uuid.New(),
		ValidUntil: time.Now().Add(time.Hour * 24 * 7),
	}
	tokenRepo.On("Store", mock.Anything, mock.MatchedBy(func(rt domain.RefreshToken) bool {
		return rt.ValidUntil.After(time.Now()) && rt.Token != uuid.Nil
	})).Once().Return(rt, nil)
	tokenRepo.On("Delete", mock.Anything, rt.Id).Once().Return(errors.New("boom2"))

	userRepo.On("SaveRefreshToken", mock.Anything, user, rt).Once().Return(errors.New("boom"))

	res, err := newService(tokenRepo, userRepo).
		GenerateRefreshToken(context.TODO(), user)
	assert.Empty(t, res)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom2")
	tokenRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestGenerateRegreshToken_UserDoesntHaveToken_ErrorSavingTokenReference(t *testing.T) {
	user := &domain.User{}
	tokenRepo := new(mocks.RefreshTokenRepository)
	userRepo := new(mocks.UserRepository)
	rt := domain.RefreshToken{
		Id:         uuid.New(),
		Token:      uuid.New(),
		ValidUntil: time.Now().Add(time.Hour * 24 * 7),
	}
	tokenRepo.On("Store", mock.Anything, mock.MatchedBy(func(rt domain.RefreshToken) bool {
		return rt.ValidUntil.After(time.Now()) && rt.Token != uuid.Nil
	})).Once().Return(rt, nil)
	tokenRepo.On("Delete", mock.Anything, rt.Id).Once().Return(nil)

	userRepo.On("SaveRefreshToken", mock.Anything, user, rt).Once().Return(errors.New("boom"))

	res, err := newService(tokenRepo, userRepo).
		GenerateRefreshToken(context.TODO(), user)
	assert.Empty(t, res)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
	tokenRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestGenerateRegreshToken_UserDoesntHaveToken_Success(t *testing.T) {
	user := &domain.User{}
	tokenRepo := new(mocks.RefreshTokenRepository)
	userRepo := new(mocks.UserRepository)
	rt := domain.RefreshToken{
		Id:         uuid.New(),
		Token:      uuid.New(),
		ValidUntil: time.Now().Add(time.Hour * 24 * 7),
	}
	tokenRepo.On("Store", mock.Anything, mock.MatchedBy(func(rt domain.RefreshToken) bool {
		return rt.ValidUntil.After(time.Now()) && rt.Token != uuid.Nil
	})).Once().Return(rt, nil)

	userRepo.On("SaveRefreshToken", mock.Anything, user, rt).Once().Return(nil)

	res, err := newService(tokenRepo, userRepo).
		GenerateRefreshToken(context.TODO(), user)
	assert.Equal(t, res, rt)
	assert.Nil(t, err)
	tokenRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestGenerateRegreshToken_UserHasAToken_CantFindTokenInRepo(t *testing.T) {
	user := &domain.User{
		RefreshTokenId: uuid.NullUUID{
			UUID:  uuid.New(),
			Valid: true,
		},
	}
	tokenRepo := new(mocks.RefreshTokenRepository)
	tokenRepo.On("GetByUUID", mock.Anything, user.RefreshTokenId.UUID).
		Once().Return(domain.RefreshToken{}, errors.New("boom"))

	res, err := newService(tokenRepo, nil).
		GenerateRefreshToken(context.TODO(), user)
	assert.Empty(t, res)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
	tokenRepo.AssertExpectations(t)
}

func TestGenerateRegreshToken_UserHasAToken_TokenIsExpired(t *testing.T) {
	user := &domain.User{
		RefreshTokenId: uuid.NullUUID{
			UUID:  uuid.New(),
			Valid: true,
		},
	}
	tokenRepo := new(mocks.RefreshTokenRepository)
	tokenRepo.On("GetByUUID", mock.Anything, user.RefreshTokenId.UUID).
		Once().
		Return(domain.RefreshToken{
			ValidUntil: time.Now().Add(time.Hour * -24),
		}, nil)

	res, err := newService(tokenRepo, nil).
		GenerateRefreshToken(context.TODO(), user)
	assert.Empty(t, res)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid refresh token")
	tokenRepo.AssertExpectations(t)
}

func TestGenerateRegreshToken_UserHasAToken_ErrorDeletingOldToken(t *testing.T) {
	user := &domain.User{
		RefreshTokenId: uuid.NullUUID{
			UUID:  uuid.New(),
			Valid: true,
		},
	}
	tokenRepo := new(mocks.RefreshTokenRepository)
	tokenRepo.On("GetByUUID", mock.Anything, user.RefreshTokenId.UUID).
		Once().
		Return(domain.RefreshToken{
			Id:         user.RefreshTokenId.UUID,
			ValidUntil: time.Now().Add(time.Hour * 7),
		}, nil)
	tokenRepo.On("Delete", mock.Anything, user.RefreshTokenId.UUID).
		Once().
		Return(errors.New("boom"))

	res, err := newService(tokenRepo, nil).
		GenerateRefreshToken(context.TODO(), user)
	assert.Empty(t, res)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
	tokenRepo.AssertExpectations(t)
}

func TestGenerateRegreshToken_UserHasAToken_Success(t *testing.T) {
	user := &domain.User{
		RefreshTokenId: uuid.NullUUID{
			UUID:  uuid.New(),
			Valid: true,
		},
	}
	tokenRepo := new(mocks.RefreshTokenRepository)
	userRepo := new(mocks.UserRepository)
	rt := domain.RefreshToken{
		Id:         uuid.New(),
		Token:      uuid.New(),
		ValidUntil: time.Now().Add(time.Hour * 24 * 7),
	}
	tokenRepo.On("GetByUUID", mock.Anything, user.RefreshTokenId.UUID).
		Once().
		Return(domain.RefreshToken{
			Id:         user.RefreshTokenId.UUID,
			ValidUntil: time.Now().Add(time.Hour * 7),
		}, nil)
	tokenRepo.On("Delete", mock.Anything, user.RefreshTokenId.UUID).
		Once().
		Return(nil)

	tokenRepo.On("Store", mock.Anything, mock.MatchedBy(func(rt domain.RefreshToken) bool {
		return rt.ValidUntil.After(time.Now()) && rt.Token != uuid.Nil
	})).Once().Return(rt, nil)

	userRepo.On("SaveRefreshToken", mock.Anything, user, rt).Once().Return(nil)

	res, err := newService(tokenRepo, userRepo).
		GenerateRefreshToken(context.TODO(), user)
	assert.Equal(t, res, rt)
	assert.Nil(t, err)
	tokenRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}
