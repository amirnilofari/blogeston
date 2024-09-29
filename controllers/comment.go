package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/amirnilofari/hash-go-mysql/models"
	"github.com/amirnilofari/hash-go-mysql/utils"
	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	var comment models.Comment

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify that the post exists
	post_id := c.Param("id")

	var exists bool
	err := utils.DB.QueryRowContext(c, "SELECT EXISTS(SELECT 1 FROM posts WHERE post_id=$1)", post_id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post not found"})
		return
	} else {
		// string to int
		i, err := strconv.Atoi(post_id)
		if err != nil {
			// ... handle error
			panic(err)
		}

		comment.PostID = i
	}

	//var commentAuthor models.User
	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := user_id.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	//commentAuthor.ID = userID
	comment.UserID = userID

	// Insert the new comment
	query := `
        INSERT INTO comments (post_id, user_id, body, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING comment_id, created_at, updated_at;
    `
	err = utils.DB.QueryRowContext(c, query, comment.PostID, comment.UserID, comment.Body, time.Now(), time.Now()).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment" + err.Error()})
		return
	} else {
		var user models.User
		findAuthorQuery := "SELECT first_name, last_name FROM users WHERE user_id = $1"
		err := utils.DB.QueryRowContext(c, findAuthorQuery, comment.UserID).Scan(&user.FirstName, &user.LastName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error in retrive author: " + err.Error()})
			return
		} else {
			comment.Author = user.FirstName + " " + user.LastName
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Comment successfully created!",
			"data":    comment,
		})
	}

	//c.JSON(http.StatusOK, gin.H{"message": "Comment created successfully", "comment": comment})

}

func GetComments(c *gin.Context) {
	post_id := c.Param("id")

	// Query to retrieve comments along with user details
	query := `
        SELECT
            c.comment_id, c.post_id, c.user_id, c.body, c.created_at, c.updated_at, c.thumbs_up_count, c.thumbs_down_count,
            u.user_id, u.first_name, u.last_name, u.email, u.created_at, u.updated_at
        FROM comments c
        JOIN users u ON c.user_id = u.user_id
        WHERE c.post_id = $1
        ORDER BY c.created_at ASC;
    `
	rows, err := utils.DB.QueryContext(c, query, post_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments" + err.Error()})
		return
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		var user models.User
		if err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.UserID, &comment.Body, &comment.CreatedAt, &comment.UpdatedAt, &comment.ThumbsUpCount, &comment.ThumbsDownCount,
			&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning comment"})
			return
		}
		comment.Author = user.FirstName + " " + user.LastName
		comments = append(comments, comment)
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})

}

func ReactToComment(c *gin.Context) {
	var reaction models.CommentReaction
	if err := c.ShouldBindJSON(&reaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(reaction)

	// Get user ID from context (authentication middleware)
	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := user_id.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	reaction.UserID = userID
	fmt.Println(reaction.UserID)

	if reaction.ReactionType != 1 && reaction.ReactionType != -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reaction type!"})
	}

	comment_id := c.Param("comment_id")
	var commentExists bool
	err := utils.DB.QueryRowContext(c, "SELECT EXISTS(SELECT 1 FROM comments WHERE post_id=$1)", comment_id).Scan(&commentExists)
	if err != nil || !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post not found"})
		return
	} else {
		// string to int
		c, err := strconv.Atoi(comment_id)
		if err != nil {
			// ... handle error
			panic(err)
		}

		reaction.CommentID = c
	}

	fmt.Println("reaction comment id:", reaction.CommentID)

	query := `
		insert into comment_reactions (comment_id, user_id, reaction_type, created_at)
		values ($1, $2, $3, $4)
		on conflict (comment_id, user_id) do update
		set reaction_type = excluded.reaction_type,
			created_at = excluded.created_at
		returning comment_reaction_id, created_at;
	`

	err = utils.DB.QueryRowContext(c, query, reaction.CommentID, reaction.UserID, reaction.ReactionType, time.Now()).Scan(&reaction.ID, &reaction.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit reaction"})
		return
	}

	err = updateCommentReactionCounts(c, reaction.CommentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reaction counts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reaction submitted", "reaction": reaction})
}

func updateCommentReactionCounts(ctx context.Context, commentID int) error {
	var thumbsUpCount, thumbsDownCount int
	query := `
		select
			sum(case when reaction_type = 1 then 1 else 0 end) as thumbs_up_count,
			sum(case when reaction_type = -1 then 1 else 0 end) as thumbs_down_count
		from comment_reactions
		where comment_id = $1

	`

	err := utils.DB.QueryRowContext(ctx, query, commentID).Scan(&thumbsUpCount, &thumbsDownCount)
	if err != nil {
		fmt.Println(err)
		return err
	}

	updateQuery := `
		update comments
		set thumbs_up_count = $1,
			thumbs_down_count = $2
		where comment_id = $3;
	`

	_, err = utils.DB.ExecContext(ctx, updateQuery, thumbsUpCount, thumbsDownCount, commentID)
	fmt.Println(err)
	return err
}
