package repo

import (
	"go.mongodb.org/mongo-driver/bson"
	"pricescraper-worker/internal/apperrors"
	"pricescraper-worker/internal/models"
	"time"
)

type MetaQuery interface {
	FindByName(cn string) (*models.Meta, error)
	UpdateFloorPrice(cn string, fp float64)
}

type metaQuery struct {
}

func (m *metaQuery) FindByName(cn string) (*models.Meta, error) {
	c := DB.Collection(metaCollection)

	var collection *models.Meta

	filter := bson.M{
		"collection": cn,
	}

	result := c.FindOne(ctx, filter)
	if result == nil {
		return nil, apperrors.ErrCollectionNotFound
	}

	err := result.Decode(&models.Meta{})
	if err != nil {
		return nil, apperrors.ErrDecode
	}

	return collection, nil

}

func (m *metaQuery) UpdateFloorPrice(cn string, fp float64) {
	c := DB.Collection(metaCollection)

	filter := bson.M{
		"collection": cn,
	}

	floorPrice := models.FloorPrice{
		Price:          fp,
		PriceEntryTime: time.Now(),
	}

	update := bson.M{
		"$set": bson.M{
			"floorPrice": floorPrice,
		},
	}

	c.FindOneAndUpdate(ctx, filter, update)
}
