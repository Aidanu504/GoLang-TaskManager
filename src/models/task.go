package models

import "time"

type Task struct {
	TaskID int
	TaskName string
	TaskDescription string
	IsCompleted bool
	CreatedAt time.Time
}
