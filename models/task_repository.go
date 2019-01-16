package models // import github.com/mlimaloureiro/golog/models

type TaskRepositoryInterface interface {
	StartTask(identifier string) error
	PauseTask(identifier string) error
	// EditTask(task Task) error
	GetTasks( /*from *time.Time, to *time.Time*/ ) (Tasks, error)
	GetTask(identifier string) (*Task, error)
}
