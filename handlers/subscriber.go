package handlers

import (
	"encoding/json"
	"fmt"
	"payment-service/database"
	"time"
)

type OrderMessage struct {
	OrderID string  `json:"order_id"`
	UserID  string  `json:"user_id"`
	Amount  float64 `json:"amount"`
}

func HandleOrderCreated(msgData []byte) error {
	var order OrderMessage
	if err := json.Unmarshal(msgData, &order); err != nil {
		return err
	}

	fmt.Printf("Processing payment for order: %d\n", order.OrderID)

	// Simulate payment success
	status := "success"
	_, err := database.DB.Exec(`
        INSERT INTO payments (order_id, user_id, amount, status, created_at)
        VALUES ($1, $2, $3, $4, $5)`,
		order.OrderID, order.UserID, order.Amount, status, time.Now())

	return err
}
