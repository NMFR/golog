package csv // import github.com/mlimaloureiro/golog/repositories/tasks/file/csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"

	taskModel "github.com/mlimaloureiro/golog/models/tasks"
)

const (
	taskStart taskAction = "start"
	taskStop  taskAction = "stop"
)

type taskAction string

func formatTime(at time.Time) string {
	return at.Format(time.RFC3339)
}

func parseTime(at string) (time.Time, error) {
	then, err := time.Parse(time.RFC3339, at)
	return then, err
}

// TaskRepository is a Task repository that stores its data in the CSV format
type TaskRepository struct {
	filePath string
}

// New creates a new csv TaskRepository
func New(filePath string) TaskRepository {
	return TaskRepository{filePath}
}

func (repository TaskRepository) writeTaskAction(
	writer io.Writer,
	identifier string,
	action taskAction,
	at time.Time,
) error {
	csvWriter := csv.NewWriter(writer)
	if err := csvWriter.Write([]string{identifier, string(action), formatTime(at)}); err != nil {
		return err
	}

	csvWriter.Flush()
	err := csvWriter.Error()

	return err
}

// StartTask starts the Task with the identifier if the Task is not already running
//  If the Task does not exist it will be created
func (repository TaskRepository) StartTask(identifier string) (err error) {
	file, err := os.OpenFile(repository.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); err == nil {
			err = closeErr
		}
	}()

	return repository.writeTaskAction(file, identifier, taskStart, time.Now())
}

// PauseTask pauses the Task with the identifier if the Task is already running
func (repository TaskRepository) PauseTask(identifier string) (err error) {
	file, err := os.OpenFile(repository.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); err == nil {
			err = closeErr
		}
	}()

	return repository.writeTaskAction(file, identifier, taskStop, time.Now())
}

// SetTask will create or update the Task in the rerpository,
//  if the task already exists in the repository its data will be overriden by the new Task
func (repository TaskRepository) SetTask(task taskModel.Task) error {
	tasks, err := repository.GetTasks()
	if err != nil {
		// Ignore "no such file" errors here:
		if pathErr, isPathErr := err.(*os.PathError); !isPathErr || pathErr.Err != syscall.ENOENT {
			return err
		}

		tasks = taskModel.Collection{}
	}

	if taskPtr := tasks.GetByIdentifier(task.Identifier); taskPtr == nil {
		tasks = append(tasks, task)
	} else {
		*taskPtr = task
	}

	err = repository.SetTasks(tasks)
	return err
}

// SetTasks will delete all tasks of the repository and insert the tasks passed by parameter
func (repository TaskRepository) SetTasks(tasks taskModel.Collection) (err error) {
	file, err := os.OpenFile(repository.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); err == nil {
			err = closeErr
		}
	}()

	for _, task := range tasks {
		for _, taskActivity := range task.Activity {
			if err := repository.writeTaskAction(file, task.Identifier, taskStart, taskActivity.StartDate); err != nil {
				return err
			}
			if !taskActivity.EndDate.IsZero() {
				if err := repository.writeTaskAction(file, task.Identifier, taskStop, taskActivity.EndDate); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// GetTasks returns all Tasks of the repository
func (repository TaskRepository) GetTasks() (tasks taskModel.Collection, err error) {
	file, err := os.OpenFile(repository.filePath, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := file.Close(); err == nil {
			err = closeErr
		}
	}()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	rawCsvData, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	tasks = taskModel.Collection{}
	for i, line := range rawCsvData {
		if len(line) != 3 {
			return nil, fmt.Errorf("csvfile: malformed line %d: %q", i, line)
		}

		identifier, action, timeString := line[0], taskAction(line[1]), line[2]
		actionTime, err := parseTime(timeString)
		if err != nil {
			return nil, err
		}

		task := tasks.GetByIdentifier(identifier)
		if task == nil {
			tasks = append(tasks, taskModel.Task{Identifier: identifier})
			task = &tasks[len(tasks)-1]
		}

		taskActivity := task.GetRunningTaskActivity()

		switch action {
		case taskStart:
			if taskActivity == nil {
				taskActivity := taskModel.TaskActivity{StartDate: actionTime}
				task.Activity = append(task.Activity, taskActivity)
			}
		case taskStop:
			if taskActivity != nil {
				taskActivity.EndDate = actionTime
			}
		}
	}

	return tasks, nil
}

// GetTask returns a Task from the repository by identifier
func (repository TaskRepository) GetTask(identifier string) (*taskModel.Task, error) {
	tasks, err := repository.GetTasks()
	if err != nil {
		return nil, err
	}

	return tasks.GetByIdentifier(identifier), nil
}
