package main

import (
	"math/rand"
	"time"
)

// RollDie simulates rolling n dice with x sides.
func RollDie(numDie int, sides int) int {
	// TODO Should seed once at beginning of app
	rand.Seed(time.Now().UTC().UnixNano())
	i := 0
	result := 0
	for i < numDie {
		result = result + (1 + rand.Intn(sides))
		i++
	}
	return result
}
