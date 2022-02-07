package handler

import (
	"io/ioutil"
	"log"

	"github.com/plagioriginal/user-microservice/domain"
	users "github.com/plagioriginal/users-service-grpc/users"
)

type UserGRPCHandler struct {
	users.UnimplementedUsersServer
	l            *log.Logger
	tokenManager domain.AccessTokenHandler
	userService  domain.UserService
}

func NewUserGRPCHandler(
	l *log.Logger,
	tokenManager domain.AccessTokenHandler,
	userService domain.UserService,
) users.UsersServer {
	return UserGRPCHandler{
		l:            l,
		tokenManager: tokenManager,
		userService:  userService,
	}
}

func newHandler(tokenManager domain.AccessTokenHandler, userService domain.UserService) users.UsersServer {
	return UserGRPCHandler{
		l:            log.New(ioutil.Discard, "tests: ", log.Flags()),
		tokenManager: tokenManager,
		userService:  userService,
	}
}
