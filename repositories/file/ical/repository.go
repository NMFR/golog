package ical // import github.com/mlimaloureiro/golog/repositories/file/ical

import (
	"io"
	"time"

	tasksModel "github.com/mlimaloureiro/golog/models/tasks"

	"github.com/jordic/goics"
)

const prodID = "-//mlimaloureiro/golog"

func formatTime(at time.Time) string {
	return at.Format(time.RFC3339)
}

func parseTime(at string) (time.Time, error) {
	then, err := time.Parse(time.RFC3339, at)
	return then, err
}

func tryParseTime(str string) time.Time {
	date, _ := parseTime(str)
	return date
}

func trySeek(readWriter io.ReadWriter, offset int64, whence int) (int64, error) {
	seeker, isSeekable := readWriter.(io.Seeker)
	if isSeekable == false {
		return 0, nil
	}
	return seeker.Seek(offset, whence)
}

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
	readWriter io.ReadWriter
	Version    string
	CALSCALE   string
	TimeFormat string // UTC time format string (layout string passed to time.Format)
}

// NewTaskRepository creates a new TaskRepository
func NewTaskRepository(readWriter io.ReadWriter) TaskRepository {
	repository := TaskRepository{readWriter: readWriter}
	repository.Version = "2.0"
	repository.CALSCALE = "GREGORIAN"
	repository.TimeFormat = "20060102T150405Z"
	return repository
}

func (repository TaskRepository) writeICalHeader() error {
	err := writeStrings(
		repository.readWriter,
		"BEGIN:VCALENDAR\n",
		"SUMMARY:", repository.Version, "\n",
		"PRODID:", prodID, "\n",
		"CALSCALE:", repository.CALSCALE, "\n",
	)
	return err
}

func (repository TaskRepository) writeICalFooter() error {
	_, err := io.WriteString(repository.readWriter, "END:VCALENDAR")
	return err
}

func (repository TaskRepository) writeICalEvent(task tasksModel.Task, taskActivity tasksModel.TaskActivity) error {
	if err := writeStrings(
		repository.readWriter,
		"BEGIN:VEVENT\n",
		"SUMMARY:", task.Identifier, "\n",
		"DTSTART:", taskActivity.StartDate.UTC().Format(repository.TimeFormat), "\n",
	); err != nil {
		return err
	}

	if !taskActivity.IsRunning() {
		if err := writeStrings(
			repository.readWriter,
			"DTEND:", taskActivity.EndDate.UTC().Format(repository.TimeFormat), "\n",
		); err != nil {
			return err
		}
	}

	if err := writeStrings(repository.readWriter, "END:VEVENT\n"); err != nil {
		return err
	}

	return nil
}

// StartTask starts the Task with the identifier if the Task is not already running
//  If the Task does not exist it will be created
func (repository TaskRepository) StartTask(identifier string) error {
	tasks, err := repository.GetTasks()
	if err != nil {
		return err
	}

	task := tasks.GetByIdentifier(identifier)
	if task == nil {
		tasks = append(tasks, tasksModel.Task{Identifier: identifier, Activity: []tasksModel.TaskActivity{}})
		task = &tasks[len(tasks)-1]
	}

	if (*task).IsRunning() {
		return nil
	}
	(*task).Activity = append((*task).Activity, tasksModel.TaskActivity{StartDate: time.Now()})

	err = repository.SetTasks(tasks)
	return err
}

// PauseTask pauses the Task with the identifier if the Task is already running
func (repository TaskRepository) PauseTask(identifier string) error {
	tasks, err := repository.GetTasks()
	if err != nil {
		return err
	}

	task := tasks.GetByIdentifier(identifier)
	if task == nil {
		return nil
	}

	taskActivity := (*task).GetRunningTaskActivity()
	if taskActivity == nil {
		return nil
	}

	(*taskActivity).EndDate = time.Now()

	err = repository.SetTasks(tasks)
	return err
}

// SetTask will create or update the Task in the rerpository, if the task already exists in the repository its data will be overriden by the new Task
func (repository TaskRepository) SetTask(task tasksModel.Task) error {
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
func (repository TaskRepository) SetTasks(tasks tasksModel.Collection) error {
	if _, err := trySeek(repository.readWriter, 0, io.SeekStart); err != nil {
		return err
	}

	if err := repository.writeICalHeader(); err != nil {
		return err
	}

	for _, task := range tasks {
		for _, taskActivity := range task.Activity {
			if err := repository.writeICalEvent(task, taskActivity); err != nil {
				return err
			}
		}
	}

	err := repository.writeICalFooter()
	return err
}

// GetTasks returns all Tasks of the repository
func (repository TaskRepository) GetTasks() (tasksModel.Collection, error) {
	if _, err := trySeek(repository.readWriter, 0, io.SeekStart); err != nil {
		return nil, err
	}

	icalDecoder := goics.NewDecoder(repository.readWriter)
	consumer := calendarConsumer{}
	err := icalDecoder.Decode(&consumer)
	if err != nil {
		if err == goics.ErrCalendarNotFound {
			return tasksModel.Collection{}, nil
		}
		return nil, err
	}

	tasks := tasksModel.Collection{}
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
			tasks = append(tasks, tasksModel.Task{Identifier: summaryNode.Val})
			task = &tasks[len(tasks)-1]
		}

		(*task).Activity = append((*task).Activity, tasksModel.TaskActivity{StartDate: startDate, EndDate: endDate})
	}

	return tasks, nil
}

// GetTask returns a Task from the repository by identifier
func (repository TaskRepository) GetTask(identifier string) (*tasksModel.Task, error) {
	tasks, err := repository.GetTasks()
	if err != nil {
		return nil, err
	}

	return tasks.GetByIdentifier(identifier), nil
}
