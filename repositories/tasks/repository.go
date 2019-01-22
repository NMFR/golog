package tasks // import github.com/mlimaloureiro/golog/repositories/tasks

import taskModel "github.com/mlimaloureiro/golog/models/tasks"

// TaskRepositoryInterface represents a repository of Tasks
type TaskRepositoryInterface interface {
	SetTask(task taskModel.Task) error
	SetTasks(tasks taskModel.Collection) error
	GetTask(identifier string) (*taskModel.Task, error)
	GetTasks( /*from *time.Time, to *time.Time*/ ) (taskModel.Collection, error)
}
