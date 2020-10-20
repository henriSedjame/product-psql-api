package main

import (
	"github.com/hsedjame/product-psql-api/src/main/repository"
	"github.com/hsedjame/product-psql-api/src/main/web"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	app := App{Log: logger}

	app.Initialize()

	productRepo := repository.NewProductRepository(app.DB)
	handler := web.NewProductHandler(*productRepo)
	prodControl := web.NewProductController(*handler)
	app.Controllers = []web.RestController{prodControl}

	app.Run()

}
