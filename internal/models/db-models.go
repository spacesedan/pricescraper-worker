package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CollectionDB struct {
	ID              primitive.ObjectID `bson:"_id" json:"_id"`
	Image           string             `bson:"image" json:"image"`
	Contract        string             `bson:"contract" json:"contract"`
	Metadata        string             `bson:"metadata" json:"metadata"`
	Number          int32              `bson:"number" json:"number"`
	Attributes      []interface{}      `bson:"attributes" json:"attributes"`
	RarityScore     float64            `bson:"rarityScore" json:"rarityScore"`
	RarityScoreRank int32              `bson:"rarityScoreRank" json:"rarityScoreRank"`
	OpenSeaLink     string             `bson:"OpenSeaLink" json:"OpenSeaLink"`
	Collection      string             `bson:"collection" json:"collection"`
	AttributeSize   int32              `bson:"attributeSize" json:"attributeSize"`
	Price           float64            `bson:"price" json:"price"`
}

type TypeBreakdown struct {
	Type  string `json:"type" bson:"type"`
	Count int    `json:"count" bson:"count"`
}

type OccurrenceBreakdown struct {
	Type       string `json:"type" bson:"type"`
	Occurrence int    `json:"occurrence" bson:"occurrence"`
}

type Meta struct {
	ID            primitive.ObjectID `bson:"_id"`
	Collection    string             `bson:"collection" `
	Contract      string             `bson:"contract"`
	InProgress    bool               `bson:"inProgress"`
	PriceScraping bool               `bson:"priceScraping"`
}

type FloorPrice struct {
	Price          float64   `bson:"price"`
	PriceEntryTime time.Time `bson:"priceEntryTime"`
}

type TraitCombo struct {
	TraitType string `bson:"trait_type"`
	Value     string `bson:"value"`
}

type TraitCollection struct {
	ID          TraitCombo `bson:"_id"`
	Count       int        `bson:"count"`
	RarityScore float64    `bson:"rarityScore"`
}
