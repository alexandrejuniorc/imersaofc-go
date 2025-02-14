package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
}

var products = []Product{
	{ID: 1, Name: "Laptop", Description: "Macbook Pro", Price: 2000, Category: "Electronics"},
	{ID: 2, Name: "Smartphone", Description: "iPhone 12", Price: 1000, Category: "Electronics"},
	{ID: 3, Name: "Headphones", Description: "AirPods Pro", Price: 250, Category: "Electronics"},
	{ID: 4, Name: "T-Shirt", Description: "White T-Shirt", Price: 20, Category: "Clothing"},
	{ID: 5, Name: "Jeans", Description: "Blue Jeans", Price: 50, Category: "Clothing"},
	{ID: 6, Name: "Sneakers", Description: "White Sneakers", Price: 80, Category: "Clothing"},
}

func getProducts(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(products)
}

func counter(count int) {
	for i := 0; i < count; i++ {
		fmt.Println(i)
		time.Sleep(time.Second)
	}
}

// Thread 1 (goroutine)
func main() {
	channel := make(chan string)

	// Thread 2 (goroutine)
	go func() {
		channel <- "Hello from goroutine" // Channel is full until a message is received
	}()

	// The main goroutine will wait for the other goroutine to finish
	message := <-channel // Channel is empty until a message is received
	fmt.Println(message) // Hello from goroutine
}
