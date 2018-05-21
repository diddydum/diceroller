package main

import (
	"log"
	"net/http"
	"strconv"

	"database/sql"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func main() {
	// Attempt to connect to sql db, and don't really do anything after that.
	connStr := "postgres://diceroller:diceroller@localhost/diceroller?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error while connecting to postgres: ", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Error while testing connection to db: ", err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/dice/roll", diceRollHandler)
	err = r.Run(":8080")
	if err != nil {
		log.Fatal("Error while running api server: ", err)
	}
}

func diceRollHandler(c *gin.Context) {
	numStr := c.DefaultQuery("num", "1")
	sidesStr := c.Query("sides")

	num, err := strconv.Atoi(numStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad value for query parameter 'num', expected uint",
		})
		return
	}
	sides, err := strconv.Atoi(sidesStr)
	if err != nil || sides == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad value for query paramater 'sides', expected uint and > 0",
		})
		return
	}

	// Done validating, proceed
	roll := RollDie(num, sides)
	c.JSON(http.StatusOK, gin.H{
		"result": roll,
	})
}
