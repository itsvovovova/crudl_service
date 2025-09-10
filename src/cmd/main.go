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
	"crudl_service/src/api"
	"crudl_service/src/config"
	"net/http"

	"log"

	_ "crudl_service/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
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

	log.Println("Registering Swagger documentation")
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	log.Println("Starting HTTP server", "port", config.CurrentConfig.Server.Port)
	err := http.ListenAndServe(":"+config.CurrentConfig.Server.Port, r)
	if err != nil {
		log.Fatal("Failed to start HTTP server")
	}
}
