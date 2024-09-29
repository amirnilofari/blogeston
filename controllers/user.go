package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/amirnilofari/blogeston/models"
	"github.com/amirnilofari/blogeston/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check user exists in DB
	var exists bool
	queryExistUser := "select exists(select 1 from users where email=$1);"
	err := utils.DB.QueryRowContext(c, queryExistUser, input.Email).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
	}

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered!"})
		return
	}

	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	input.Password = hashedPassword

	// Save user to the database
	query := `insert into users (first_name,last_name, email, password, created_at, updated_at, role) 
			  values ($1, $2, $3, $4, $5, $6, $7)
			  returning user_id`

	err = utils.DB.QueryRow(query, input.FirstName, input.LastName, input.Email, input.Password, time.Now(), time.Now(), input.Role).Scan(&input.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Registration successful!",
		})
	}

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("error to compare:", err)
		return false
	} else {
		return true
	}
}

func Login(c *gin.Context) {

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Fetch user from database
	var storedUser models.User
	err := utils.DB.QueryRow("SELECT user_id, email, password FROM users WHERE email=$1", input.Email).Scan(&storedUser.ID, &storedUser.Email, &storedUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query user"})
		return
	}
	fmt.Println(storedUser)

	// Compare password
	isValid := CheckPasswordHash(input.Password, storedUser.Password)
	fmt.Println(isValid)
	if isValid {
		token, err := utils.GenerateJWT(storedUser.ID, storedUser.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
		//c.JSON(http.StatusOK, gin.H{"token": "hellooooooooo"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// // Generate JWT token
	// tokenString, err := GenerateJWT(storedUser)
	// if err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
	//     return
	// }

	// c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func GetAllUsers(c *gin.Context) {
	// var users []models.User
	// var user models.User

	query := `
		SELECT
            u.user_id, u.first_name, u.last_name, u.email, u.created_at, u.updated_at, u.role,
            p.post_id, p.title, p.body, p.created_at, p.updated_at
        FROM users u
        LEFT JOIN posts p ON u.user_id = p.user_id
		ORDER BY u.user_id
	`

	rows, err := utils.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"DB query error": err.Error()})
		return
	}
	defer rows.Close()

	usersMap := make(map[int]*models.User)
	for rows.Next() {
		var (
			userID        int
			userFirstname string
			userLastname  string
			email         string
			userCreatedAt time.Time
			userUpdatedAt time.Time
			userRole      string
			postID        sql.NullInt64
			title         sql.NullString
			body          sql.NullString
			postCreatedAt sql.NullTime
			postUpdatedAt sql.NullTime
		)

		err := rows.Scan(
			&userID,
			&userFirstname,
			&userLastname,
			&email,
			&userCreatedAt,
			&userUpdatedAt,
			&userRole,
			&postID,
			&title,
			&body,
			&postCreatedAt,
			&postUpdatedAt,
		)
		if err != nil {
			fmt.Println("Error in scan get all users" + err.Error())
			return
		}

		user, exists := usersMap[userID]
		if !exists {
			user = &models.User{
				ID:        userID,
				FirstName: userFirstname,
				LastName:  userLastname,
				Email:     email,
				Role:      userRole,
				CreatedAt: userCreatedAt,
				UpdatedAt: userUpdatedAt,
				Posts:     []models.Post{},
			}
			usersMap[userID] = user
		}

		if postID.Valid {
			post := models.Post{
				ID:        int(postID.Int64),
				Title:     title.String,
				Body:      body.String,
				CreatedAt: postCreatedAt.Time,
				UpdatedAt: postUpdatedAt.Time,
			}
			user.Posts = append(user.Posts, post)
		}
	}

	// Convert the map to a slice
	users := make([]models.User, 0, len(usersMap))
	for _, user := range usersMap {
		users = append(users, *user)
	}

	c.JSON(http.StatusOK, gin.H{"data": users})

	// rows, err := utils.DB.Query(query)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"DB query error": err.Error()})
	// 	return
	// } else {

	// 	for rows.Next() {
	// 		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.Role)
	// 		if err != nil {
	// 			if err == sql.ErrNoRows {
	// 				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 				return
	// 			}
	// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users:" + err.Error()})
	// 			return
	// 		}
	// 		user = models.User{ID: user.ID, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Posts: user.Posts, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Role: user.Role}
	// 		users = append(users, user)
	// 	}
	// }

	// c.JSON(http.StatusOK, gin.H{"data": user})

}
