package main

import (
	"log"
	"os"

	"github.com/amirnilofari/blogeston/routes"
	"github.com/amirnilofari/blogeston/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	connStr := os.Getenv("DB_CONNECTION")
	err = utils.ConnectDB(connStr)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	//utils.ConnectDB()

	r := gin.Default()

	routes.PublicRoutes(r)
	routes.PrivateRoutes(r)

	r.Run(":8081")
}
