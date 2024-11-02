package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tugas/database"
	"tugas/models"
)

type PaginationResponse struct {
	StatusCode int           `json:"statusCode"`
	Message    string        `json:"message"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalItems int           `json:"total items"`
	TotalPages int           `json:"total pages"`
	Data       []models.Item `json:"data"`
}

type CreatedResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       models.Item `json:"data"`
}

type NotFoundResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// GetItems retrieves all items
func GetItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	rows, err := database.DB.Query("SELECT id, name, price, quantity, category, location FROM items LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []models.Item{}
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Quantity, &item.Category, &item.Location); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	// Count total items for pagination
	var totalItems int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM items").Scan(&totalItems)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := PaginationResponse{
		StatusCode: 200,
		Message:    "Data retrieved successfully",
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: (totalItems + limit - 1) / limit, // Calculate total pages
		Data:       items,
	}

	json.NewEncoder(w).Encode(response)
}

// CreateItem adds a new item
func CreateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert item to database
	err := database.DB.QueryRow("INSERT INTO items (name, price, quantity, category, location) VALUES ($1, $2, $3, $4, $5) RETURNING id", item.Name, item.Price, item.Quantity, item.Category, item.Location).Scan(&item.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := CreatedResponse{
		StatusCode: 201,
		Message:    "Data created successfully",
		Data:       item,
	}

	json.NewEncoder(w).Encode(response)
}

// NotFoundHandler for handling 404 responses
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	response := NotFoundResponse{
		StatusCode: 404,
		Message:    "Data not found",
	}

	json.NewEncoder(w).Encode(response)
}
