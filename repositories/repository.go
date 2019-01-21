package repositories // import github.com/mlimaloureiro/golog/repositories

import tasksModel "github.com/mlimaloureiro/golog/models/tasks"

// TaskRepositoryInterface represents a repository of Tasks
// TODO: should we split this into a TaskRepositoryInterface and a TaskServiceInterface?
type TaskRepositoryInterface interface {
	StartTask(identifier string) error
	PauseTask(identifier string) error
	SetTask(task tasksModel.Task) error
	SetTasks(tasks tasksModel.Collection) error // also serves as Clear() => SetTasks(nil)
	// DeleteTask(identifier string) error
	GetTasks( /*from *time.Time, to *time.Time*/ ) (tasksModel.Collection, error)
	GetTask(identifier string) (*tasksModel.Task, error)
}
