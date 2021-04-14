package routes

import (
	"encoding/json"
	"richangfan/forum/middleware"
	"richangfan/forum/model"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

type PostResult struct {
	List List `json:"list"`
}

func AddPostRoute(rg *gin.RouterGroup) {
	group := rg.Group("")

	group.POST("create", func(c *gin.Context) {
		user := model.GetUserByToken(c)
		raw, err := c.GetRawData()
		if err == nil {
			var p Post
			err = json.Unmarshal(raw, &p)
			if err == nil {
				tlen := utf8.RuneCountInString(p.Title)
				clen := utf8.RuneCountInString(p.Content)
				if tlen == 0 || tlen > 50 || clen == 0 || clen > 10000 {
					sendErrorJson(c, "输入出错")
					return
				}
				db, err := middleware.GetMysqlClient()
				if err == nil {
					defer db.Close()
					stmt, err := db.Prepare("INSERT INTO post (user_id, title, content, created) VALUES (?, ?, ?, ?)")
					if err == nil {
						defer stmt.Close()
						_, err = stmt.Exec(user.Id, p.Title, p.Content, time.Now().String()[0:19])
						if err == nil {
							sendSuccessJson(c, nil)
						}
					}
				}
			}
		}
		sendErrorJson(c, err.Error())
	})

	group.GET("list", func(c *gin.Context) {
		// page, err := getPage(c)
		// if err == nil {
		// 	db, err := middleware.GetMySQLClient()
		// 	if err == nil {
		// 		defer db.Close()
		// 		stmt, err := db.Prepare("SELECT count(*) AS total FROM post")
		// 		if err == nil {
		// 			defer stmt.Close()
		// 			res, err := stmt.Query()
		// 			if err == nil {
		// 				var start, end int64
		// 				if total == 0 {
		// 					sendSuccessJson(c, PostResult{List: List{PAGE_SIZE, 1, 0, nil}})
		// 					return
		// 				} else if total <= int64((page-1)*PAGE_SIZE) {
		// 					sendSuccessJson(c, PostResult{List: List{PAGE_SIZE, page, total, nil}})
		// 					return
		// 				} else {
		// 					start = int64((page-1)*PAGE_SIZE + 1)
		// 					if total <= int64(page*PAGE_SIZE) {
		// 						end = total
		// 					} else {
		// 						end = int64(page * PAGE_SIZE)
		// 					}
		// 				}
		// 				raw, err := client.LRange(ctx, CACHE_POST, start-1, end-1).Result()
		// 				if err != nil {
		// 					fmt.Println(err.Error())
		// 					return
		// 				}
		// 				var data = make([]Post, end-start+1)
		// 				for i, v := range raw {
		// 					var p Post
		// 					err = json.Unmarshal([]byte(v), &p)
		// 					if err != nil {
		// 						sendErrorJson(c, err.Error())
		// 						return
		// 					}
		// 					data[i] = p
		// 				}
		// 				sendSuccessJson(c, PostResult{List: List{PAGE_SIZE, page, total, data}})
		// 			}
		// 		}
		// 	}
		// }
		// sendErrorJson(c, err.Error())
		// return
	})
}
