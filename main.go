// code below from gin-gonic documentation
package main

import (
	"log"
	"github.com/aidanu504/golang-taskmanager/database"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// Untested since need to push for "database" import above 
func main() {
	db := database.DatabaseConnect()
	defer db.Close()

	router := gin.Default()
	
	router.GET("/", func(c *gin.Context) {
        c.String(200, "Home Route Successful")
    })

    log.Println("Server running")
    router.Run(":8080")
}