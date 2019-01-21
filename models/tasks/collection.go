package tasks // import github.com/mlimaloureiro/golog/models/tasks

// Collection represents a collection of Task structs
type Collection []Task

// GetByIdentifier returns a Task from the collection by its identifier
func (tasks Collection) GetByIdentifier(identifier string) *Task {
	for i := range tasks {
		if tasks[i].Identifier == identifier {
			return &tasks[i]
		}
	}
	return nil
}
