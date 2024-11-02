package main

import (
	"fmt"
	"log"
	"net/http"
	"tugas/database"
	"tugas/handlers"
)

func main() {
	// Koneksi ke database
	if err := database.Connect(); err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer database.Close()

	// Daftarkan handler
	http.HandleFunc("/items", handlers.GetItems)                        // GET untuk mengambil data
	http.HandleFunc("/items/create", handlers.CreateItem)               // POST untuk membuat data
	http.HandleFunc("/transactions", handlers.GetTransactions)          // GET untuk mengambil transaksi
	http.HandleFunc("/transactions/create", handlers.CreateTransaction) // POST untuk membuat transaksi
	http.HandleFunc("/", handlers.NotFoundHandler)                      // 404 handler

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
