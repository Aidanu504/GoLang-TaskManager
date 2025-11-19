// code below from gin-gonic documentation
package main

import (
	"log"
	"goLang-taskmanager/database"
	"goLang-taskmanager/src/handlers"
	"goLang-taskmanager/src/routes"
	_ "modernc.org/sqlite" 
)

// Server connection worked now need to test databse connection with GET and test newly implented router
func main() {
	db := database.DatabaseConnect()
	defer db.Close()

	taskHandler := handlers.NewTaskHandler(db)
	router := routes.Routes(taskHandler)

	log.Println("Server running")
	router.Run(":8080")
}
