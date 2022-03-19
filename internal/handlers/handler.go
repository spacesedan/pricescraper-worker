package handlers

import (
	"github.com/gin-gonic/gin"
	"pricescraper-worker/internal/service"
)

type Handler struct {
	PriceService service.PriceService
}

type Config struct {
	Router       *gin.Engine
	PriceService service.PriceService
}

func NewHandler(c *Config) {
	h := &Handler{
		PriceService: c.PriceService,
	}

	api := c.Router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/", h.Health)
			v1.POST("/scrape", h.ScrapePrices)

		}
	}
}
