package utils

import (
	"pricescraper-worker/internal/dto"
)

// GENERICS ARE COOL!

type SingleTask interface {
	dto.OpenseaTask | dto.ReservoirTask
}

func Producer[S SingleTask](ch chan S, tasks []S) {
	for _, task := range tasks {
		ch <- task
	}
}

func Worker[S SingleTask](chA, chB chan S) {
	task := <-chA
	chB <- task
}
