package orders

import "time"

type Position struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	AvgPrice   uint64    `json:"averagePrice"`
	Quantity   uint32    `json:"quantity"`
	Direction  string    `json:"direction"`
	OrderID    string    `json:"orderId"`
	ProfitLoss int64     `json:"profitLoss"`
	Symbol     string    `json:"symbol"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// NewPosition is a function that creates a new position.
