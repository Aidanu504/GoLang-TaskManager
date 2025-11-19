package handlers

import (
    "database/sql"
    "net/http"
    "github.com/gin-gonic/gin"
    "goLang-taskmanager/src/models"
	"goLang-taskmanager/src/utils"
    "strconv"
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
// No input from client just returns JSON of all tasks in database
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

// Create new tasks function
// Takes in JSON body with new task information
// Insters this into database and returns created task with ID
func (h *TaskHandler) CreateTask(c *gin.Context) {
    var task models.Task

    // Bind JSON body into task
    if err := c.ShouldBindJSON(&task); err != nil {
        utils.BadRequest(c, "invalid request body")
        return
    }

    // SQL query for inserting task
    // CreatedAt will be set by datbase and TaskID autoincremented
    res, err := h.DB.Exec(`
        INSERT INTO Tasks (TaskName, TaskDescription, IsCompleted)
        VALUES (?, ?, ?)
    `, task.TaskName, task.TaskDescription, task.IsCompleted)

    // if this fails return 500 response 
    if err != nil {
        utils.ServerError(c, "failed to create task")
        return
    }

    // Get the autoincremented taskID and put it back on the struct
    id, err := res.LastInsertId()
    if err != nil {
        utils.ServerError(c, "failed to get last insert id")
        return
    }
    task.TaskID = int(id)

    c.JSON(http.StatusCreated, task)
}

// Delete a task from the database
// It deletes the task with the given ID and returns 204 if successful
func (h *TaskHandler) DeleteTask(c *gin.Context) {

    // Grab task ID from url and convert to int
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        utils.BadRequest(c, "taskID not valid")
        return
    }

    // Delete task from the database
    // Return 500 response if this fails
    res, err := h.DB.Exec(`DELETE FROM Tasks WHERE TaskID = ?`, id)
    if err != nil {
        utils.ServerError(c, "failed to delete task")
        return
    }

    // Check if a row was actually deleted
    affected, err := res.RowsAffected()
    if err != nil || affected == 0 {
        utils.NotFound(c, "task not found in database")
        return
    }

    // Successful delete
    c.Status(http.StatusNoContent)
}