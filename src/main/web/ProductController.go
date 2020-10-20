package web

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/hsedjame/product-psql-api/src/main/core"
	"github.com/hsedjame/product-psql-api/src/main/models"
	"net/http"
)

type ProductKey struct{}

type ProductController struct {
	Handler ProductHandler
}

func NewProductController(handler ProductHandler) *ProductController {
	return &ProductController{Handler: handler}
}

func (controller *ProductController) Path() string {
	return "/products"
}

func (controller *ProductController) AddRoutes(router *mux.Router) {

	router.HandleFunc("", controller.Handler.All).Methods(http.MethodGet)
	router.HandleFunc("", controller.Handler.Create).Methods(http.MethodPost)
	router.HandleFunc("", controller.Handler.Update).Methods(http.MethodPut)
	router.HandleFunc("/{id:[0-9]+}", controller.Handler.Delete).Methods(http.MethodDelete)

}

func (controller *ProductController) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {
		wr.Header().Add("Content-Type", "application/json")

		/*
		 * Dans le cas d'une methode POST ou PUT
		 * Valider la requête
		 */
		if rq.Method == http.MethodPost || rq.Method == http.MethodPut {
			var product models.Product
			if err := core.FromJson(&product, rq.Body); err != nil {
				wr.WriteHeader(http.StatusBadRequest)
				_ = core.ToJson(core.AppError{Message: err.Error()}, wr)
				return
			} else if err := core.IsValid(product); err != nil {
				wr.WriteHeader(http.StatusBadRequest)
				_ = core.ToJson(core.AppError{Message: err.Error()}, wr)
				return
			}

			ctx := context.WithValue(rq.Context(), ProductKey{}, product)

			rq := rq.WithContext(ctx)

			next.ServeHTTP(wr, rq)

			return
		}

		next.ServeHTTP(wr, rq)
	})
}
