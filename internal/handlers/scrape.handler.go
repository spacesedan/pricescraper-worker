package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pricescraper-worker/internal/models"
)

func (h *Handler) ScrapePrices(c *gin.Context) {
	var requestBody models.ChunkedCollection

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "failed to parse request body",
		})
	}

	c.AbortWithStatusJSON(200, gin.H{
		"success": "ok",
	})

	h.PriceService.HandleCollections(requestBody)

}
