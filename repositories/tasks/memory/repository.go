package memory // import github.com/mlimaloureiro/golog/repositories/tasks/memory

import (
	tasksModel "github.com/mlimaloureiro/golog/models/tasks"
)

func cloneTask(task tasksModel.Task) tasksModel.Task {
	newTask := tasksModel.Task{
		Identifier: task.Identifier,
		Activity:   make([]tasksModel.TaskActivity, len(task.Activity)),
	}
	copy(newTask.Activity, task.Activity)
	return newTask
}

func cloneTaskCollection(tasks tasksModel.Collection) tasksModel.Collection {
	newTasks := tasksModel.Collection{}
	for _, task := range tasks {
		newTasks = append(newTasks, cloneTask(task))
	}
	return newTasks
}

// TaskRepository is a Task repository that stores its data in memory
type TaskRepository struct {
	tasks tasksModel.Collection
}

// New creates a new csv TaskRepository
func New() TaskRepository {
	return TaskRepository{}
}

// SetTask will create or update the Task in the rerpository, if the task already exists in the repository its data will be overriden by the new Task
func (repository *TaskRepository) SetTask(task tasksModel.Task) error {
	taskClone := cloneTask(task)

	if taskPtr := repository.tasks.GetByIdentifier(task.Identifier); taskPtr == nil {
		repository.tasks = append(repository.tasks, taskClone)
	} else {
		*taskPtr = taskClone
	}

	return nil
}

// SetTasks will delete all tasks of the repository and insert the tasks passed by parameter
func (repository *TaskRepository) SetTasks(tasks tasksModel.Collection) (err error) {
	repository.tasks = cloneTaskCollection(tasks)
	return nil
}

// GetTask returns a Task from the repository by identifier
func (repository *TaskRepository) GetTask(identifier string) (*tasksModel.Task, error) {
	taskPtr := repository.tasks.GetByIdentifier(identifier)
	if taskPtr == nil {
		return nil, nil
	}

	task := cloneTask(*taskPtr)
	return &task, nil
}

// GetTasks returns all Tasks of the repository
func (repository *TaskRepository) GetTasks() (tasks tasksModel.Collection, err error) {
	return cloneTaskCollection(repository.tasks), nil
}
