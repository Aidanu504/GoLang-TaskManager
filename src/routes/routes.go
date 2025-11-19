package routes

import (
	"github.com/gin-gonic/gin"
	"goLang-taskmanager/src/handlers"
)

// Routes func
// All CRUD endpoint routes added
func Routes(h *handlers.TaskHandler) *gin.Engine {
    router := gin.Default()
    
    router.GET("/", func(c *gin.Context) {
		c.String(200, "Home Route Successful")
	})

    router.GET("/tasks", h.GetTasks)
    router.POST("/tasks", h.CreateTask)
    router.DELETE("/tasks/:id", h.DeleteTask)
    router.PUT("/tasks/:id", h.UpdateTask)

    return router
}