package internal

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"pricescraper-worker/internal/handlers"
	"pricescraper-worker/internal/repo"
	"pricescraper-worker/internal/service"
)

func Inject(db *mongo.Database) (*gin.Engine, error) {
	dao := repo.NewDAO(db)

	priceService := service.NewPriceService(dao)

	priceService.UpdateTraitFloorPrices("little-lemon-friends")

	app := gin.Default()

	handlers.NewHandler(&handlers.Config{
		Router:       app,
		PriceService: priceService,
	})

	return app, nil
}
