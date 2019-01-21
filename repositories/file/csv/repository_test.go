package csv // import github.com/mlimaloureiro/golog/repositories/file/csv

import (
	"bytes"
	"fmt"
	"math"
	"regexp"
	"testing"
	"time"

	tasksModel "github.com/mlimaloureiro/golog/models/tasks"
	"github.com/mlimaloureiro/golog/repositories"

	"github.com/stretchr/testify/assert"
)

const (
	testCsvContent = `first-task,start,2019-01-01T10:00:00Z
first-task,stop,2019-01-01T10:01:00Z
second-task,start,2019-01-01T10:00:00Z
second-task,stop,2019-01-01T10:04:00Z
first-task,start,2019-01-01T10:02:00Z
first-task,stop,2019-01-01T10:06:00Z
last-task,start,2019-01-01T10:00:00Z
`
	otherTestCsvContent = `first-task,start,2010-06-01T15:00:00Z
second-task,start,2019-01-01T10:00:00Z
second-task,stop,2019-01-01T10:04:00Z
last-task,start,2019-01-01T10:00:00Z
`
)

func TestTaskRepository(t *testing.T) {
	t.Run("implements TaskRepositoryInterface", func(t *testing.T) {
		buffer := bytes.Buffer{}
		repository := NewTaskRepository(&buffer)
		assert.Implements(t, (*repositories.TaskRepositoryInterface)(nil), repository)
	})

	t.Run("GetTasks()", func(t *testing.T) {
		buffer := bytes.Buffer{}
		buffer.WriteString(testCsvContent)
		taskRepository := NewTaskRepository(&buffer)
		tasks, err := taskRepository.GetTasks()

		assert.NoError(t, err)
		assert.Equal(t, 3, len(tasks))

		taskMap := make(map[string]tasksModel.Task)
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
			buffer.WriteString(testCsvContent)
			taskRepository := NewTaskRepository(&buffer)
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
			buffer.WriteString(testCsvContent)
			taskRepository := NewTaskRepository(&buffer)
			task, err := taskRepository.GetTask("unknown-task")

			assert.NoError(t, err)
			assert.Nil(t, task)
		})
	})

	t.Run("StartTask() PauseTask()", func(t *testing.T) {
		buffer := bytes.Buffer{}
		taskRepository := NewTaskRepository(&buffer)

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

		dateRegExpPattern := "\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z"
		testCsvContentRegexpString := regexp.MustCompile(dateRegExpPattern).ReplaceAllString(testCsvContent, dateRegExpPattern)
		testCsvContentRegexpString = "^" + regexp.MustCompile("\\n").ReplaceAllString(testCsvContentRegexpString, "[\\s\\n]*") + "$"

		re := regexp.MustCompile(testCsvContentRegexpString)

		assert.True(t, re.MatchString(buffer.String()), "incorrect csv generation", testCsvContentRegexpString, buffer.String())
	})

	t.Run("SetTask()", func(t *testing.T) {
		buffer := bytes.Buffer{}
		buffer.WriteString(testCsvContent)
		taskRepository := NewTaskRepository(&buffer)

		err := taskRepository.SetTask(tasksModel.Task{
			Identifier: "first-task",
			Activity: []tasksModel.TaskActivity{
				{StartDate: tryParseTime("2010-06-01T15:00:00Z")},
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, otherTestCsvContent, buffer.String(), "incorrect csv generation", otherTestCsvContent, buffer.String())
	})

	t.Run("SetTask()", func(t *testing.T) {
		t.Run("3 tasks", func(t *testing.T) {
			buffer := bytes.Buffer{}
			taskRepository := NewTaskRepository(&buffer)

			tasks := tasksModel.Collection{
				{Identifier: "first-task", Activity: []tasksModel.TaskActivity{
					{StartDate: tryParseTime("2010-06-01T15:00:00Z")},
				}},
				{Identifier: "second-task", Activity: []tasksModel.TaskActivity{
					{StartDate: tryParseTime("2019-01-01T10:00:00Z"), EndDate: tryParseTime("2019-01-01T10:04:00Z")},
				}},
				{Identifier: "last-task", Activity: []tasksModel.TaskActivity{
					{StartDate: tryParseTime("2019-01-01T10:00:00Z")},
				}},
			}

			err := taskRepository.SetTasks(tasks)
			assert.NoError(t, err)
			assert.Equal(t, otherTestCsvContent, buffer.String(), "incorrect csv generation", otherTestCsvContent, buffer.String())
		})

		t.Run("0 tasks", func(t *testing.T) {
			buffer := bytes.Buffer{}
			taskRepository := NewTaskRepository(&buffer)

			tasks := tasksModel.Collection{}

			err := taskRepository.SetTasks(tasks)
			assert.NoError(t, err)
			assert.Equal(t, "", buffer.String(), "incorrect csv generation", "", buffer.String())
		})
	})
}
