package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	health = "alive"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, health)
}
