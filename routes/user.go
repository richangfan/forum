package routes

import (
	"encoding/json"
	"richangfan/forum/model"

	"github.com/gin-gonic/gin"
)

func AddUserRoute(rg *gin.RouterGroup) {
	group := rg.Group("")

	group.POST("register", func(c *gin.Context) {
		var user model.User
		raw, err := c.GetRawData()
		if err == nil {
			if err = json.Unmarshal(raw, &user); err == nil {
				if err = user.Register(); err == nil {
					sendSuccessJson(c, user)
					return
				}
			}
		}
		sendErrorJson(c, err.Error())
	})

	group.POST("login", func(c *gin.Context) {
		var user model.User
		raw, err := c.GetRawData()
		if err == nil {
			if err = json.Unmarshal(raw, &user); err == nil {
				if err = user.Login(); err == nil {
					sendSuccessJson(c, user)
					return
				}
			}
		}
		sendErrorJson(c, err.Error())
	})

	group.GET("logout", func(c *gin.Context) {
	})
}
