package utils

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func Startup(log *logrus.Logger) {
	_ = godotenv.Load()
	InitRedis()
}
