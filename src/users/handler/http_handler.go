package handler

import (
	"fmt"
	"log"
	"net/http"
)

type HttpHandler struct {
	logger *log.Logger
}

func NewHttpHandler(logger *log.Logger) http.Handler {
	return HttpHandler{logger}
}

func (handler HttpHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)

	name := "World"

	if len(r.URL.Path[1:]) > 0 {
		name = r.URL.Path[1:]
	}

	returnString := fmt.Sprintf("Hello, %s!", name)
	rw.Write([]byte(returnString))
}
