package server

import (
	"log/slog"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"

	"scholatheologiae-api/controller"
	"scholatheologiae-api/models"
)

func Run(c *controller.Controller) {
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
		v1.GET("/read/:book/:chapter", rh.HandleRead)
		v1.GET("/read/:book/:chapter/:verse", rh.HandleRead)

		v1.GET("/summa-theologiae", rh.HandleSummaTheologiae)
		v1.GET("/summa-theologiae/:part", rh.HandleSummaTheologiae)
		v1.GET("/summa-theologiae/:part/:question", rh.HandleSummaTheologiae)
		v1.GET("/summa-theologiae/:part/:question/:article", rh.HandleSummaTheologiae)
	}

	// Graceful shutdown
	endless.ListenAndServe(":8080", server)
}

type RouterHandler struct {
	Controller *controller.Controller
}

func NewRouterHandler(c *controller.Controller) *RouterHandler {
	return &RouterHandler{
		Controller: c,
	}
}

func (r *RouterHandler) HandleSummaTheologiae(c *gin.Context) {
	summa_theologiae := models.SummaTheologiaeRequest{
		Part:     c.Param("part"),
		Question: c.Param("question"),
		Article:  c.Param("article"),
	}

	summa_theologiae.Validate()

	response, err := r.Controller.SummaTheologiae(summa_theologiae)
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
	book := c.Param("book")
	chapter := c.Param("chapter")
	verse := c.Param("verse")

	c.JSON(200, gin.H{
		"book":    book,
		"chapter": chapter,
		"verse":   verse,
	})
}

func (r *RouterHandler) HandleHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
