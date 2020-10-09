package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hsedjame/product-psql-api/src/main/core"
	"github.com/hsedjame/product-psql-api/src/main/web"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type App struct {
	Server *http.Server
	DB *pg.DB
	Properties core.AppProperties
	Log *log.Logger
	Controllers []web.RestController
	Classpath string
	IsInitialized bool
}


// Initialize the application
func (app *App) Initialize() {

	if err := app.setClasspath(); err != nil {
		app.Log.Fatal(err)
	}

	if err, msg := app.setAppProperties(); err != nil {
		if msg == "" {
			app.Log.Fatal(err)
		} else {
			app.Log.Fatalf("%s => %s", msg, err)
		}
	}

	if err := app.openDatabase(); err != nil {
		app.Log.Fatal(err)
	}

	app.configureServer()

	app.IsInitialized = true
}

// Run the application
func (app *App) Run() {
	if !app.IsInitialized {
		app.Initialize()
	}

	go func() {
		app.Log.Fatal(app.Server.ListenAndServe())
	}()

	app.shutDown()
}

func (app *App) shutDown() {

	signalChannel := make(chan os.Signal)

	// Envoyer un signal au channel lors :
	// * d'une interruption de la machine
	// * d'un arrêt de la machine
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	// Récupérer le signal envoyé

	_ = <-signalChannel

	app.Log.Println(" ===> Arrêt du serveur")

	deadline, _ := context.WithTimeout(context.Background(), 30 * time.Second)

	if err := app.DB.Close(); err != nil {
		app.Log.Fatal(err)
	}

	if err := app.Server.Shutdown(deadline); err != nil {
		app.Log.Fatal(err)
	}

}

func (app *App) setClasspath() error {
	if app.Classpath == "" {
		if rootDir, err := os.Getwd(); err != nil {
			return err
		} else {
			app.Classpath = fmt.Sprintf("%s/src/resources", rootDir)

			app.Log.Printf(" ===> Application classpath configured to : %s", app.Classpath)
		}
	}

	return nil
}

func (app *App) setAppProperties() (error, string) {

	props := core.DefaultAppProperties()

	err, s, done := app.setProfileProps("", props)

	if done {
		return err, s
	} else {
		profiles := app.Properties.ActiveProfiles
		if profiles != "" {

			app.Log.Printf(" ===> Application active profiles : %s", profiles)

			for _, profile := range strings.Split(profiles, ",") {
				err, s, done := app.setProfileProps(profile, app.Properties)
				if done {
					return err, s
				}
			}
		}

	}

	app.Log.Println(" ===> Application Properties retrieval succeeded.", )

	return nil, ""
}

func (app *App) setProfileProps(profile string, props core.AppProperties) (error, string, bool) {

	suffix := "application"

	if profile != "" {
		suffix = fmt.Sprintf("%s-%s", suffix, profile)
	}

	propertiesLocation := fmt.Sprintf("%s/%s.json", app.Classpath, suffix)

	if propertiesFile, err := os.OpenFile(propertiesLocation, os.O_RDONLY, 0777); err != nil {
		return err, "", true
	} else {
		defer propertiesFile.Close()

		if bytes, err := ioutil.ReadAll(propertiesFile); err != nil {
			return err, "An error occured while retrieving application properties", true
		} else {
			if err := json.Unmarshal(bytes, &props); err != nil {
				return err, "An error occured while retrieving application properties => %s", true
			}
		}
	}

	app.Properties = props

	return nil, "", false
}

func (app *App) openDatabase() error {

	dbProps := app.Properties.Datasource

	if dbProps.Dbname != "" {

		optStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			dbProps.Username, dbProps.Password, dbProps.Host, dbProps.Port, dbProps.Dbname)

		if opt, err := pg.ParseURL(optStr); err != nil {
			return err
		} else {
			app.DB = pg.Connect(opt)

			ctx := context.Background()

			var version string
			if _, err := app.DB.QueryOneContext(ctx, pg.Scan(&version), "SELECT version()"); err != nil {
				app.Log.Fatal(err)
			} else {
				app.Log.Println(" ===> Connection to Database succeeded.")
				app.Log.Printf(" ===> Database version : %s", version)
			}
		}
	} else {
		return core.AppError{
			Message: "Database is not set. Please consider to set property { \"datasource\" : {\"database\" : ****}}",
		}
	}
	return nil
}

func (app *App) configureServer() {

	router := mux.NewRouter()

	for _,controller := range app.Controllers {
		subRouter := router.PathPrefix(controller.Path()).Subrouter()
		subRouter.Use(controller.Middleware)
		controller.AddRoutes(subRouter)
	}

	/* Confifure redocs */
	redocOpts := middleware.RedocOpts{SpecURL: "../../swagger.yaml"}
	redoc := middleware.Redoc(redocOpts, nil)
	router.Handle("/docs", redoc).Methods(http.MethodGet)
	router.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	/* Configure CORS */
	var opts []handlers.CORSOption
	cors := app.Properties.Cors
	origins := cors.AllowedOrigins
	headers := cors.AllowedHeaders
	methods := cors.AllowedMethods

	if origins != "" {
		opts = append(opts, handlers.AllowedOrigins(strings.Split(origins, ",")))
	}
	if headers != "" {
		opts = append(opts, handlers.AllowedHeaders(strings.Split(headers, ",")))
	}
	if methods != "" {
		opts = append(opts, handlers.AllowedHeaders(strings.Split(methods, ",")))
	}

	corsHandlers := handlers.CORS(opts...)

	/* Configure server */
	port := app.Properties.Server.Port
	if port == 0 {
		port = 8080
	}

	app.Server = &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: corsHandlers(router),
		IdleTimeout:       120 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

}


