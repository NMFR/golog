package repositories // import github.com/mlimaloureiro/golog/repositories

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/mlimaloureiro/golog/models"
)

const (
	taskStart taskAction = "start"
	taskStop  taskAction = "stop"
)

type taskAction string

// CsvTaskRepository is a Task repository that stores its data in the CSV format
type CsvTaskRepository struct {
	readWriter io.ReadWriter
}

// NewCsvTaskRepository creates a new CsvTaskRepository
func NewCsvTaskRepository(readWriter io.ReadWriter) CsvTaskRepository {
	return CsvTaskRepository{readWriter}
}

func (repository CsvTaskRepository) writeTaskAction(identifier string, action taskAction, at time.Time) error {
	writer := csv.NewWriter(repository.readWriter)
	if err := writer.Write([]string{identifier, string(action), formatTime(at)}); err != nil {
		return err
	}

	writer.Flush()
	err := writer.Error()

	return err
}

// StartTask starts the Task with the identifier if the Task is not already running
//  If the Task does not exist it will be created
func (repository CsvTaskRepository) StartTask(identifier string) error {
	if _, err := trySeek(repository.readWriter, 0, io.SeekEnd); err != nil {
		return err
	}
	return repository.writeTaskAction(identifier, taskStart, time.Now())
}

// PauseTask pauses the Task with the identifier if the Task is already running
func (repository CsvTaskRepository) PauseTask(identifier string) error {
	if _, err := trySeek(repository.readWriter, 0, io.SeekEnd); err != nil {
		return err
	}
	return repository.writeTaskAction(identifier, taskStop, time.Now())
}

// SetTask will create or update the Task in the rerpository, if the task already exists in the repository its data will be overriden by the new Task
func (repository CsvTaskRepository) SetTask(task models.Task) error {
	tasks, err := repository.GetTasks()
	if err != nil {
		return err
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
func (repository CsvTaskRepository) SetTasks(tasks models.Tasks) error {
	if _, err := trySeek(repository.readWriter, 0, io.SeekStart); err != nil {
		return err
	}

	for _, task := range tasks {
		for _, taskActivity := range task.Activity {
			if err := repository.writeTaskAction(task.Identifier, taskStart, taskActivity.StartDate); err != nil {
				return err
			}
			if !taskActivity.EndDate.IsZero() {
				if err := repository.writeTaskAction(task.Identifier, taskStop, taskActivity.EndDate); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// GetTasks returns all Tasks of the repository
func (repository CsvTaskRepository) GetTasks() (models.Tasks, error) {
	if _, err := trySeek(repository.readWriter, 0, io.SeekStart); err != nil {
		return nil, err
	}

	reader := csv.NewReader(repository.readWriter)
	reader.FieldsPerRecord = -1
	rawCsvData, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	tasks := models.Tasks{}
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
			tasks = append(tasks, models.Task{Identifier: identifier})
			task = &tasks[len(tasks)-1]
		}

		taskActivity := task.GetRunningTaskActivity()

		switch action {
		case taskStart:
			if taskActivity == nil {
				taskActivity := models.TaskActivity{StartDate: actionTime}
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
func (repository CsvTaskRepository) GetTask(identifier string) (*models.Task, error) {
	tasks, err := repository.GetTasks()
	if err != nil {
		return nil, err
	}

	return tasks.GetByIdentifier(identifier), nil
}
