package handler

import (
	"github.com/gin-gonic/gin"

	"scholatheologiae-api/controller"
	"scholatheologiae-api/models"
)

func Server() {
	server := gin.Default()

	// V1
	v1 := server.Group("v1")
	{
		v1.GET("/health", HandleHealth)
		v1.GET("/read/:book", HandleRead)
		v1.GET("/read/:book/:chapter", HandleRead)
		v1.GET("/read/:book/:chapter/:verse", HandleRead)

		v1.GET("/summa_theologiae", HandleSummaTheologiae)
		v1.GET("/summa_theologiae/:part", HandleSummaTheologiae)
		v1.GET("/summa_theologiae/:part/:question", HandleSummaTheologiae)
		v1.GET("/summa_theologiae/:part/:question/:article", HandleSummaTheologiae)
	}

	server.Run(":8080")
}

func HandleSummaTheologiae(c *gin.Context) {
	summa_theologiae := models.SummaTheologiaeRequest{
		Part:     c.Param("part"),
		Question: c.Param("question"),
		Article:  c.Param("article"),
	}

	summa_theologiae.Validate()

	response, err := controller.SummaTheologiae(summa_theologiae)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, response)
}

func HandleRead(c *gin.Context) {
	book := c.Param("book")
	chapter := c.Param("chapter")
	verse := c.Param("verse")

	// Here you would typically fetch the requested data from a database or another source
	// For demonstration purposes, we'll just return the parameters
	c.JSON(200, gin.H{
		"book":    book,
		"chapter": chapter,
		"verse":   verse,
	})
}

func HandleHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
