package main

import (
	"log"
	"payment-service/database"
	"payment-service/pubsub"

	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
)

func main() {
	// Load .env variables
	_ = godotenv.Load()

	// Connect to DB
	if err := database.Connect(); err != nil {
		log.Fatal("DB connection failed:", err)
	}

	// Start listening to Pub/Sub in a goroutine (non-blocking)
	go func() {
		if err := pubsub.ListenForOrders(); err != nil {
			log.Fatalf("Pub/Sub listen error: %v", err)
		}
	}()

	// Start HTTP server for health check or future use
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.WriteString("Payment Service is running")
	})
	if err := app.Listen(":3001"); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
