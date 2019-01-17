package models // import github.com/mlimaloureiro/golog/models

import "time"

// TaskActivity represents a date interval of a task execution
type TaskActivity struct {
	StartDate time.Time
	EndDate   time.Time // EndDate.IsZero() indicates a running task, meaning the 0001-01-01 EndDate cannot be represented
}

// IsRunning indicate if the TaskActivity date interval is in progress
func (taskActivity TaskActivity) IsRunning() bool {
	return taskActivity.EndDate.IsZero()
}

// Duration returns the time.Duration of the TaskActivity date interval
func (taskActivity TaskActivity) Duration() time.Duration {
	endDate := taskActivity.EndDate
	if endDate.IsZero() {
		endDate = time.Now()
	}
	return endDate.Sub(taskActivity.StartDate)
}

// Task represents a tasks execution time
type Task struct {
	Identifier string
	// Description string
	Activity []TaskActivity
}

// GetRunningTaskActivity returns the Task's current running TaskActivty or nil if the Task does not have a TaskActivity in progress
func (task Task) GetRunningTaskActivity() *TaskActivity {
	for i := range task.Activity {
		if task.Activity[i].IsRunning() {
			return &task.Activity[i]
		}
	}
	return nil
}

// IsRunning indicates if the Task is in progress or if it is paused
func (task Task) IsRunning() bool {
	return task.GetRunningTaskActivity() != nil
}

// Duration returns the time.Duration of the Task execution time
func (task Task) Duration() time.Duration {
	var duration time.Duration
	for i := range task.Activity {
		duration += task.Activity[i].Duration()
	}
	return duration
}

// Tasks represents a collection of Task structs
type Tasks []Task

// GetByIdentifier returns a Task from the collection by its identifier
func (tasks Tasks) GetByIdentifier(identifier string) *Task {
	for i := range tasks {
		if tasks[i].Identifier == identifier {
			return &tasks[i]
		}
	}
	return nil
}
