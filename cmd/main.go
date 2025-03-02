package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vicjeremy/blog-api/api"
)

func main() {
	api.InitDB()
	r := gin.Default()

	//routes
	r.POST("/posts", api.CreateBlog)
	r.PUT("/posts/:id", api.UpdateBlog)
	r.DELETE("/posts/:id", api.DeleteBlog)
	r.GET("/posts/:id", api.GetBlog)
	r.GET("/posts", api.GetBlogs)

	r.Run(":8080")
}
