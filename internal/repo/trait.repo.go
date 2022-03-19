package repo

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"pricescraper-worker/internal/apperrors"
	"pricescraper-worker/internal/dto"
	"pricescraper-worker/internal/models"
	"time"
)

type TraitQuery interface {
	GetAllTraits(cn string) ([]*models.TraitCollection, error)
	UpdateTraitPrices(cn string, length int, tokens chan dto.CollectionWithTraitCombo)
}

type traitQuery struct {
}

func (t *traitQuery) GetAllTraits(cn string) ([]*models.TraitCollection, error) {
	m := DB.Collection(cn + traitCollection)

	query := bson.M{}

	var traits []*models.TraitCollection

	results, err := m.Find(ctx, query)
	if err != nil {
		return nil, apperrors.ErrCollectionNotFound
	}

	defer results.Next(ctx)

	for results.Next(ctx) {
		var trait *models.TraitCollection
		results.Decode(&trait)
		traits = append(traits, trait)
	}

	return traits, nil

}

func (t *traitQuery) UpdateTraitPrices(cn string, length int, tokens chan dto.CollectionWithTraitCombo) {
	m := DB.Collection(cn + traitCollection)

	ops := make([]mongo.WriteModel, length)
	for t := 0; t < length; t++ {
		// fmt.Printf("Iteration: %v\n", t)
		trait := <-tokens
		if len(trait.Collection) != 0 {
			// fmt.Printf("Token Trait: %v\n  Trait Value: %v\n Price: %v\n", toke.TraitType, toke.TraitValue, toke.Collection[0].Price)
			for _, token := range trait.Collection {
				if len(trait.Collection) != 0 {
					fp := models.FloorPrice{
						Price:          token.Price,
						PriceEntryTime: time.Now(),
					}

					filter := bson.M{"_id": bson.M{
						"trait_type": trait.TraitType,
						"value":      trait.TraitValue,
					},
					}
					update := bson.M{"$set": bson.M{
						"floorPrice": fp,
					},
					}

					// add the operation to the bulk writer
					ops[t] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)

				} else {

					filter := bson.M{"_id": bson.M{
						"trait_type": trait.TraitType,
						"value":      trait.TraitValue,
					},
					}
					update := bson.M{"$set": bson.M{
						"floorPrice": bson.M{
							"price":          "NA",
							"priceEntryTime": time.Now(),
						},
					},
					}

					thisOp := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
					if thisOp != nil {

						ops = append(ops, thisOp)
					}
				}
			}
		}

	}

	// take in all the operations and remove any that are nil/null
	newOps := make([]mongo.WriteModel, 0, len(ops))
	for _, item := range ops {
		if item != nil {
			newOps = append(newOps, item)
		}
	}

	bulkOps := options.BulkWriteOptions{}
	bulkOps.SetOrdered(false)
	// fmt.Println(len(newOps))
	// run the operation and update the trait floor prices
	if len(newOps) > 0 {
		_, err := m.BulkWrite(ctx, newOps, &bulkOps)
		if err != nil {
			fmt.Print("ERROR: LINE 142")
			log.Fatal(err.Error())
		}
	}
}
