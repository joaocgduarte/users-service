package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	_posgresConnection "github.com/plagioriginal/user-microservice/database/connection/postgres"
	"github.com/plagioriginal/user-microservice/database/migrations"
	"github.com/plagioriginal/user-microservice/domain"
	"github.com/plagioriginal/user-microservice/helpers"
	_refreshTokensMigrations "github.com/plagioriginal/user-microservice/refresh_tokens/migrations"
	_refreshTokensRepo "github.com/plagioriginal/user-microservice/refresh_tokens/repository/postgres"
	_refreshTokensService "github.com/plagioriginal/user-microservice/refresh_tokens/service"
	_rolesMigrations "github.com/plagioriginal/user-microservice/roles/migrations"
	_rolesRepo "github.com/plagioriginal/user-microservice/roles/repository/postgres"
	"github.com/plagioriginal/user-microservice/server"
	"github.com/plagioriginal/user-microservice/users/handler"
	_usersMigrations "github.com/plagioriginal/user-microservice/users/migrations"
	_usersRepo "github.com/plagioriginal/user-microservice/users/repository/postgres"

	_usersService "github.com/plagioriginal/user-microservice/users/service"
	"github.com/plagioriginal/user-microservice/users/tokens"
)

func doMigrations(l *log.Logger, db *sql.DB) {
	migrations := migrations.MigrationsHandler{
		Logger: l,
		Db:     db,
		Migrations: []migrations.Migration{
			// Migrations
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

	timeoutContext := time.Duration(2) * time.Second
	jwtTokenSecret := os.Getenv("JWT_GENERATOR_SECRET")

	// Creating all the repos
	userRepo := _usersRepo.New(db)
	roleRepo := _rolesRepo.New(db)
	refreshTokenRepo := _refreshTokensRepo.New(db)

	// Creating all the services.
	refreshTokenService := _refreshTokensService.New(refreshTokenRepo, userRepo, timeoutContext)
	tokenManager := tokens.NewTokenManager(jwtTokenSecret, refreshTokenService, roleRepo)
	userService := _usersService.New(userRepo, roleRepo, timeoutContext)

	// Doing all the actions
	ctx := context.Background()

	// refreshToken(tokenManager, ctx, userService)
	generateTokensFromLogin(userService, tokenManager, ctx)

	server := server.New(serverConfigs, logger)

	server.Init()
}

func refreshToken(tokenManager tokens.TokenManager, ctx context.Context, userService domain.UserService) {
	//Refreshes tokens by ID. Deletes the Refresh token if invalid
	oldRefreshToken, _ := uuid.Parse("18c2779d-2319-4bdc-9295-9da3edc9ddee")
	tokens, err := tokenManager.RefreshAllTokens(ctx, oldRefreshToken)
	fmt.Println(tokens, err)
}

func generateTokensFromLogin(userService domain.UserService, tokenManager tokens.TokenManager, ctx context.Context) {
	//Gets the user by login
	user, err := userService.GetUserByLogin(ctx, "admin", "password")

	fmt.Println(user, err)

	if err != nil {
		return
	}

	// Generates the tokens of said user.
	token, err := tokenManager.GenerateTokens(ctx, user)

	newToken, _ := tokenManager.ParseJWT(token.AccessToken)
	userId, err := tokenManager.GetUserIDFromToken(newToken)

	loggedinUser, _ := userService.GetUserByUUID(ctx, userId)
	fmt.Println(loggedinUser, "Role of said user: "+loggedinUser.Role.RoleLabel)
}
