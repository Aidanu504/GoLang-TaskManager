package handlers

import (
    "database/sql"
    "net/http"
    "github.com/gin-gonic/gin"
    "golang-taskmanager/src/models"
	"golang-taskmanager/src/utils"
)

// struct to hold database connection 
// funcions access through *sql.DB
type TaskHandler struct {
    DB *sql.DB
}

// function returns new task handler inststance from database connection 
func NewTaskHandler(db *sql.DB) *TaskHandler {
    return &TaskHandler{DB: db}
}

// Get all tasks function
func (h *TaskHandler) GetTasks(c *gin.Context) {

    // Fetch all tasks sorted by most recent
    rows, err := h.DB.Query(`
        SELECT TaskID, TaskName, TaskDescription, IsCompleted, CreatedAt
        FROM Tasks
        ORDER BY CreatedAt DESC
    `)

    // If the query fails return 500 response 
    if err != nil {
        utils.ServerError(c, "failed to query tasks")
        return
    }
    defer rows.Close()

    var tasks []models.Task

	// loop though results rows
    for rows.Next() {
        var t models.Task
        
        // scan row into task struct 
		// if this fails retrurn 500 response 
        if err := rows.Scan(&t.TaskID, &t.TaskName, &t.TaskDescription, &t.IsCompleted, &t.CreatedAt); err != nil {
            utils.ServerError(c, "failed to scan task")
            return
        }

		// add this to the slice of tasks
        tasks = append(tasks, t)
    }

    c.JSON(http.StatusOK, tasks)
}
