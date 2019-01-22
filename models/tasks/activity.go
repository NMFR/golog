package tasks // import github.com/mlimaloureiro/golog/models/tasks

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

// GetDuration returns the time.Duration of the TaskActivity date interval
func (taskActivity TaskActivity) GetDuration() time.Duration {
	endDate := taskActivity.EndDate
	if endDate.IsZero() {
		endDate = time.Now()
	}
	return endDate.Sub(taskActivity.StartDate)
}
