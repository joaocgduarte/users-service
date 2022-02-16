package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/plagioriginal/user-microservice/database"
	_posgresConnection "github.com/plagioriginal/user-microservice/database/connection/postgres"
	_refreshTokensRepo "github.com/plagioriginal/user-microservice/refresh-tokens/repository/postgres"
	_refreshTokensService "github.com/plagioriginal/user-microservice/refresh-tokens/service"
	_rolesRepo "github.com/plagioriginal/user-microservice/roles/repository/postgres"
	"github.com/plagioriginal/user-microservice/users/handler"
	_usersRepo "github.com/plagioriginal/user-microservice/users/repository/postgres"
	users "github.com/plagioriginal/users-service-grpc/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_usersService "github.com/plagioriginal/user-microservice/users/service"
	"github.com/plagioriginal/user-microservice/users/tokens"
)

func main() {
	logger := log.New(os.Stdout, "users-auth: ", log.Flags())

	db, err := _posgresConnection.Get(_posgresConnection.PostgresConnectionSettings{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
	})
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	jwtTokenSecret := os.Getenv("JWT_GENERATOR_SECRET")
	timeoutContext := time.Duration(2) * time.Second

	database.DoMigrations(logger, db, database.MigrationSettings{
		DefaultUserUsername: os.Getenv("DEFAULT_USER_USERNAME"),
		DefaultUserPassword: os.Getenv("DEFAULT_USER_PASSWORD"),
		JwtSecret:           jwtTokenSecret,
		Timeout:             timeoutContext,
	})

	// Creating all the repos
	userRepo := _usersRepo.New(db)
	roleRepo := _rolesRepo.New(db)
	refreshTokenRepo := _refreshTokensRepo.New(db)

	// Creating all the services.
	refreshTokenService := _refreshTokensService.New(logger, refreshTokenRepo, userRepo, timeoutContext)
	tokenManager := tokens.NewTokenManager(jwtTokenSecret, refreshTokenService, roleRepo)
	userService := _usersService.New(userRepo, roleRepo, timeoutContext)

	// @todo: refactor server instantiation.
	gs := grpc.NewServer()
	grpcServer := handler.NewUserGRPCHandler(logger, tokenManager, userService)
	users.RegisterUsersServer(gs, grpcServer)

	reflection.Register(gs)
	l, err := net.Listen("tcp", ":"+os.Getenv("APP_PORT"))

	if err != nil {
		logger.Fatalln("Unable to listen to port")
		os.Exit(1)
	}

	logger.Println("gRPC Server running at port: " + os.Getenv("APP_PORT"))
	if err = gs.Serve(l); err != nil {
		logger.Fatalln(err)
	}
}
