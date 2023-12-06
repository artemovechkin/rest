package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"rest/db"
	"rest/delivery"
)

func main() {
	database := db.InitDatabase()
	service := delivery.NewService(database)

	router := gin.Default()
	delivery.InitEndPoints(router, service)

	router.Run(":8000")
}
