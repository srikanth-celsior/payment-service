package main

import (
	"context"
	"os"
	"os/signal"
	"payment-service/database"
	"payment-service/pubsub"
	"payment-service/utils"
	"syscall"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	utils.Startup(log)

	if err := database.Connect(); err != nil {
		log.WithError(err).Fatal("DB connection failed")
	}
	go func() {
		if err := pubsub.ListenForOrders(); err != nil {
			log.WithError(err).Error("Pub/Sub listen error")
		}
	}()

	app := iris.New()
	SetupRouter(app)

	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		if err := app.Listen(":" + port); err != nil {
			log.WithError(err).Fatal("Failed to start HTTP server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Server forced to shutdown")
	}
	log.Info("Server exited cleanly")
}
