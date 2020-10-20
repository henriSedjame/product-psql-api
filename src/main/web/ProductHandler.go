package web

import (
	"github.com/gorilla/mux"
	"github.com/hsedjame/product-psql-api/src/main/core"
	"github.com/hsedjame/product-psql-api/src/main/models"
	"github.com/hsedjame/product-psql-api/src/main/repository"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	Repository repository.ProductRepository
}

func NewProductHandler(repository repository.ProductRepository) *ProductHandler {
	return &ProductHandler{Repository: repository}
}

func (handler *ProductHandler) All(wr http.ResponseWriter, _ *http.Request) {

	if products, err := handler.Repository.FindAll(); err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		_ = core.ToJson(core.AppError{Message: err.Error()}, wr)
		return
	} else if err := core.ToJson(products, wr); err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		_ = core.ToJson(core.AppError{Message: err.Error()}, wr)
		return
	}

}

func (handler *ProductHandler) Create(wr http.ResponseWriter, rq *http.Request) {

	product := rq.Context().Value(ProductKey{}).(models.Product)

	if err := handler.Repository.Create(&product); err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		_ = core.ToJson(core.AppError{Message: err.Error()}, wr)
		return
	} else {
		wr.WriteHeader(http.StatusOK)
		return
	}

}

func (handler *ProductHandler) Update(wr http.ResponseWriter, rq *http.Request) {

	product := rq.Context().Value(ProductKey{}).(models.Product)

	if err := handler.Repository.Update(&product); err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		_ = core.ToJson(core.AppError{Message: err.Error()}, wr)
		return
	} else {
		wr.WriteHeader(http.StatusOK)
		return
	}
}

func (handler *ProductHandler) Delete(wr http.ResponseWriter, rq *http.Request) {

	pathParams := mux.Vars(rq)
	if idString := pathParams["id"]; idString != "" {

		if id, err := strconv.Atoi(idString); err != nil {

			wr.WriteHeader(http.StatusBadRequest)
			_ = core.ToJson(core.AppError{Message: err.Error()}, wr)
			return

		} else if err := handler.Repository.Delete(id); err != nil {

			wr.WriteHeader(http.StatusBadRequest)
			_ = core.ToJson(core.AppError{Message: err.Error()}, wr)
			return

		} else {
			wr.WriteHeader(http.StatusOK)
			return
		}
	}

}
