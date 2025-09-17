// Package main CRUD Service API
//
//	@title			CRUD Service API
//	@version		1.0
//	@description	A subscription management service API
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
//	@host		localhost:8080
//	@BasePath	/
//
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
package main

import (
	"context"
	"crudl_service/src/api"
	"crudl_service/src/config"
	"crudl_service/src/db"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "crudl_service/docs"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

func setupLogger() {
	level, err := log.ParseLevel(config.CurrentConfig.Server.LogLevel)
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)
}

func main() {
	setupLogger()
	log.Info("Starting CRUD service application")

	log.Info("Initializing configuration")
	config.InitConfig()

	log.Info("Initializing database connection")
	db.InitDBConnection()

	log.Info("Initializing repository")
	repository := db.NewPostgresRepository()

	log.Info("Initializing API layer")
	api.InitAPI(repository)

	log.Info("Initializing router and endpoints")
	var r = chi.NewRouter()

	log.Info("Registering API endpoints")
	r.Post("/subscription", api.CreateSubscription)
	r.Get("/subscription/{id}", api.ReadSubscription)
	r.Put("/subscription/{id}", api.UpdateSubscription)
	r.Delete("/subscription/{id}", api.DeleteSubscription)
	r.Get("/subscriptionList", api.ListSubscription)
	r.Get("/sum_subscriptions", api.SumUserSubscriptions)

	log.Info("Registering Swagger documentation")
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	server := &http.Server{
		Addr:    ":" + config.CurrentConfig.Server.Port,
		Handler: r,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.WithField("port", config.CurrentConfig.Server.Port).Info("Starting HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start HTTP server")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Info("Received shutdown signal, gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown")
	}

	log.Info("Waiting for goroutines to finish...")
	wg.Wait()

	log.Info("Shutting down API layer...")
	api.ShutdownAPI()

	log.Info("Closing database connection...")
	db.CloseDB()

	log.Info("Shutting down configuration...")
	config.ShutdownConfig()

	log.Info("Application shutdown complete")
}
