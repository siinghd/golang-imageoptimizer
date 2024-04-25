package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/siinghd/golang-imageoptimizer/internal/handler"
)



func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", handler.ProcessImageHandler).Methods("GET")
	router.HandleFunc("/", handler.NotFoundHandler).Methods("POST", "PUT", "DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3010"
	}

	fmt.Printf("Server running on port %s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}








