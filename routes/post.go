package routes

import (
	"encoding/json"
	"richangfan/forum/tool"
	"unicode/utf8"

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
		tlen := utf8.RuneCountInString(p.Title)
		clen := utf8.RuneCountInString(p.Content)
		if tlen == 0 || tlen > 50 || clen == 0 || clen > 10000 {
			tool.SendErrorJson(c, "输入出错")
			return
		}
		tool.SendSuccessJson(c, nil)
	})
}
