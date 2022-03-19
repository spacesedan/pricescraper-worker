package apperrors

import "errors"

var (
	ErrCollectionNotFound = errors.New("collection not found")
	ErrDecode             = errors.New("could not decode")
)
