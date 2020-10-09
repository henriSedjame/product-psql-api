package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

type RestController interface {
	Path() string
	AddRoutes(router *mux.Router)
	Middleware(next http.Handler) http.Handler
}
