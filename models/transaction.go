package models

import "time"

type Transaction struct {
	ID          int       `json:"id"`
	ItemID      int       `json:"item_id"`
	Quantity    int       `json:"quantity"`
	Type        string    `json:"type"` // "in" for incoming, "out" for outgoing
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
}
