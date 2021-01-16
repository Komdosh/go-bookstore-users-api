package app

import (
	"github.com/Komdosh/go-bookstore-utils/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApp() {
	mapUrls()

	logger.Info("about to start the application...")
	router.Run(":8081")
}
