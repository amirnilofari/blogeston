package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/amirnilofari/blogeston/models"
	"github.com/amirnilofari/blogeston/utils"
	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	var post models.Post

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := user_id.(int)
	fmt.Println(userID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Set the user ID of the post to the authenticated user and author
	post.UserID = userID
	post.Status = "draft"

	createPostQuery := `insert into posts (user_id, title, body,rating, created_at, updated_at, status) 
			  values ($1, $2, $3, $4, $5, $6, $7)
			  returning post_id`

	err := utils.DB.QueryRow(createPostQuery, post.UserID, post.Title, post.Body, post.Rating, time.Now(), time.Now(), post.Status).Scan(&post.ID)
	if post.Status != "draft" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Post status is invalid!"})
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post!"})
		fmt.Println("Database error: ", err)
		return
	} else {
		var user models.User
		findAuthorQuery := "SELECT first_name, last_name, FROM users WHERE user_id = $1"
		err := utils.DB.QueryRowContext(c, findAuthorQuery, post.UserID).Scan(&user.FirstName, &user.LastName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error in retrive author: " + err.Error()})
			return
		} else {
			post.Author = &user
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Post successfully created!",
			"data":    post,
		})
	}

}

func GetPosts(c *gin.Context) {

	var posts []models.Post
	var post models.Post

	query := "select * from posts"
	rows, err := utils.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"DB query error": err.Error()})
		return
	} else {

		for rows.Next() {
			err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Body, &post.Rating, &post.CreatedAt, &post.UpdatedAt, &post.Status)
			if err != nil {
				if err == sql.ErrNoRows {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts:" + err.Error()})
				return
			}
			post = models.Post{ID: post.ID, Title: post.Title, Body: post.Body, Author: post.Author, Rating: post.Rating, CreatedAt: post.CreatedAt, UpdatedAt: post.UpdatedAt}
			posts = append(posts, post)
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": posts})
}

func GetPost(c *gin.Context) {
	post_id := c.Param("id")
	var post models.Post
	findAuthorQuery := "SELECT * FROM posts WHERE post_id = $1"
	err := utils.DB.QueryRowContext(c, findAuthorQuery, post_id).Scan(&post.ID, &post.UserID, &post.Title, &post.Body, &post.Rating, &post.CreatedAt, &post.UpdatedAt, &post.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error in retrive a post: " + err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": post})
	}
}

func PublishPost(c *gin.Context) {
	post_id := c.Param("id")
	var post models.Post
	query := "SELECT post_id, status FROM posts WHERE post_id = $1"
	err := utils.DB.QueryRowContext(c, query, post_id).Scan(&post.ID, &post.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error in retrive a post: " + err.Error()})
		return
	} else {
		post.Status = "published"
		c.JSON(http.StatusOK, gin.H{"data": post})
	}
}
