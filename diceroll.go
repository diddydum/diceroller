package main

import (
    "fmt"
    "math/rand"
    "time"
)


func RollDie(num int, max int, player string) int {
	rand.Seed(time.Now().UTC().UnixNano())
	i := 0
	result := 0
	for i < num {
		result = result + (1 + rand.Intn(max))
		i++
	}
    fmt.Printf("%s rolls a %dd%d for %d\n", player, num, max, result)

    dice := fmt.Sprintf("%dd%d", num, max)
    q := fmt.Sprintf("INSERT INTO user_tbl (name, dice, result) VALUES (%s, %s, %d)", player, dice, result)
    fmt.Println(dice)
    fmt.Println(q)
	return result // or whatever, query, 2d6... 

}



func main() {
    RollDie(2, 12, "nate")
    RollDie(2, 12, "nate")
    RollDie(2, 12, "nate")
    RollDie(2, 12, "nate")
    RollDie(2, 12, "nate")
    RollDie(2, 12, "nate")
    RollDie(2, 12, "nate")
    RollDie(2, 12, "nate")
    RollDie(2, 12, "nate")
    RollDie(2, 12, "nate")
    RollDie(2, 12, "nate")
}

