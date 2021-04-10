package tool

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result map[string]interface{}

func SendErrorJson(c *gin.Context, m string) {
	c.JSON(http.StatusOK, Result{
		"error":   1,
		"message": m,
		"data":    nil,
	})
}

func SendSuccessJson(c *gin.Context, d interface{}) {
	c.JSON(http.StatusOK, Result{
		"error":   0,
		"message": "",
		"data":    d,
	})
}
