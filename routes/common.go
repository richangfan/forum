package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Result map[string]interface{}

type List struct {
	PageSize int         `json:"pageSize"`
	Page     int         `json:"page"`
	Total    int64       `json:"total"`
	Data     interface{} `json:"data"`
}

type Pagination struct {
	Page   int
	Limit  int
	Offset int
}

// 默认页面列表长度
const PAGE_SIZE = 10

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

func getPagination(c *gin.Context, total int64) (Pagination, error) {
	var p Pagination
	raw := c.Query("pageSize")
	if raw == "" {
		p.Limit = PAGE_SIZE
	} else {
		pageSize, err := strconv.Atoi(raw)
		if err != nil {
			return p, errors.New("获取页长错误")
		}
		if pageSize <= 0 || pageSize > 100 {
			p.Limit = PAGE_SIZE
		} else {
			p.Limit = pageSize
		}
	}
	raw = c.Query("page")
	if raw == "" {
		p.Page = 1
	} else {
		page, err := strconv.Atoi(raw)
		if err != nil {
			return p, err
		}
		if page < 1 {
			p.Page = 1
		} else {
			p.Page = page
		}
	}
	p.Offset = (p.Page - 1) * p.Limit
	return p, nil
}
