package tasks // import github.com/mlimaloureiro/golog/models/tasks

import "time"

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
