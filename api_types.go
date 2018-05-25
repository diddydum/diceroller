package main

import "time"

type User struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
}

type CreateUserRequest struct {
	Name string `form:"name" json:"name" binding:"required,max=32,min=2,alphanum"`
}

type CreateUserResponse struct {
	User User `json:"user"`
}
