package repositories // import github.com/mlimaloureiro/golog/repositories

import (
	"bytes"
	"fmt"
	"math"
	"regexp"
	"testing"
	"time"

	"github.com/mlimaloureiro/golog/models"

	"github.com/stretchr/testify/assert"
)

const (
	testICalEmptyContent = `BEGIN:VCALENDAR
SUMMARY:2.0
PRODID:-//mlimaloureiro/golog
CALSCALE:GREGORIAN
END:VCALENDAR`
	testICalContent = `BEGIN:VCALENDAR
SUMMARY:2.0
PRODID:-//mlimaloureiro/golog
CALSCALE:GREGORIAN
BEGIN:VEVENT
SUMMARY:first-task
DTSTART:20190101T100000Z
DTEND:20190101T100100Z
END:VEVENT
BEGIN:VEVENT
SUMMARY:first-task
DTSTART:20190101T100200Z
DTEND:20190101T100600Z
END:VEVENT
BEGIN:VEVENT
SUMMARY:second-task
DTSTART:20190101T100000Z
DTEND:20190101T100400Z
END:VEVENT
BEGIN:VEVENT
SUMMARY:last-task
DTSTART:20190101T100000Z
END:VEVENT
END:VCALENDAR`
	otherTestICalContent = `BEGIN:VCALENDAR
SUMMARY:2.0
PRODID:-//mlimaloureiro/golog
CALSCALE:GREGORIAN
BEGIN:VEVENT
SUMMARY:first-task
DTSTART:20100601T150000Z
END:VEVENT
BEGIN:VEVENT
SUMMARY:second-task
DTSTART:20190101T100000Z
DTEND:20190101T100400Z
END:VEVENT
BEGIN:VEVENT
SUMMARY:last-task
DTSTART:20190101T100000Z
END:VEVENT
END:VCALENDAR`
)

func TestICalTaskRepository(t *testing.T) {
	t.Run("implements TaskRepositoryInterface", func(t *testing.T) {
		buffer := bytes.Buffer{}
		repository := NewICalTaskRepository(&buffer)
		assert.Implements(t, (*models.TaskRepositoryInterface)(nil), repository)
	})

	t.Run("GetTasks()", func(t *testing.T) {
		buffer := bytes.Buffer{}
		buffer.WriteString(testICalContent)
		taskRepository := NewICalTaskRepository(&buffer)
		tasks, err := taskRepository.GetTasks()

		assert.NoError(t, err)
		assert.Equal(t, 3, len(tasks))

		taskMap := make(map[string]models.Task)
		for _, task := range tasks {
			taskMap[task.Identifier] = task
		}

		taskCases := []struct {
			taskIdentifier string
			isRunning      bool
			duration       time.Duration
		}{
			{"first-task", false, 5 * time.Minute},
			{"second-task", false, 4 * time.Minute},
			{"last-task", true, time.Now().Sub(tryParseTime("2019-01-01T10:00:00Z"))},
		}
		for _, taskCase := range taskCases {
			taskCase := taskCase
			t.Run(fmt.Sprintf("task %s", taskCase.taskIdentifier), func(t *testing.T) {
				t.Parallel()
				assert.Contains(t, taskMap, taskCase.taskIdentifier)
				task := taskMap[taskCase.taskIdentifier]
				assert.Equal(t, taskCase.isRunning, task.IsRunning(), "task.IsRunning()")
				caseDuration := taskCase.duration
				taskDuration := task.Duration()
				if task.IsRunning() {
					if diff := time.Duration(math.Abs(float64(taskCase.duration - taskDuration))); diff < (1 * time.Second) {
						caseDuration = taskDuration
					}
				}
				assert.Equal(t, caseDuration, taskDuration, "task.Duration()")
			})
		}
	})

	t.Run("GetTask()", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			buffer := bytes.Buffer{}
			buffer.WriteString(testICalContent)
			taskRepository := NewICalTaskRepository(&buffer)
			task, err := taskRepository.GetTask("first-task")

			assert.NoError(t, err)
			if assert.NotNil(t, task) == false {
				return
			}

			assert.Equal(t, "first-task", task.Identifier, "task.Identifier")
			assert.Equal(t, false, task.IsRunning(), "task.IsRunning()")
			assert.Equal(t, 5*time.Minute, task.Duration(), "task.Duration()")
		})

		t.Run("success (unknown task)", func(t *testing.T) {
			buffer := bytes.Buffer{}
			buffer.WriteString(testICalContent)
			taskRepository := NewICalTaskRepository(&buffer)
			task, err := taskRepository.GetTask("unknown-task")

			assert.NoError(t, err)
			assert.Nil(t, task)
		})
	})

	t.Run("StartTask() PauseTask()", func(t *testing.T) {
		buffer := bytes.Buffer{}
		taskRepository := NewICalTaskRepository(&buffer)

		err := taskRepository.StartTask("first-task")
		assert.NoError(t, err)
		err = taskRepository.PauseTask("first-task")
		assert.NoError(t, err)

		err = taskRepository.StartTask("second-task")
		assert.NoError(t, err)
		err = taskRepository.PauseTask("second-task")
		assert.NoError(t, err)

		err = taskRepository.StartTask("first-task")
		assert.NoError(t, err)
		err = taskRepository.PauseTask("first-task")
		assert.NoError(t, err)

		err = taskRepository.StartTask("last-task")
		assert.NoError(t, err)

		dateRegExpPattern := "\\d{4}\\d{2}\\d{2}T\\d{2}\\d{2}\\d{2}Z"
		testICalContentRegexpString := regexp.MustCompile(dateRegExpPattern).ReplaceAllString(testICalContent, dateRegExpPattern)
		testICalContentRegexpString = "^" + regexp.MustCompile("\\n").ReplaceAllString(testICalContentRegexpString, "[\\s\\n]*") + "$"

		re := regexp.MustCompile(testICalContentRegexpString)

		assert.True(t, re.MatchString(buffer.String()), "incorrect csv generation", testICalContentRegexpString, buffer.String())
	})

	t.Run("SetTask()", func(t *testing.T) {
		buffer := bytes.Buffer{}
		buffer.WriteString(testICalContent)
		taskRepository := NewICalTaskRepository(&buffer)

		err := taskRepository.SetTask(models.Task{
			Identifier: "first-task",
			Activity: []models.TaskActivity{
				{StartDate: tryParseTime("2010-06-01T15:00:00Z")},
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, otherTestICalContent, buffer.String(), "incorrect csv generation", otherTestICalContent, buffer.String())
	})

	t.Run("SetTask()", func(t *testing.T) {
		t.Run("3 tasks", func(t *testing.T) {
			buffer := bytes.Buffer{}
			taskRepository := NewICalTaskRepository(&buffer)

			tasks := models.Tasks{
				{Identifier: "first-task", Activity: []models.TaskActivity{
					{StartDate: tryParseTime("2010-06-01T15:00:00Z")},
				}},
				{Identifier: "second-task", Activity: []models.TaskActivity{
					{StartDate: tryParseTime("2019-01-01T10:00:00Z"), EndDate: tryParseTime("2019-01-01T10:04:00Z")},
				}},
				{Identifier: "last-task", Activity: []models.TaskActivity{
					{StartDate: tryParseTime("2019-01-01T10:00:00Z")},
				}},
			}

			err := taskRepository.SetTasks(tasks)
			assert.NoError(t, err)
			assert.Equal(t, otherTestICalContent, buffer.String(), "incorrect csv generation", otherTestICalContent, buffer.String())
		})

		t.Run("0 tasks", func(t *testing.T) {
			buffer := bytes.Buffer{}
			taskRepository := NewICalTaskRepository(&buffer)

			tasks := models.Tasks{}

			err := taskRepository.SetTasks(tasks)
			assert.NoError(t, err)
			assert.Equal(t, testICalEmptyContent, buffer.String(), "incorrect csv generation", testICalEmptyContent, buffer.String())
		})
	})
}
