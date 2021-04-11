package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Result map[string]interface{}

type Post struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type List struct {
	PageSize int         `json:"pageSize"`
	Page     int         `json:"page"`
	Total    int64       `json:"total"`
	Data     interface{} `json:"data"`
}

// 默认页面列表长度
const PAGE_SIZE = 10

// 帖子列表在缓存中的key
const CACHE_POST = "post"

func sendErrorJson(c *gin.Context, m string) {
	c.JSON(http.StatusOK, Result{
		"error":   1,
		"message": m,
		"data":    nil,
	})
}

func sendSuccessJson(c *gin.Context, d interface{}) {
	c.JSON(http.StatusOK, Result{
		"error":   0,
		"message": "",
		"data":    d,
	})
}

func getPage(c *gin.Context) (int, error) {
	raw := c.Query("page")
	if raw == "" {
		return 1, nil
	}
	page, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}
	if page <= 1 {
		page = 1
	}
	return page, nil
}

func checkLogin(c *gin.Context) {
	raw := c.Query("token")
	if raw == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
