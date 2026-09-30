package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"

	"scholatheologiae-api/controllers"
	"scholatheologiae-api/models"
)

func Run(c *controllers.Controllers) {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	rh := NewRouterHandler(c)

	// V1
	v1 := server.Group("v1")
	{
		v1.GET("/health", rh.HandleHealth)
		v1.GET("/books/:name", rh.HandleBooks)
		v1.GET("/books/:name/:part", rh.HandleBooks)
		v1.GET("/books/:name/:part/:chapter", rh.HandleBooks)
		// Reserved for article-level reads; answers 404 while no book ships articles.
		v1.GET("/books/:name/:part/:chapter/:article", rh.HandleBooks)

		v1.GET("/search", rh.HandleSearch)
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
	bookRequest := models.BookRequest{
		Name:    c.Param("name"),
		Part:    c.Param("part"),
		Chapter: c.Param("chapter"),
		Article: c.Param("article"),
	}
	bookRequest.Validate()

	response, err := r.Controllers.Read(bookRequest)
	if err != nil {
		r.fail(c, err)
		return
	}

	if response == nil {
		r.fail(c, fmt.Errorf("%w: no content for this request", models.ErrNotFound))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (r *RouterHandler) HandleSearch(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		r.fail(c, fmt.Errorf("%w: query parameter 'q' is required", models.ErrBadRequest))
		return
	}

	response, err := r.Controllers.Search(query)
	if err != nil {
		r.fail(c, err)
		return
	}

	// An empty result set is a valid answer to a valid query, not an error.
	if response == nil {
		response = []models.SearchResult{}
	}

	c.JSON(http.StatusOK, response)
}

func (r *RouterHandler) HandleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// fail logs the error and answers with the status code the error maps to, so that
// a missing book is not reported as a malformed request and a real failure is not
// reported as "no results".
func (r *RouterHandler) fail(c *gin.Context, err error) {
	slog.Error("Request failed", "path", c.Request.URL.Path, "error", err)

	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, models.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, models.ErrBadRequest):
		status = http.StatusBadRequest
	}

	c.JSON(status, gin.H{
		"status": false,
		"error":  err.Error(),
	})
}
