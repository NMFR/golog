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

func formatTime(at time.Time) string {
	return at.Format(time.RFC3339)
}

func parseTime(at string) (time.Time, error) {
	then, err := time.Parse(time.RFC3339, at)
	if err != nil {
		return then, err
	}
	return then, nil
}

type CsvTaskRepository struct {
	readWriter io.ReadWriter
}

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

func (repository CsvTaskRepository) StartTask(identifier string) error {
	return repository.writeTaskAction(identifier, taskStart, time.Now())
}

func (repository CsvTaskRepository) PauseTask(identifier string) error {
	return repository.writeTaskAction(identifier, taskStop, time.Now())
}

func (repository CsvTaskRepository) GetTasks() (models.Tasks, error) {
	// TODO: check if repository.readWriter should be a Seeker or a file path to allow multiple reads from the same repository.readWriter?
	reader := csv.NewReader(repository.readWriter)
	reader.FieldsPerRecord = -1
	rawCsvData, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	taskMap := make(map[string]models.Task)
	for i, line := range rawCsvData {
		if len(line) != 3 {
			return nil, fmt.Errorf("csvfile: malformed line %d: %q", i, line)
		}

		identifier, action, timeString := line[0], taskAction(line[1]), line[2]
		actionTime, err := parseTime(timeString)
		if err != nil {
			return nil, err
		}

		task, inMap := taskMap[identifier]
		if inMap == false {
			task = models.Task{Identifier: identifier}
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

		taskMap[identifier] = task
	}

	tasks := models.Tasks{}
	for _, task := range taskMap {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (repository CsvTaskRepository) GetTask(identifier string) (*models.Task, error) {
	tasks, err := repository.GetTasks()
	if err != nil {
		return nil, err
	}

	return tasks.GetByIdentifier(identifier), nil
}
