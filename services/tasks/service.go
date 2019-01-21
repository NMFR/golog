package tasks // import github.com/mlimaloureiro/golog/services/tasks

import (
	"time"

	tasksModel "github.com/mlimaloureiro/golog/models/tasks"
	tasksRepositories "github.com/mlimaloureiro/golog/repositories/tasks"
	"github.com/mlimaloureiro/golog/repositories/tasks/file"
)

// TaskService represents a service that can perform several operations on a repositories.TaskRepositoryInterface
type TaskService struct {
	repository tasksRepositories.TaskRepositoryInterface
}

// New creates a new TaskService
func New(repository tasksRepositories.TaskRepositoryInterface) TaskService {
	return TaskService{repository}
}

// SetTask will create or update the Task in the rerpository, if the task already exists in the repository its data will be overriden by the new Task
func (service TaskService) SetTask(task tasksModel.Task) error {
	return service.repository.SetTask(task)
}

// SetTasks will delete all tasks of the repository and insert the tasks passed by parameter
func (service TaskService) SetTasks(tasks tasksModel.Collection) error {
	return service.repository.SetTasks(tasks)
}

// GetTask returns a Task from the repository by identifier
func (service TaskService) GetTask(identifier string) (*tasksModel.Task, error) {
	return service.repository.GetTask(identifier)
}

// GetTasks returns all Tasks of the repository
func (service TaskService) GetTasks() (tasksModel.Collection, error) {
	return service.repository.GetTasks()
}

// StartTask starts the Task with the identifier if the Task is not already running
//  If the Task does not exist it will be created
func (service TaskService) StartTask(identifier string) error {
	if starterPauserRepository, ok := service.repository.(tasksRepositories.StarterPauserTaskRepositoryInterface); ok {
		return starterPauserRepository.StartTask(identifier)
	}

	task, err := service.GetTask(identifier)
	if err != nil {
		return err
	}

	if task == nil {
		task = &tasksModel.Task{Identifier: identifier, Activity: []tasksModel.TaskActivity{}}
	}

	if (*task).IsRunning() {
		return nil
	}
	(*task).Activity = append((*task).Activity, tasksModel.TaskActivity{StartDate: time.Now()})

	err = service.SetTask(*task)
	return err
}

// PauseTask pauses the Task with the identifier if the Task is already running
func (service TaskService) PauseTask(identifier string) error {
	if starterPauserRepository, ok := service.repository.(tasksRepositories.StarterPauserTaskRepositoryInterface); ok {
		return starterPauserRepository.PauseTask(identifier)
	}

	task, err := service.GetTask(identifier)
	if err != nil || task == nil {
		return err
	}

	taskActivity := (*task).GetRunningTaskActivity()
	if taskActivity == nil {
		return nil
	}

	(*taskActivity).EndDate = time.Now()

	err = service.SetTask(*task)
	return err
}

// DeleteTask removes the task with the identifier from the repository
func (service TaskService) DeleteTask(identifier string) error {
	if deleterRepository, ok := service.repository.(tasksRepositories.DeleterTaskRepositoryInterface); ok {
		return deleterRepository.DeleteTask(identifier)
	}

	tasks, err := service.GetTasks()
	if err != nil {
		return err
	}

	newTasks := tasksModel.Collection{}
	for _, task := range tasks {
		if task.Identifier != identifier {
			newTasks = append(newTasks, task)
		}
	}

	err = service.SetTasks(newTasks)

	return nil
}

// DeleteTasks removes all tasks from the repository
func (service TaskService) DeleteTasks() error {
	if deleterRepository, ok := service.repository.(tasksRepositories.DeleterTaskRepositoryInterface); ok {
		return deleterRepository.DeleteTasks()
	}

	return service.SetTasks(nil)
}

// Export the tasks to filePath in the specified format
func (service TaskService) Export(format file.Format, filePath string) error {
	repository, err := file.GetTaskFileRepository(format, filePath)
	if err != nil {
		return err
	}

	tasks, err := service.repository.GetTasks()
	if err != nil {
		return err
	}

	err = repository.SetTasks(tasks)
	if err != nil {
		return err
	}

	return nil
}
