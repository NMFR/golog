package tasks // import github.com/mlimaloureiro/golog/repositories/tasks

// StarterPauserTaskRepositoryInterface is a TaskRepositoryInterface that can also start and stop tasks directly
type StarterPauserTaskRepositoryInterface interface {
	TaskRepositoryInterface
	StartTask(identifier string) error
	PauseTask(identifier string) error
}
