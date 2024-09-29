package models

import "time"

type Comment struct {
	ID              int       `json:"comment_id"`
	PostID          int       `json:"post_id"`
	UserID          int       `json:"user_id"`
	Author          string    `json:"author"`
	Body            string    `json:"body" binding:"required"`
	ThumbsUpCount   int       `json:"thumbs_up_count"`
	ThumbsDownCount int       `json:"thumbs_down_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CommentReaction struct {
	ID           int       `json:"comment_reaction_id"`
	CommentID    int       `json:"comment_id"`
	UserID       int       `json:"user_id"`
	ReactionType int       `json:"reaction_type" binding:"required"` // 1 for thumbs-up, -1 for thumbs-down
	CreatedAt    time.Time `json:"created_at"`
}
