package ical // import github.com/mlimaloureiro/golog/repositories/tasks/file/ical

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"testing"
	"time"

	taskModel "github.com/mlimaloureiro/golog/models/tasks"
	taskRepositories "github.com/mlimaloureiro/golog/repositories/tasks"

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

func parseTime(at string) (time.Time, error) {
	then, err := time.Parse(time.RFC3339, at)
	return then, err
}

func tryParseTime(str string) time.Time {
	date, _ := parseTime(str)
	return date
}

func getFileContent(t *testing.T, filePath string) string {
	stringBytes, err := ioutil.ReadFile(filePath)
	assert.NoError(t, err)
	return string(stringBytes)
}

func setFileContent(t *testing.T, filePath string, content string) {
	file, err := os.Create(filePath)
	assert.NoError(t, err)
	_, err = file.WriteString(content)
	assert.NoError(t, err)
	err = file.Close()
	assert.NoError(t, err)
}

func TestTaskRepository(t *testing.T) {
	t.Run("implements TaskRepositoryInterface", func(t *testing.T) {
		setFileContent(t, "fixtures/test.ics", testICalContent)
		repository := New("fixtures/test.ics")
		assert.Implements(t, (*taskRepositories.TaskRepositoryInterface)(nil), repository)
	})

	t.Run("GetTasks()", func(t *testing.T) {
		setFileContent(t, "fixtures/test.ics", testICalContent)
		taskRepository := New("fixtures/test.ics")
		tasks, err := taskRepository.GetTasks()

		assert.NoError(t, err)
		assert.Equal(t, 3, len(tasks))

		taskMap := make(map[string]taskModel.Task)
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
				taskDuration := task.GetDuration()
				if task.IsRunning() {
					if diff := time.Duration(math.Abs(float64(taskCase.duration - taskDuration))); diff < (1 * time.Second) {
						caseDuration = taskDuration
					}
				}
				assert.Equal(t, caseDuration, taskDuration, "task.GetDuration()")
			})
		}
	})

	t.Run("GetTask()", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			setFileContent(t, "fixtures/test.ics", testICalContent)
			taskRepository := New("fixtures/test.ics")
			task, err := taskRepository.GetTask("first-task")

			assert.NoError(t, err)
			if assert.NotNil(t, task) == false {
				return
			}

			assert.Equal(t, "first-task", task.Identifier, "task.Identifier")
			assert.Equal(t, false, task.IsRunning(), "task.IsRunning()")
			assert.Equal(t, 5*time.Minute, task.GetDuration(), "task.GetDuration()")
		})

		t.Run("success (unknown task)", func(t *testing.T) {
			taskRepository := New("fixtures/test.ics")
			task, err := taskRepository.GetTask("unknown-task")

			assert.NoError(t, err)
			assert.Nil(t, task)
		})
	})

	t.Run("SetTask()", func(t *testing.T) {
		setFileContent(t, "fixtures/write_test.ics", testICalContent)
		taskRepository := New("fixtures/write_test.ics")
		defer func() {
			os.Remove("fixtures/write_test.ics")
		}()

		err := taskRepository.SetTask(taskModel.Task{
			Identifier: "first-task",
			Activity: []taskModel.TaskActivity{
				{StartDate: tryParseTime("2010-06-01T15:00:00Z")},
			},
		})
		assert.NoError(t, err)
		fileContent := getFileContent(t, "fixtures/write_test.ics")
		assert.Equal(t, otherTestICalContent, fileContent, "incorrect ics generation", otherTestICalContent, fileContent)
	})

	t.Run("SetTask()", func(t *testing.T) {
		t.Run("check repository with tasks", func(t *testing.T) {
			setFileContent(t, "fixtures/write_test.ics", testICalContent)
			taskRepository := New("fixtures/write_test.ics")
			defer func() {
				os.Remove("fixtures/write_test.ics")
			}()

			tasks := taskModel.Collection{
				{Identifier: "first-task", Activity: []taskModel.TaskActivity{
					{StartDate: tryParseTime("2010-06-01T15:00:00Z")},
				}},
				{Identifier: "second-task", Activity: []taskModel.TaskActivity{
					{StartDate: tryParseTime("2019-01-01T10:00:00Z"), EndDate: tryParseTime("2019-01-01T10:04:00Z")},
				}},
				{Identifier: "last-task", Activity: []taskModel.TaskActivity{
					{StartDate: tryParseTime("2019-01-01T10:00:00Z")},
				}},
			}

			err := taskRepository.SetTasks(tasks)
			assert.NoError(t, err)
			fileContent := getFileContent(t, "fixtures/write_test.ics")
			assert.Equal(t, otherTestICalContent, fileContent, "incorrect ics generation", otherTestICalContent, fileContent)
		})

		t.Run("check empty repository", func(t *testing.T) {
			setFileContent(t, "fixtures/write_test.ics", testICalContent)
			taskRepository := New("fixtures/write_test.ics")
			defer func() {
				os.Remove("fixtures/write_test.ics")
			}()

			tasks := taskModel.Collection{}

			err := taskRepository.SetTasks(tasks)
			assert.NoError(t, err)
			fileContent := getFileContent(t, "fixtures/write_test.ics")
			assert.Equal(t, testICalEmptyContent, fileContent, "incorrect ics generation", testICalEmptyContent, fileContent)
		})
	})
}
