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

	go func() {
		if err := pubsub.ListenForOrders(); err != nil {
			log.WithError(err).Error("Pub/Sub listen error")
		}
	}()

	utils.Startup(log)
	if err := database.Connect(); err != nil {
		log.WithError(err).Fatal("DB connection failed")
	}

	app := iris.New()
	SetupRouter(app)

	go func() {
		if err := app.Listen(":3001", iris.WithoutInterruptHandler); err != nil {
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
