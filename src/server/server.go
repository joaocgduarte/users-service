package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/plagioriginal/user-microservice/helpers"
)

type ServerConfigs struct {
	Mux          *http.ServeMux
	Port         string
	IdleTimeout  int
	WriteTimeout int
	ReadTimeout  int
}

type Server struct {
	srv    *http.Server
	logger *log.Logger
}

// Returns a new instance of an http.Server based on the environment variables
func New(configs ServerConfigs, logger *log.Logger) Server {
	serveAt := fmt.Sprintf(":%s", configs.Port)

	srv := &http.Server{
		Addr:         serveAt,
		Handler:      configs.Mux,
		IdleTimeout:  time.Duration(configs.IdleTimeout) * time.Second,
		WriteTimeout: time.Duration(configs.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(configs.ReadTimeout) * time.Second,
	}

	return Server{srv, logger}
}

func (s Server) Init() {
	// Go routine to begin the server
	go func() {
		s.logger.Println("Listening to port " + s.srv.Addr)
		err := s.srv.ListenAndServe()

		if err != nil {
			s.logger.Fatalln(err)
		}
	}()

	// Wait for an interrupt
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit

	// Attempt a graceful shutdown
	timeoutContext, cancel := context.
		WithTimeout(
			context.Background(),
			time.Duration(helpers.ConvertToInt(os.Getenv("SERVER_WRITE_TIMEOUT"), 2))*time.Second,
		)
	defer cancel()

	log.Println("Shutting down server...")

	if err := s.srv.Shutdown(timeoutContext); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}
