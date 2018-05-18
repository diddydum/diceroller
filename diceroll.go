package main

import (
    "fmt"
    "math/rand"
    "time"
)


func rollDie(numDie int, sides int, player string) int {
	rand.Seed(time.Now().UTC().UnixNano())
	i := 0
	result := 0
	for i < numDie {
		result = result + (1 + rand.Intn(sides))
		i++
	}
    fmt.Printf("%s rolls a %dd%d for %d\n", player, numDie, sides, result)

    dice := fmt.Sprintf("%dd%d", numDie, sides)
    q := fmt.Sprintf("INSERT INTO user_tbl (name, dice, result) VALUES (%s, %s, %d);", player, dice, result)
    fmt.Println(dice)
    fmt.Println(q)
	return result // or whatever, query, 2d6... 
}


func main() {
    rollDie(2, 12, "nate")
    rollDie(2, 12, "nate")
    rollDie(2, 12, "nate")
    rollDie(2, 12, "nate")
    rollDie(2, 12, "nate")
    rollDie(2, 12, "nate")
    rollDie(2, 12, "nate")
    rollDie(2, 12, "nate")
    rollDie(2, 12, "nate")
    rollDie(2, 12, "nate")
    rollDie(2, 12, "nate")
}

