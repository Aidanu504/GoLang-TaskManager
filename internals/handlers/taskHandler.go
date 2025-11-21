package handlers

import (
	"database/sql"
	"goLang-taskmanager/internals/models"
	"goLang-taskmanager/internals/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	// Get the completed query parameter
    completedParam := c.Query("completed")

	// Fetch all tasks sorted by most recent
	// Now they are sorted by if completed or not
	var query string
	var section string
    if completedParam == "true" {
        query = `SELECT TaskID, TaskName, TaskDescription, IsCompleted, CreatedAt
                 FROM Tasks WHERE IsCompleted = 1 ORDER BY CreatedAt DESC`
		section = "completed"
    } else if completedParam == "false" {
        query = `SELECT TaskID, TaskName, TaskDescription, IsCompleted, CreatedAt
                 FROM Tasks WHERE IsCompleted = 0 ORDER BY CreatedAt DESC`
		section = "active"
    } else {
        query = `SELECT TaskID, TaskName, TaskDescription, IsCompleted, CreatedAt
                 FROM Tasks ORDER BY CreatedAt DESC`
		section = "all"
    }

	rows, err := h.DB.Query(query)

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

	//c.JSON(http.StatusOK, tasks)
	//Render using template
	c.HTML(http.StatusOK, "tasks.html", gin.H{
		"tasks": tasks,
		"section": section,
	})
}

// Create new tasks function
// Takes in JSON body with new task information
// Insters this into database and returns created task with ID
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task

	// Bind JSON body into task
	if err := c.ShouldBind(&task); err != nil {
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

	//c.JSON(http.StatusCreated, task)
	// Render single task using template
	c.HTML(http.StatusCreated, "task-item.html", gin.H{
		"task": task,
	})
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
	//c.Status(http.StatusNoContent)
	c.Header("HX-Trigger", "refresh-tasks")
	c.String(http.StatusOK, "")
}

// Put to update a task within the database
// It updates an existing task and returns the updated task JSON
func (h *TaskHandler) UpdateTask(c *gin.Context) {

	// Grab task ID from url and convert to int
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.BadRequest(c, "TaskID not valid")
		return
	}

	// Bind request body into the Task struct
	var task models.Task
	if err := c.ShouldBind(&task); err != nil {
		utils.BadRequest(c, "invalid request body")
		return
	}

	// Run update query
	// CreatedAt and TaskID aren't updated
	res, err := h.DB.Exec(`
        UPDATE Tasks
        SET TaskName = ?, TaskDescription = ?, IsCompleted = ?
        WHERE TaskID = ?
    `, task.TaskName, task.TaskDescription, task.IsCompleted, id)
	if err != nil {
		utils.ServerError(c, "failed to update task")
		return
	}

	// Ensure a row was updated
	affected, err := res.RowsAffected()
	if err != nil || affected == 0 {
		utils.NotFound(c, "task not found")
		return
	}

	// Load the updated task
	var updated models.Task
	err = h.DB.QueryRow(`
        SELECT TaskID, TaskName, TaskDescription, IsCompleted, CreatedAt
        FROM Tasks
        WHERE TaskID = ?
    `, id).Scan(&updated.TaskID, &updated.TaskName, &updated.TaskDescription, &updated.IsCompleted, &updated.CreatedAt)
	if err != nil {
		utils.ServerError(c, "failed to load updated task")
		return
	}

	// Return updated task
	//c.JSON(http.StatusOK, updated)
	// Render updated task using template
	c.HTML(http.StatusOK, "task-item.html", gin.H{
		"task": updated,
	})
}

// ToggleTask function
// Flips the IsCompleted status of a task and returns the updated task HTML
func (h *TaskHandler) ToggleTask(c *gin.Context) {
	// Grab task ID from url and convert to int
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.BadRequest(c, "taskID not valid")
		return
	}

	// Toggle the IsCompleted status in the database
	_, err = h.DB.Exec(`
        UPDATE Tasks
        SET IsCompleted = NOT IsCompleted
        WHERE TaskID = ?
    `, id)

	// Return 500 if the update fails
	if err != nil {
		utils.ServerError(c, "failed to toggle task")
		return
	}

	// Load the updated task back from the database
	var updated models.Task
	err = h.DB.QueryRow(`
        SELECT TaskID, TaskName, TaskDescription, IsCompleted, CreatedAt
        FROM Tasks
        WHERE TaskID = ?
    `, id).Scan(&updated.TaskID, &updated.TaskName, &updated.TaskDescription, &updated.IsCompleted, &updated.CreatedAt)
	if err != nil {
		utils.ServerError(c, "failed to load task")
		return
	}

	// Determine target section based on completion status
    targetSection := "active-tasks"
    if updated.IsCompleted {
        targetSection = "completed-tasks"
    }

	// Render the updated task HTML fragment
	c.HTML(http.StatusOK, "task-toggle.html", gin.H{
		"task": updated,
		"targetSection": targetSection,
	})
}

// EditTask function
// Loads a single task from the database and renders the edit form
func (h *TaskHandler) EditTask(c *gin.Context) {
	// Grab task ID from url and convert to int
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.BadRequest(c, "taskID not valid")
		return
	}

	// Load the task from the database
	var task models.Task
	err = h.DB.QueryRow(`
        SELECT TaskID, TaskName, TaskDescription, IsCompleted, CreatedAt
        FROM Tasks
        WHERE TaskID = ?
    `, id).Scan(&task.TaskID, &task.TaskName, &task.TaskDescription, &task.IsCompleted, &task.CreatedAt)

	// If this fails return 500 response
	if err != nil {
		utils.ServerError(c, "failed to load task")
		return
	}

	// Render edit page and pass in task data
	c.HTML(http.StatusOK, "task-edit.html", gin.H{
		"task": task,
	})
}

// ViewTask function
// Loads a single task from the database and renders the view template
func (h *TaskHandler) ViewTask(c *gin.Context) {
	// Grab task ID from url and convert to int
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.BadRequest(c, "taskID not valid")
		return
	}

	// Load task data from the database
	var task models.Task
	err = h.DB.QueryRow(`
        SELECT TaskID, TaskName, TaskDescription, IsCompleted, CreatedAt
        FROM Tasks
        WHERE TaskID = ?
    `, id).Scan(&task.TaskID, &task.TaskName, &task.TaskDescription, &task.IsCompleted, &task.CreatedAt)

	// If this fails return 500 response
	if err != nil {
		utils.ServerError(c, "failed to load task")
		return
	}

	// Render task item page
	c.HTML(http.StatusOK, "task-item.html", gin.H{
		"task": task,
	})
}

// ShowTaskDetails function
// Loads a single task from the database and renders the detailed view page
func (h *TaskHandler) ShowTaskDetails(c *gin.Context) {
	// Grab task ID from url and convert to int
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.BadRequest(c, "taskID not valid")
		return
	}

    // Query database for this task
	var task models.Task
	err = h.DB.QueryRow(`
        SELECT TaskID, TaskName, TaskDescription, IsCompleted, CreatedAt
        FROM Tasks
        WHERE TaskID = ?
    `, id).Scan(&task.TaskID, &task.TaskName, &task.TaskDescription, &task.IsCompleted, &task.CreatedAt)

    // If this fails return 500 response
	if err != nil {
		utils.ServerError(c, "failed to load task")
		return
	}

    // Render the task details page with the loaded task data
	c.HTML(http.StatusOK, "task-details.html", gin.H{
		"task": task,
	})
}

// ConfirmDeleteTask function
// Renders a confirmation modal for deleting a task
func (h *TaskHandler) ConfirmDeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.BadRequest(c, "taskID not valid")
		return
	}

	c.HTML(http.StatusOK, "task-delete.html", gin.H{
		"taskID": id,
	})
}