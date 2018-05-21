package main

import "github.com/gin-gonic/gin"
import "net/http"
import "strconv"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/dice/roll", diceRollHandler)
	r.Run(":8080")
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
