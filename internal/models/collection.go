package models

type Collection struct {
	OSSlug          string `bson:"slug"`
	ContractAddress string `bson:"contract"`
	Collection      string `bson:"collection"`
}

type ChunkedCollection struct {
	Collections []Collection
	ID          int
}
