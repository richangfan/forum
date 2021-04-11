package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"richangfan/forum/middleware"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

type PostList struct {
	List List `json:"list"`
}

func AddPost(rg *gin.RouterGroup) {
	group := rg.Group("")

	group.POST("", func(c *gin.Context) {
		checkLogin(c)
		raw, err := c.GetRawData()
		if err != nil {
			sendErrorJson(c, err.Error())
			return
		}
		var p Post
		err = json.Unmarshal(raw, &p)
		if err != nil {
			sendErrorJson(c, err.Error())
			return
		}
		tlen := utf8.RuneCountInString(p.Title)
		clen := utf8.RuneCountInString(p.Content)
		if tlen == 0 || tlen > 50 || clen == 0 || clen > 10000 {
			sendErrorJson(c, "输入出错")
			return
		}
		ctx := context.Background()
		client := middleware.GetRedisClient()
		val, err := json.Marshal(p)
		if err != nil {
			sendErrorJson(c, err.Error())
			return
		}
		_, err = client.LPush(ctx, CACHE_POST, string(val)).Result()
		if err != nil {
			sendErrorJson(c, err.Error())
			return
		}
		sendSuccessJson(c, nil)
	})

	group.GET("list", func(c *gin.Context) {
		page, err := getPage(c)
		if err != nil {
			sendErrorJson(c, err.Error())
			return
		}
		ctx := context.Background()
		client := middleware.GetRedisClient()
		total, err := client.LLen(ctx, CACHE_POST).Result()
		if err != nil {
			sendErrorJson(c, err.Error())
			return
		}
		var start, end int64
		if total <= 0 {
			sendSuccessJson(c, PostList{List: List{PAGE_SIZE, 1, 0, nil}})
			return
		} else if total <= int64((page-1)*PAGE_SIZE) {
			sendSuccessJson(c, PostList{List: List{PAGE_SIZE, page, total, nil}})
			return
		} else {
			start = int64((page-1)*PAGE_SIZE + 1)
			if total <= int64(page*PAGE_SIZE) {
				end = total
			} else {
				end = int64(page * PAGE_SIZE)
			}
		}
		raw, err := client.LRange(ctx, CACHE_POST, start-1, end-1).Result()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var data = make([]Post, end-start+1)
		for i, v := range raw {
			var p Post
			err = json.Unmarshal([]byte(v), &p)
			if err != nil {
				sendErrorJson(c, err.Error())
				return
			}
			data[i] = p
		}
		sendSuccessJson(c, PostList{List: List{PAGE_SIZE, page, total, data}})
	})
}
