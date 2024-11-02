package handlers

import (
	"encoding/json"
	"net/http"
	"tugas/database"
	"tugas/models"
)

// CreateTransaction records a new transaction
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert transaction to database
	_, err := database.DB.Exec("INSERT INTO transactions (item_id, quantity, type, timestamp, description) VALUES ($1, $2, $3, $4, $5)",
		transaction.ItemID, transaction.Quantity, transaction.Type, transaction.Timestamp, transaction.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := CreatedResponse{
		StatusCode: 201,
		Message:    "Transaction created successfully",
	}

	json.NewEncoder(w).Encode(response)
}

// GetTransactions retrieves all transactions
func GetTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := database.DB.Query("SELECT id, item_id, quantity, type, timestamp, description FROM transactions")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	transactions := []models.Transaction{}
	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.ItemID, &transaction.Quantity, &transaction.Type, &transaction.Timestamp, &transaction.Description); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, transaction)
	}

	json.NewEncoder(w).Encode(transactions)
}
