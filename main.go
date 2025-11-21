// code below from gin-gonic documentation
package main

import (
	"log"
	"goLang-taskmanager/database"
	"goLang-taskmanager/internals/handlers"
	"goLang-taskmanager/internals/routes"
	_ "modernc.org/sqlite" 
)

// Server connection worked now need to test databse connection with GET and test newly implented router
func main() {
	db := database.DatabaseConnect()

	// Migrate DB to create tables if not exist
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Migration failed %v", err)
	}
	defer db.Close()

	taskHandler := handlers.NewTaskHandler(db)
	router := routes.Routes(taskHandler)

	log.Println("Server running")
	router.Run(":8080")
}
