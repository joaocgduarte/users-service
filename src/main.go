package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	_posgresConnection "github.com/plagioriginal/user-microservice/database/connection/postgres"
	"github.com/plagioriginal/user-microservice/database/migrations"
	_refreshTokensMigrations "github.com/plagioriginal/user-microservice/refresh_tokens/migrations"
	_refreshTokensRepo "github.com/plagioriginal/user-microservice/refresh_tokens/repository/postgres"
	_refreshTokensService "github.com/plagioriginal/user-microservice/refresh_tokens/service"
	_rolesMigrations "github.com/plagioriginal/user-microservice/roles/migrations"
	_rolesRepo "github.com/plagioriginal/user-microservice/roles/repository/postgres"
	"github.com/plagioriginal/user-microservice/users/handler"
	_usersMigrations "github.com/plagioriginal/user-microservice/users/migrations"
	_usersRepo "github.com/plagioriginal/user-microservice/users/repository/postgres"
	users "github.com/plagioriginal/users-service-grpc/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_usersService "github.com/plagioriginal/user-microservice/users/service"
	"github.com/plagioriginal/user-microservice/users/tokens"
)

func doMigrations(l *log.Logger, db *sql.DB) {
	migrations := migrations.MigrationsHandler{
		Logger: l,
		Db:     db,
		Migrations: []migrations.Migration{
			// Migrations
			migrations.NewAddUuidExtensionMigration(),
			_rolesMigrations.NewCreateRolesMigration(),
			_usersMigrations.NewCreateUsersMigration(),
			_refreshTokensMigrations.NewCreateRefreshTokensMigration(),
			_usersMigrations.NewAddRefreshTokenReferenceMigration(),

			// Seeds
			_rolesMigrations.NewAddRolesMigration(),
			_usersMigrations.NewAddDefaultUserMigration(),
		},
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	ctx = context.WithValue(ctx, _usersMigrations.DefaultUserNameKey, os.Getenv("DEFAULT_USER_USERNAME"))
	ctx = context.WithValue(ctx, _usersMigrations.DefaultUserPasswordKey, os.Getenv("DEFAULT_USER_PASSWORD"))
	ctx = context.WithValue(ctx, "jwtSecret", os.Getenv("JWT_GENERATOR_SECRET"))

	defer cancelfunc()
	migrations.DoAll(ctx)
}

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

	doMigrations(logger, db)

	timeoutContext := time.Duration(2) * time.Second
	jwtTokenSecret := os.Getenv("JWT_GENERATOR_SECRET")

	// Creating all the repos
	userRepo := _usersRepo.New(db)
	roleRepo := _rolesRepo.New(db)
	refreshTokenRepo := _refreshTokensRepo.New(db)

	// Creating all the services.
	refreshTokenService := _refreshTokensService.New(logger, refreshTokenRepo, userRepo, timeoutContext)
	tokenManager := tokens.NewTokenManager(jwtTokenSecret, refreshTokenService, roleRepo)
	userService := _usersService.New(logger, userRepo, roleRepo, timeoutContext)

	// @todo: refactor server instantiation.
	gs := grpc.NewServer()
	grpcServer := handler.NewUserGRPCHandler(logger, tokenManager, userService)
	users.RegisterUsersServer(gs, grpcServer)

	reflection.Register(gs)
	l, err := net.Listen("tcp", ":"+os.Getenv("API_PORT"))

	if err != nil {
		logger.Fatalln("Unable to listen to port")
		os.Exit(1)
	}

	logger.Println("gRPC Server running at port: " + os.Getenv("API_PORT"))
	if err = gs.Serve(l); err != nil {
		logger.Fatalln(err)
	}
}
