// code below from gin-gonic documentation
package main

import (
	"log"

	"goLang-taskmanager/database"
	"goLang-taskmanager/src/handlers"
	"goLang-taskmanager/src/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// Server connection worked now need to test databse connection with GET and test newly implented router
func main() {
	db := database.DatabaseConnect()
	defer db.Close()

	taskHandler := handlers.NewTaskHandler(db)
	router := routes.Routes(taskHandler)

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Home Route Successful")
	})

	log.Println("Server running")
	router.Run(":8080")
}
