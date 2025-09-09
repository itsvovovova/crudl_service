package main

import (
	"crudl_service/src/api"
	"crudl_service/src/config"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	// metrics := metrics.NewMetrics()
	var r = chi.NewRouter()
	// r.Use(metrics.Middleware)
	r.Post("/subscription", api.CreateSubscription)
	r.Get("/subscription", api.ReadSubscription)
	r.Put("/subscription", api.UpdateSubscription)
	r.Delete("/subscription", api.DeleteSubscription)
	r.Get("/subscriptionList", api.ListSubscription)
	r.Get("/sum_subscriptions", api.SumUserSubscriptions)
	// r.Get("/metrics", promhttp.Handler().ServeHTTP)
	err := http.ListenAndServe(":"+config.CurrentConfig.Server.Port, r)
	if err != nil {
		panic(err)
	}
}
