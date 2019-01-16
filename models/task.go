package models // import github.com/mlimaloureiro/golog/models

import "time"

type TaskActivity struct {
	StartDate time.Time
	EndDate   time.Time // EndDate.IsZero() indicates a running task, meaning the 0001-01-01 EndDate cannot be represented
}

func (taskActivity TaskActivity) IsRunning() bool {
	return taskActivity.EndDate.IsZero()
}

func (taskActivity TaskActivity) Duration() time.Duration {
	endDate := taskActivity.EndDate
	if endDate.IsZero() {
		endDate = time.Now()
	}
	return endDate.Sub(taskActivity.StartDate)
}

type Task struct {
	Identifier string
	// Description string
	Activity []TaskActivity
}

func (task Task) GetRunningTaskActivity() *TaskActivity {
	for i, _ := range task.Activity {
		if task.Activity[i].IsRunning() {
			return &task.Activity[i]
		}
	}
	return nil
}

func (task Task) IsRunning() bool {
	return task.GetRunningTaskActivity() != nil
}

func (task Task) Duration() time.Duration {
	var duration time.Duration
	for i, _ := range task.Activity {
		duration += task.Activity[i].Duration()
	}
	return duration
}

type Tasks []Task

func (tasks Tasks) GetByIdentifier(identifier string) *Task {
	for i := range tasks {
		if tasks[i].Identifier == identifier {
			return &tasks[i]
		}
	}
	return nil
}
