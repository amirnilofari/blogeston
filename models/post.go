package models

import "time"

type Post struct {
	ID        int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title" binding:"required"`
	Body      string    `json:"body" binding:"required"`
	Status    string    `json:"status"`
	Author    *User     `json:"author,omitempty"`
	Rating    float64   `json:"rating,omitempty"`
	Comments  []Comment `json:"comments,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Rating struct {
	ID        int       `json:"rating_id"`
	PostID    int       `json:"post_id" binding:"required"`
	UserID    int       `json:"user_id" binding:"required"`
	Rating    int       `json:"rating" binding:"required,min=1,max=5"`
	CreatedAt time.Time `json:"created_at"`
}
