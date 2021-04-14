package routes

import (
	"encoding/json"
	"net/http"
	"richangfan/forum/model"

	"github.com/gin-gonic/gin"
)

type UserResult struct {
	User model.User `json:"user"`
}

func AddUserRoute(rg *gin.RouterGroup) {
	group := rg.Group("")

	group.POST("register", func(c *gin.Context) {
		var user model.User
		raw, err := c.GetRawData()
		if err == nil {
			if err = json.Unmarshal(raw, &user); err == nil {
				if err = user.Register(); err == nil {
					sendSuccessJson(c, UserResult{user})
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
					sendSuccessJson(c, UserResult{user})
					return
				}
			}
		}
		sendErrorJson(c, err.Error())
	})

	group.GET("logout", func(c *gin.Context) {
		user, err := model.GetUserByToken(c.Query("token"))
		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			return
		}
		err = user.Logout()
		if err != nil {
			sendErrorJson(c, err.Error())
			return
		}
		sendSuccessJson(c, nil)
	})
}
