package main

import (
	"net/http"

	"github.com/DaniilKalts/market-rest-api/logger"
)

func main() {
	logger.Info("Market REST API is launched!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to my Marketplace!"))
	})

	logger.Info("Server is running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Fatal("Failed to run the server: %v")
	}
}
