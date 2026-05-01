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
	"crudl_service/src/closer"
	"crudl_service/src/config"
	"crudl_service/src/db"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "crudl_service/docs"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

func setupLogger(cfg *config.Config) {
	level, err := log.ParseLevel(cfg.Server.LogLevel)
	if err != nil {
		level = log.InfoLevel
		log.Warn("Invalid LOG_LEVEL value, defaulting to info")
	}
	log.SetLevel(level)
}

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}
	setupLogger(cfg)
	log.Info("Starting CRUD service application")

	sqlDB, err := db.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	cl := &closer.Closer{}
	cl.Add(func() error {
		log.Info("Closing database connection")
		return sqlDB.Close()
	})

	repo := db.NewPostgresRepository(sqlDB)
	app := api.NewApp(repo, cfg.JWT.SecretKey)

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Post("/register", app.RegisterUser)
	r.Post("/login", app.LoginUser)

	r.Post("/subscription", app.ValidateJWT(app.CreateSubscription))
	r.Get("/subscription/{id}", app.ValidateJWT(app.ReadSubscription))
	r.Put("/subscription/{id}", app.ValidateJWT(app.UpdateSubscription))
	r.Delete("/subscription/{id}", app.ValidateJWT(app.DeleteSubscription))
	r.Get("/subscriptionList", app.ValidateJWT(app.ListSubscription))
	r.Post("/sum_subscriptions", app.ValidateJWT(app.SumUserSubscriptions))

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	cl.Add(func() error {
		log.Info("Shutting down HTTP server")
		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.WithField("port", cfg.Server.Port).Info("Starting HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Info("Received shutdown signal, gracefully shutting down...")
	stop() // повторный Ctrl+C убьёт процесс немедленно

	if err := cl.Close(); err != nil {
		log.WithError(err).Error("Error during shutdown")
	}
	wg.Wait()
	log.Info("Application shutdown complete")
}
