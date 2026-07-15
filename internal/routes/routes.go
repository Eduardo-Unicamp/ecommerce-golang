package routes

import (
	"first-api/internal/auth"
	"first-api/internal/handler"
	"first-api/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r *chi.Mux, handler *handler.AuthHandler) {
	r.Post("/register", handler.Register)
	r.Post("/login", handler.Login)
	r.Post("/refresh", handler.Refresh)

	//OAUTH2
	r.Get("/auth/{provider}", handler.BeginOAuth)
	r.Get("/auth/{provider}/callback", handler.CallbackOAuth)
}

func CustomerRoutes(r *chi.Mux, handler *handler.CustomerHandler, jwtConfig *auth.JWTConfig) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(jwtConfig))

		r.Get("/customer", handler.GetCustomers)
		r.Get("/customer/{customer_id}", handler.GetCustomerByID)
		r.Post("/customer", handler.CreateCustomer)
		r.Put("/customer/{customer_id}", handler.UpdateCustomer)
		r.Delete("/customer/{customer_id}", handler.DeleteCustomer)
	})

}

func ProductRoutes(r *chi.Mux, handler *handler.ProductHandler, jwtConfig *auth.JWTConfig) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(jwtConfig))
		r.Get("/product", handler.GetProducts)
		r.Get("/product/{product_id}", handler.GetProductByID)
		r.Post("/product", handler.CreateProduct)
		r.Put("/product/{product_id}", handler.UpdateProduct)
		r.Delete("/product/{product_id}", handler.DeleteProduct)
	})

}

func OrderRoutes(r *chi.Mux, handler *handler.OrderHandler, jwtConfig *auth.JWTConfig) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(jwtConfig))

		r.Get("/order/all/{customer_id}", handler.GetOrders)
		r.Get("/order/{order_id}", handler.GetOrderByID)
		r.Post("/order", handler.CreateOrder)
		r.Post("/order/cancel/{order_id}", handler.CancelOrder)
		r.Post("/order/pay/{order_id}", handler.PayOrder)
	})
}
