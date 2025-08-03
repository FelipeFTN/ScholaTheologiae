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
		v1.GET("/read/:book", rh.HandleRead)
		v1.GET("/read/:book/:part", rh.HandleRead)
		v1.GET("/read/:book/:part/:chapter", rh.HandleRead)
		v1.GET("/read/:book/:part/:chapter/:article", rh.HandleRead)

		v1.GET("/summa-theologiae", rh.HandleSummaTheologiae)
		v1.GET("/summa-theologiae/:part", rh.HandleSummaTheologiae)
		v1.GET("/summa-theologiae/:part/:question", rh.HandleSummaTheologiae)
		v1.GET("/summa-theologiae/:part/:question/:article", rh.HandleSummaTheologiae)

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

func (r *RouterHandler) HandleSummaTheologiae(c *gin.Context) {
	summa_theologiae := models.SummaTheologiaeRequest{
		Part:     c.Param("part"),
		Question: c.Param("question"),
		Article:  c.Param("article"),
	}

	summa_theologiae.Validate()

	response, err := r.Controllers.SummaTheologiae(summa_theologiae)
	if err != nil {
		slog.Error("Error in SummaTheologiae", "error", err)
		c.JSON(400, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(200, response)
}

func (r *RouterHandler) HandleRead(c *gin.Context) {
	book_request := models.BookRequest{
		Name:    c.Param("book"),
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
