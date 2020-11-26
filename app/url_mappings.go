package app

import (
	"github.com/Komdosh/go-bookstore-users-api/controllers/health"
	"github.com/Komdosh/go-bookstore-users-api/controllers/users"
)

func mapUrls() {
	router.GET("/health", health.Health)

	router.GET("/users/:user_id", users.GetUser)
	router.POST("/users", users.CreateUser)
}
