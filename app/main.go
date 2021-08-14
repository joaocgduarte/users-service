package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_posgresConnection "github.com/plagioriginal/user-microservice/database/connection/postgres"
	"github.com/plagioriginal/user-microservice/database/migrations"
	"github.com/plagioriginal/user-microservice/helpers"
	_refreshTokensMigrations "github.com/plagioriginal/user-microservice/refresh_tokens/migrations"
	_rolesMigrations "github.com/plagioriginal/user-microservice/roles/migrations"
	_rolesRepo "github.com/plagioriginal/user-microservice/roles/repository/postgres"
	"github.com/plagioriginal/user-microservice/server"
	"github.com/plagioriginal/user-microservice/users/handler"
	_usersMigrations "github.com/plagioriginal/user-microservice/users/migrations"
	_usersRepo "github.com/plagioriginal/user-microservice/users/repository/postgres"
	"github.com/plagioriginal/user-microservice/users/service"
	"github.com/plagioriginal/user-microservice/users/tokens"
)

func doMigrations(l *log.Logger, db *sql.DB) {
	migrations := migrations.MigrationsHandler{
		Logger: l,
		Db:     db,
		Migrations: []migrations.Migration{
			_rolesMigrations.NewCreateRolesMigration(),
			_rolesMigrations.NewAddRolesMigration(),
			_usersMigrations.NewCreateUsersMigration(),
			_usersMigrations.NewAddDefaultUserMigration(),
			_refreshTokensMigrations.NewCreateRefreshTokensMigration(),
			_usersMigrations.NewAddRefreshTokenReferenceMigration(),
		},
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	ctx = context.WithValue(ctx, "defaultUserUsername", os.Getenv("DEFAULT_USER_USERNAME"))
	ctx = context.WithValue(ctx, "defaultUserPassword", os.Getenv("DEFAULT_USER_PASSWORD"))
	ctx = context.WithValue(ctx, "jwtSecret", os.Getenv("JWT_GENERATOR_SECRET"))

	defer cancelfunc()

	migrations.DoAll(ctx)
}

func main() {
	logger := log.New(os.Stdout, "users-auth: ", log.Flags())

	userHandler := handler.NewHttpHandler(logger)

	serveMux := http.NewServeMux()
	serveMux.Handle("/", userHandler)

	serverConfigs := server.ServerConfigs{
		Mux:          serveMux,
		Port:         os.Getenv("API_PORT"),
		IdleTimeout:  helpers.ConvertToInt(os.Getenv("SERVER_IDLE_TIMEOUT"), 120),
		WriteTimeout: helpers.ConvertToInt(os.Getenv("SERVER_WRITE_TIMEOUT"), 2),
		ReadTimeout:  helpers.ConvertToInt(os.Getenv("SERVER_READ_TIMEOUT"), 2),
	}

	db, err := _posgresConnection.Get(_posgresConnection.PostgresConnectionSettings{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
	})

	if err != nil {
		logger.Fatal(err)
	}

	doMigrations(logger, db)

	userRepo := _usersRepo.New(db)
	roleRepo := _rolesRepo.New(db)
	tokenManager := tokens.NewTokenManager(os.Getenv("JWT_GENERATOR_SECRET"))
	timeoutContext := time.Duration(2) * time.Second

	ctx := context.Background()
	userService := service.New(userRepo, roleRepo, tokenManager, timeoutContext)
	token, err := userService.GetLoginJWT(ctx, "admin", "password")
	fmt.Println(token, err)

	newToken, _ := tokenManager.ParseJWT(token)
	userId, err := tokenManager.GetUserIDFromToken(newToken)

	loggedinUser, _ := userRepo.GetByUUID(ctx, userId)
	userRole, _ := roleRepo.GetByUUID(ctx, loggedinUser.RoleId)
	loggedinUser.Role = &userRole
	fmt.Println(loggedinUser, "Role of said user: "+loggedinUser.Role.RoleLabel)

	server := server.New(serverConfigs, logger)

	server.Init()
}
