package tasks // import github.com/mlimaloureiro/golog/repositories/tasks

import tasksModel "github.com/mlimaloureiro/golog/models/tasks"

// TaskRepositoryInterface represents a repository of Tasks
type TaskRepositoryInterface interface {
	SetTask(task tasksModel.Task) error
	SetTasks(tasks tasksModel.Collection) error
	GetTask(identifier string) (*tasksModel.Task, error)
	GetTasks( /*from *time.Time, to *time.Time*/ ) (tasksModel.Collection, error)
}
