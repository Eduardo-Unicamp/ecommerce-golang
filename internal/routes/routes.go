package routes

import (
	"first-api/internal/handler"

	"github.com/go-chi/chi/v5"
)

func CustomerRoutes(r *chi.Mux, handler *handler.CustomerHandler) {
	r.Get("/client", handler.GetCustomers)
	r.Post("/client", handler.CreateCustomer)
	r.Put("/client/{customerId}", handler.UpdateCustomer)
	r.Delete("/client/{customerId}", handler.DeleteCustomer)

}
