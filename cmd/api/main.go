package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"first-api/db"
	"first-api/internal/handler"
	"first-api/internal/repository"
	"first-api/internal/routes"
	"first-api/internal/usecases"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar as variáveis de ambiente: %s", err)
	}

	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	CustomerRepository := repository.NewCustomerRepository(dbConnection)
	CustomerUseCase := usecases.NewCustomerUseCase(CustomerRepository)
	CustomerHandler := handler.NewCustomerHandler(CustomerUseCase)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	routes.CustomerRoutes(r, CustomerHandler)

	log.Println("API rodando em http://localhost:8080")
	log.Println("GET    /client      -> listar clientes")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %s", err)
	}

}
