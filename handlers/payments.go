package handlers

import (
	"net/http"
	"payment-service/database"
	"payment-service/models"

	"github.com/kataras/iris/v12"
)

// GetAllPayments handles GET /payments
func GetAllPayments(ctx iris.Context) {
	rows, err := database.DB.Query("SELECT id, order_id, user_id, amount, status, created_at FROM payments")
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to fetch payments"})
		return
	}
	defer rows.Close()

	payments := []models.Payment{}
	for rows.Next() {
		var p models.Payment
		err := rows.Scan(&p.ID, &p.OrderID, &p.UserID, &p.Amount, &p.Status, &p.CreatedAt)
		if err != nil {
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Failed to parse payment"})
			return
		}
		payments = append(payments, p)
	}

	ctx.JSON(payments)
}
