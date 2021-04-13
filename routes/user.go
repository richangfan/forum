package routes

import (
	"encoding/json"
	"richangfan/forum/middleware"

	"github.com/gin-gonic/gin"
)

func AddUserRoute(rg *gin.RouterGroup) {
	group := rg.Group("")

	group.POST("register", func(c *gin.Context) {
		var user middleware.User
		raw, err := c.GetRawData()
		if err == nil {
			if err = json.Unmarshal(raw, &user); err == nil {
				if err = middleware.Register(&user); err == nil {
					if token, err := middleware.Login(user); err == nil {
						sendSuccessJson(c, user)
						return
					}
				}
			}
		}
		sendErrorJson(c, err.Error())
	})

	group.POST("login", func(c *gin.Context) {
	})

	group.GET("logout", func(c *gin.Context) {
	})
}
