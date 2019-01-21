package tasks // import github.com/mlimaloureiro/golog/repositories/tasks

// DeleterTaskRepositoryInterface is a TaskRepositoryInterface that can also delete tasks directly
type DeleterTaskRepositoryInterface interface {
	TaskRepositoryInterface
	DeleteTask(identifier string) error
	DeleteTasks() error
}
