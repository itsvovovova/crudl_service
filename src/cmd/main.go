package main

import (
	"crudl_service/src/api"
	"crudl_service/src/config"
	"net/http"

	"log"

	"github.com/go-chi/chi/v5"
)

func main() {
	log.Println("Starting CRUD service application")

	log.Println("Initializing router and endpoints")
	var r = chi.NewRouter()

	log.Println("Registering API endpoints")
	r.Post("/subscription", api.CreateSubscription)
	r.Get("/subscription", api.ReadSubscription)
	r.Put("/subscription", api.UpdateSubscription)
	r.Delete("/subscription", api.DeleteSubscription)
	r.Get("/subscriptionList", api.ListSubscription)
	r.Get("/sum_subscriptions", api.SumUserSubscriptions)

	log.Println("Starting HTTP server", "port", config.CurrentConfig.Server.Port)
	err := http.ListenAndServe(":"+config.CurrentConfig.Server.Port, r)
	if err != nil {
		log.Fatal("Failed to start HTTP server")
	}
}
