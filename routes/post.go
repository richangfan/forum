package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddPost(rg *gin.RouterGroup) {
	group := rg.Group("")

	group.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "post list")
	})

	group.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "post list")
	})
}
