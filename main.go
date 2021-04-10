package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"richangfan/forum/routes"
)

func main() {
	r := gin.Default()

	// 静态资源
	r.Static("/dist", "./dist")

	// 动态资源
	routes.AddIndex(r.Group(""))
	routes.AddPost(r.Group("post"))

	err := r.Run()
	if err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
