package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) addPingRoutes() {
	ping := app.Router
	ping.GET("ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "v1",
		})
	})
}
