package dto

import "pricescraper-worker/internal/models"

type CollectionWithTraitCombo struct {
	Collection []models.CollectionDB
	TraitType  string
	TraitValue string
}
