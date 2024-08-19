package main

import (
	"log"
	"net/http"
	"tucows-grill-client/internal/api"
	"tucows-grill-client/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	client := api.NewClient()
	handler := handlers.NewHandler(client)
	r := mux.NewRouter()

	r.HandleFunc("/login", handler.LoginHandler).Methods("POST")
	r.HandleFunc("/ingredients/{id}", handler.GetIngredientByIDHandler).Methods("GET")
	r.HandleFunc("/ingredients", handler.PostIngredientHandler).Methods("POST")
	r.HandleFunc("/total-cost", handler.GetTotalCostHandler).Methods("GET")

	log.Println("Starting tucows-grill-client HTTP server on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
