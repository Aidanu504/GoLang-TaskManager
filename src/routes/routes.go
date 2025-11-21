package routes

import (
	"goLang-taskmanager/src/handlers"

	"github.com/gin-gonic/gin"
)

// Routes func
// All CRUD endpoint routes added
func Routes(h *handlers.TaskHandler) *gin.Engine {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static/css")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.GET("/modal/close", func(c *gin.Context) {
		c.String(200, "")
	})

	router.GET("/tasks", h.GetTasks)
	router.POST("/tasks", h.CreateTask)
	router.DELETE("/tasks/:id", h.DeleteTask)
	router.PUT("/tasks/:id", h.UpdateTask)
	router.PUT("/tasks/:id/toggle", h.ToggleTask)
	router.GET("/tasks/:id/edit", h.EditTask)
	router.GET("/tasks/:id/view", h.ViewTask)
	router.GET("/tasks/:id/details", h.ShowTaskDetails)
	router.GET("/tasks/:id/confirm-delete", h.ConfirmDeleteTask)

	return router
}
