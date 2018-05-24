package main

import (
	"bytes"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"database/sql"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func main() {
	// Seed according to time
	rand.Seed(time.Now().UnixNano())

	connStr := "postgres://diceroller:diceroller@localhost/diceroller?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error while connecting to postgres: ", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Error while testing connection to db: ", err)
	}

	server := Server{db: db}

	r := gin.Default()
	r.GET("/ping", server.pingHandler)
	r.POST("/dice/roll", server.diceRollHandler)
	r.POST("/rooms", server.addRoomHandler)
	r.POST("/users", server.addUserHandler)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal("Error while running api server: ", err)
	}
}

type Server struct {
	db *sql.DB
}

func (s *Server) pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (s *Server) diceRollHandler(c *gin.Context) {
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

func (s *Server) addRoomHandler(c *gin.Context) {
	name := c.Query("name")
	if len(name) > 64 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Name too long, must be less than 64 bytes in length",
		})
		return
	}

	// Need to choose a name
	id := generateToken()
	_, err := s.db.Exec("INSERT INTO rooms (id, name) VALUES ($1, $2)", id, name)
	if err != nil {
		log.Printf("Error while creating room: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"room_id": id,
	})
}

func (s *Server) addUserHandler(c *gin.Context) {
	name := strings.TrimSpace(c.Query("name"))
	if len(name) > 32 || len(name) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad length for name, constraint must hold: (33 > name > 1)",
		})
		return
	}

	token := generateToken()
	result, err := s.db.Exec("INSERT INTO users (token, name) VALUES ($1, $2)", token, name)
	if err != nil {
		log.Printf("Error while creating user: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Unable to extract last insert id for user: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    id,
		"name":  name,
		"token": token,
	})
}

func generateToken() string {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	var buffer bytes.Buffer

	for i := 0; i < 128; i++ {
		buffer.WriteByte(alphabet[rand.Intn(len(alphabet))])
	}

	return buffer.String()
}
