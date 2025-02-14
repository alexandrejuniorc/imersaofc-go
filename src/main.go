package main

import (
	"fmt"
	"net/http"
)

func main() {
	// CHANNEL EXAMPLES
	data := make(chan int)
	workersQuantity := 10

	// Create 10 workers
	for i := 0; i < workersQuantity; i++ {
		go worker(i, data)
	}

	for i := range 100 {
		data <- i
	}

	// HTTP EXAMPLES
	http.HandleFunc("/products", getProducts)
	http.HandleFunc("/cep", searchCEP)

	fmt.Printf("Server running at http://localhost:8080\n")

	http.ListenAndServe(":8080", nil)
}
