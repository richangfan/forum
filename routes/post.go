package routes

import (
	"encoding/json"
	"net/http"
	"richangfan/forum/model"

	"github.com/gin-gonic/gin"
)

type PostResult struct {
	List List `json:"list"`
}

func AddPostRoute(rg *gin.RouterGroup) {
	group := rg.Group("")

	group.POST("create", func(c *gin.Context) {
		user, err := model.GetUserByToken(c.Query("token"))
		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			return
		}
		raw, err := c.GetRawData()
		if err != nil {
			sendErrorJson(c, err.Error())
			return
		}
		var post model.Post
		err = json.Unmarshal(raw, &post)
		if err != nil {
			sendErrorJson(c, err.Error())
			return
		}
		post.UserId = user.Id
		err = post.AddPost()
		if err != nil {
			sendErrorJson(c, err.Error())
			return
		}
		sendSuccessJson(c, nil)
	})

	group.GET("list", func(c *gin.Context) {
		var post model.Post
		total, err := post.GetTotal()
		if err != nil {
			sendErrorJson(c, err.Error())
			return
		}
		pagination, err := getPagination(c, total)
		if err != nil {
			sendErrorJson(c, err.Error())
			return
		}
		list, err := post.GetList(pagination.Start, pagination.End)
		if err != nil {
			sendErrorJson(c, err.Error())
			return
		}
		sendSuccessJson(c, PostResult{List: List{pagination.PageSize, pagination.Page, total, list}})
	})
}
