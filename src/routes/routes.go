package routes

import (
	"github.com/gin-gonic/gin"
	"goLang-taskmanager/src/handlers"
)

// Routes func
// Will add others need to commit to test get all tasks route
func Routes(h *handlers.TaskHandler) *gin.Engine {
    router := gin.Default()

    router.GET("/tasks", h.GetTasks)

    return router
}