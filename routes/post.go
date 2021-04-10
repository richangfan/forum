package routes

import (
	"encoding/json"
	"richangfan/forum/tool"

	"github.com/gin-gonic/gin"
)

type Post struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func AddPost(rg *gin.RouterGroup) {
	group := rg.Group("")

	group.POST("", func(c *gin.Context) {
		raw, err := c.GetRawData()
		if err != nil {
			tool.SendErrorJson(c, err.Error())
			return
		}
		var p Post
		err = json.Unmarshal(raw, &p)
		if err != nil {
			tool.SendErrorJson(c, err.Error())
			return
		}
		tool.SendSuccessJson(c, nil)
	})
}
