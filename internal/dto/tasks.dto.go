package dto

import "pricescraper-worker/internal/models"

type ReservoirTask struct {
	PriceMap   models.ReservoirPriceMap
	Collection models.Collection
}

type OpenseaTask struct {
	Stats      models.Stats
	Collection models.Collection
}
