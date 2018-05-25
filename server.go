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

// Server is the Diceroller API server. Use NewServer to construct and Run to start.
type Server struct {
	db *sql.DB
	r  *gin.Engine
}

// TODO this type likely needs to be shared outside of server context
type User struct {
	id        string
	createdAt time.Time
	name      string
}

// NewServer constructs a new instance of the DiceRoller server. Use Run() to start.
func NewServer(db *sql.DB) *Server {
	s := Server{db: db}
	r := gin.Default()

	// High level routes go here
	r.GET("/ping", s.pingHandler)
	r.POST("/dice/roll", s.diceRollHandler)
	r.POST("/room", s.addRoomHandler)
	r.POST("/user", s.addUserHandler)

	// Authenticated routes go here
	authenticated := r.Group("/")
	authenticated.Use(s.authenticated())
	{
		// TODO actually use a real join handler; right now this is just a dummy implementation for testing
		authenticated.POST("/room/:name/join", s.pingHandler)
	}

	s.r = r

	return &s
}

// Run starts the server; give it a port like ':8080' or include hostname as well.
// See gin's documentation on run for more info.
func (s *Server) Run(host string) error {
	return s.r.Run(host)
}

func (s *Server) authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")

		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		token := auth[len("Bearer "):]

		var user User
		err := s.db.QueryRow("SELECT id, created_at, name FROM users WHERE token = $1", token).Scan(&user.id, &user.createdAt, &user.name)
		switch {
		case err == sql.ErrNoRows:
			c.AbortWithStatus(http.StatusForbidden)
		case err != nil:
			log.Printf("Error while looking up bearer token: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
		default:
			c.Set("user", user)
			c.Next()
		}
	}
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

	var id string
	err := s.db.QueryRow("INSERT INTO users (token, name) VALUES ($1, $2) RETURNING id", token, name).Scan(&id)
	if err != nil {
		log.Printf("Error while creating user: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    id,
		"name":  name,
		"token": token,
	})
}

// TODO rename to more generic name
func generateToken() string {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	var buffer bytes.Buffer

	for i := 0; i < 128; i++ {
		buffer.WriteByte(alphabet[rand.Intn(len(alphabet))])
	}

	return buffer.String()
}
