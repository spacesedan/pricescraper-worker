package repo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"pricescraper-worker/internal/dto"
	"pricescraper-worker/internal/models"
	"pricescraper-worker/internal/utils"
	"strconv"
	"time"
)

type CollectionQuery interface {
	FindLowestFloor(cn string) float64
	StorePriceMap(priceMap models.ReservoirPriceMap, collection string)
	FindWithTraitCombo(cn string, ttype string, tvalue string, jobs <-chan int, t chan<- dto.CollectionWithTraitCombo)
	RemoveOldPrices(cn string)
}

type collectionQuery struct {
}

func (c *collectionQuery) FindLowestFloor(cn string) float64 {
	m := DB.Collection(cn)

	lowest := options.FindOne()
	lowest.SetSort(bson.M{"price": 1})
	query := bson.M{"price": bson.M{"$exists": true}}

	var token models.CollectionDB
	result := m.FindOne(ctx, query, lowest)
	err := result.Decode(&token)
	if err != nil {
		//log.Println(err)
	}

	return token.Price
}

func (c *collectionQuery) StorePriceMap(priceMap models.ReservoirPriceMap, collection string) {
	tn := time.Now()
	var bulkOps []mongo2.WriteModel
	for k, v := range priceMap.Tokens {
		tId, _ := strconv.Atoi(k)
		filter := bson.M{"number": tId}
		update := bson.D{
			{"$set", bson.M{"price": v}},
			{"$set", bson.M{"priceEntryTime": tn}},
		}
		op := mongo2.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
		bulkOps = append(bulkOps, op)
	}

	m := DB.Collection(collection)
	ctx := context.Background()

	chunked := utils.ChunkCollections(bulkOps, 1000)

	for _, chunk := range chunked {
		res, err := m.BulkWrite(ctx, chunk)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Bulk Write Update:", res.ModifiedCount, collection)
	}

}

func (c *collectionQuery) FindWithTraitCombo(cn string, ttype string, tvalue string, jobs <-chan int, t chan<- dto.CollectionWithTraitCombo) {
	m := DB.Collection(cn)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	var tokens []models.CollectionDB

	for _ = range jobs {
		smallest := options.Find()
		smallest.SetSort(bson.D{{"price", 1}}).SetLimit(1)
		// can potentially speed this up using go routines and channels
		query := bson.M{
			"$and": []bson.M{
				{"price": bson.M{"$exists": true}},
				{
					"attributes": bson.M{
						"$elemMatch": bson.M{"trait_type": ttype, "value": tvalue},
					},
				},
			}}
		fr, err := m.Find(ctx, query, smallest)
		if err != nil {

			log.Fatal(err.Error())
		}
		for fr.Next(ctx) {
			var token models.CollectionDB
			err := fr.Decode(&token)
			if err != nil {
				fmt.Println("HERE?")
				log.Fatal(err)
			}
			tokens = append(tokens, token)
		}

		// fmt.Printf("Tokens found %v\n", len(tokens))
		// create an object that will be used to build out the bulk operation
		var d = dto.CollectionWithTraitCombo{
			Collection: tokens,
			TraitType:  ttype,
			TraitValue: tvalue,
		}

		// send the object to the token channel so that it can start building appending jobs to the
		// bulk operation
		t <- d

	}

}

func (c *collectionQuery) RemoveOldPrices(cn string) {
	m := DB.Collection(cn)

	query := bson.M{
		"price": bson.M{
			"$exists": true,
		},
	}

	update := bson.D{
		{"$unset", bson.M{"price": "NA"}},
		{"$unset", bson.M{"priceEntryTime": "NA"}},
	}

	_, err := m.UpdateMany(ctx, query, update)
	if err != nil {
		fmt.Println(err)
	}
}
