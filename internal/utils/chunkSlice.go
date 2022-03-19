package utils

import "go.mongodb.org/mongo-driver/mongo"

func ChunkCollections(slice []mongo.WriteModel, chunkSize int) [][]mongo.WriteModel {
	var chunks [][]mongo.WriteModel
	length := len(slice)

	for i := 0; i < length; i += chunkSize {
		end := i + chunkSize

		if end > length {
			end = length
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
