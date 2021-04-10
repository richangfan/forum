package routes

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddIndex(rg *gin.RouterGroup) {
	group := rg.Group("")

	group.GET("", func(c *gin.Context) {
		html, err := ioutil.ReadFile("dist/pc.html")
		if err != nil {
			c.String(http.StatusNotFound, "")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", html)
	})
}
