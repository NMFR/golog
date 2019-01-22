package ical // import github.com/mlimaloureiro/golog/repositories/tasks/file/ical

import (
	"io"
	"os"
	"syscall"
	"time"

	taskModel "github.com/mlimaloureiro/golog/models/tasks"

	"github.com/jordic/goics"
)

const prodID = "-//mlimaloureiro/golog"

func writeStrings(writer io.Writer, strings ...string) error {
	for _, str := range strings {
		_, err := io.WriteString(writer, str)
		if err != nil {
			return err
		}
	}
	return nil
}

type calendarConsumer struct {
	Calendar *goics.Calendar
}

func (consumer *calendarConsumer) ConsumeICal(calendar *goics.Calendar, err error) error {
	consumer.Calendar = calendar
	return err
}

// TaskRepository is a Task repository that stores its data in the ical (ics) format
type TaskRepository struct {
	filePath   string
	Version    string
	CALSCALE   string
	TimeFormat string // UTC time format string (layout string passed to time.Format)
}

// New creates a new ical TaskRepository
func New(filePath string) TaskRepository {
	return TaskRepository{
		filePath:   filePath,
		Version:    "2.0",
		CALSCALE:   "GREGORIAN",
		TimeFormat: "20060102T150405Z",
	}
}

func (repository TaskRepository) writeICalHeader(writer io.Writer) error {
	err := writeStrings(
		writer,
		"BEGIN:VCALENDAR\n",
		"SUMMARY:", repository.Version, "\n",
		"PRODID:", prodID, "\n",
		"CALSCALE:", repository.CALSCALE, "\n",
	)
	return err
}

func (repository TaskRepository) writeICalFooter(writer io.Writer) error {
	_, err := io.WriteString(writer, "END:VCALENDAR")
	return err
}

func (repository TaskRepository) writeICalEvent(
	writer io.Writer,
	task taskModel.Task,
	taskActivity taskModel.TaskActivity,
) error {
	if err := writeStrings(
		writer,
		"BEGIN:VEVENT\n",
		"SUMMARY:", task.Identifier, "\n",
		"DTSTART:", taskActivity.StartDate.UTC().Format(repository.TimeFormat), "\n",
	); err != nil {
		return err
	}

	if !taskActivity.IsRunning() {
		if err := writeStrings(
			writer,
			"DTEND:", taskActivity.EndDate.UTC().Format(repository.TimeFormat), "\n",
		); err != nil {
			return err
		}
	}

	if err := writeStrings(writer, "END:VEVENT\n"); err != nil {
		return err
	}

	return nil
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

	if err = repository.writeICalHeader(file); err != nil {
		return err
	}

	for _, task := range tasks {
		for _, taskActivity := range task.Activity {
			if err = repository.writeICalEvent(file, task, taskActivity); err != nil {
				return err
			}
		}
	}

	err = repository.writeICalFooter(file)
	return err
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

	icalDecoder := goics.NewDecoder(file)
	consumer := calendarConsumer{}
	err = icalDecoder.Decode(&consumer)
	if err != nil {
		if err == goics.ErrCalendarNotFound {
			return taskModel.Collection{}, nil
		}
		return nil, err
	}

	tasks = taskModel.Collection{}
	for _, event := range icalDecoder.Calendar.Events {
		summaryNode, startDateNode, endDateNode := event.Data["SUMMARY"], event.Data["DTSTART"], event.Data["DTEND"]
		if summaryNode == nil || startDateNode == nil {
			continue
		}

		startDate, err := startDateNode.DateDecode()
		if err != nil {
			return tasks, err
		}

		var endDate time.Time
		if endDateNode != nil {
			endDate, err = endDateNode.DateDecode()
			if err != nil {
				return tasks, err
			}
		}

		task := tasks.GetByIdentifier(summaryNode.Val)
		if task == nil {
			tasks = append(tasks, taskModel.Task{Identifier: summaryNode.Val})
			task = &tasks[len(tasks)-1]
		}

		(*task).Activity = append((*task).Activity, taskModel.TaskActivity{StartDate: startDate, EndDate: endDate})
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
