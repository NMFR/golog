package models // import github.com/mlimaloureiro/golog/models

// TaskRepositoryInterface represents a repository of Tasks
// TODO: should we split this into a TaskRepositoryInterface and a TaskServiceInterface?
type TaskRepositoryInterface interface {
	StartTask(identifier string) error
	PauseTask(identifier string) error
	SetTask(task Task) error
	SetTasks(tasks Tasks) error // also serves as Clear() => SetTasks(nil)
	// DeleteTask(identifier string) error
	GetTasks( /*from *time.Time, to *time.Time*/ ) (Tasks, error)
	GetTask(identifier string) (*Task, error)
}
