package server

import (
	"log/slog"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"

	"scholatheologiae-api/controllers"
	"scholatheologiae-api/models"
)

func Run(c *controllers.Controllers) {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	server.Use(gin.Recovery())
	rh := NewRouterHandler(c)

	// V1
	v1 := server.Group("v1")
	{
		v1.GET("/health", rh.HandleHealth)
		v1.GET("/books/:name", rh.HandleBooks)
		v1.GET("/books/:name/:part", rh.HandleBooks)
		v1.GET("/books/:name/:part/:chapter", rh.HandleBooks)
		v1.GET("/books/:name/:part/:chapter/:article", rh.HandleBooks)

		v1.POST("/search", rh.HandleSearch)
	}

	// Graceful shutdown
	endless.ListenAndServe(":8080", server)
}

type RouterHandler struct {
	Controllers *controllers.Controllers
}

func NewRouterHandler(c *controllers.Controllers) *RouterHandler {
	return &RouterHandler{
		Controllers: c,
	}
}

func (r *RouterHandler) HandleBooks(c *gin.Context) {
	book_request := models.BookRequest{
		Name:    c.Param("name"),
		Part:    c.Param("part"),
		Chapter: c.Param("chapter"),
		Article: c.Param("article"),
	}
	book_request.Validate()

	response, err := r.Controllers.Read(book_request)
	if err != nil {
		slog.Error("Error in Read", "error", err)
		c.JSON(400, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(200, response)
}

func (r *RouterHandler) HandleSearch(c *gin.Context) {
	c.JSON(204, gin.H{})
}

func (r *RouterHandler) HandleHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
