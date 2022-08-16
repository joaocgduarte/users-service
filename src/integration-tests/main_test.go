package integration_tests

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/plagioriginal/user-microservice/database"
	"github.com/plagioriginal/user-microservice/domain"
	_refreshTokensRepo "github.com/plagioriginal/user-microservice/refresh-tokens/repository/postgres"
	_refreshTokensService "github.com/plagioriginal/user-microservice/refresh-tokens/service"
	_rolesRepo "github.com/plagioriginal/user-microservice/roles/repository/postgres"
	"github.com/plagioriginal/user-microservice/users/handler"
	_usersRepo "github.com/plagioriginal/user-microservice/users/repository/postgres"
	_usersService "github.com/plagioriginal/user-microservice/users/service"
	"github.com/plagioriginal/user-microservice/users/tokens"
	users "github.com/plagioriginal/users-service-grpc/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var (
	db               *sql.DB
	refreshTokenRepo domain.RefreshTokenRepository
	userClient       users.UsersClient
	databaseSettings database.MigrationSettings
)

type testDatabaseSettings struct {
	PostgresUser     string
	PostgresPassword string
	PostgresDb       string
}

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	logger := log.New(os.Stdout, "integration-tests: ", log.Flags())
	pool := createPool(logger)

	settings := testDatabaseSettings{
		PostgresUser:     "user_name",
		PostgresPassword: "secret",
		PostgresDb:       "dbname",
	}
	resource := createDatabaseContainer(logger, settings, pool)
	generateDataSource(logger, pool, resource)
	defer db.Close()

	//Do the db migrations
	databaseSettings = database.MigrationSettings{
		DefaultUserUsername: "default-user",
		DefaultUserPassword: "default-password",
		JwtSecret:           "secret",
		Timeout:             time.Duration(5) * time.Second,
		BcryptCost:          _usersService.TestingBcryptCost,
	}

	database.DoMigrations(logger, db, databaseSettings)

	grpcServerStarter, grpcServerFinisher, listener := setupGrpcServer(db, logger)
	defer grpcServerFinisher()
	go grpcServerStarter()

	clientCloser := setupGrpcClient(listener, logger)
	defer clientCloser()

	//Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		logger.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createPool(logger *log.Logger) *dockertest.Pool {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}
	return pool
}

func createDatabaseContainer(logger *log.Logger, settings testDatabaseSettings, pool *dockertest.Pool) *dockertest.Resource {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_PASSWORD=" + settings.PostgresPassword,
			"POSTGRES_USER=" + settings.PostgresUser,
			"POSTGRES_DB=" + settings.PostgresDb,
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		logger.Fatalf("could not start resource: %s", err)
	}
	resource.Expire(120)
	pool.MaxWait = 120 * time.Second
	return resource
}

func generateDataSource(logger *log.Logger, pool *dockertest.Pool, resource *dockertest.Resource) {
	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)
	logger.Println("connecting to database on url: ", databaseUrl)

	var err error
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		logger.Fatalf("could not connect to docker: %s", err)
	}
}

func setupGrpcServer(dbSource *sql.DB, logger *log.Logger) (func(), func(), *bufconn.Listener) {
	// Creating all the repos
	userRepo := _usersRepo.New(db)
	roleRepo := _rolesRepo.New(db)
	refreshTokenRepo = _refreshTokensRepo.New(db)

	// Creating all the services.
	refreshTokenService := _refreshTokensService.New(logger, refreshTokenRepo, userRepo, time.Duration(10*time.Second))
	tokenManager := tokens.NewTokenManager("secret", refreshTokenService, roleRepo)
	userService := _usersService.New(userRepo, roleRepo, time.Duration(10*time.Second), _usersService.TestingBcryptCost)

	gs := grpc.NewServer()
	handler := handler.NewUserGRPCHandler(logger, tokenManager, userService)
	users.RegisterUsersServer(gs, handler)

	listener := bufconn.Listen(1024 * 1024)
	return func() {
			if err := gs.Serve(listener); err != nil {
				logger.Fatalf("failed to serve grpc test server: %v", err)
			}
		}, func() {
			logger.Println("stopping grpc server...")
			gs.Stop()
		}, listener
}

func setupGrpcClient(listener *bufconn.Listener, logger *log.Logger) func() {
	grpcDialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
	connection, err := grpc.Dial("", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(grpcDialer))
	if err != nil {
		log.Fatalf("failed to generate grpc client")
	}

	userClient = users.NewUsersClient(connection)

	return func() {
		connection.Close()
	}
}
