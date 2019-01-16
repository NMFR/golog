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

const testCsvContent = `first-task,start,2019-01-01T10:00:00Z
first-task,stop,2019-01-01T10:01:00Z
second-task,start,2019-01-01T10:00:00Z
second-task,stop,2019-01-01T10:04:00Z
first-task,start,2019-01-01T10:02:00Z
first-task,stop,2019-01-01T10:06:00Z
last-task,start,2019-01-01T10:00:00Z
`

func timeFromString(str string) time.Time {
	date, _ := parseTime(str)
	return date
}

func TestCsvTaskRepository(t *testing.T) {
	t.Run("implements TaskRepositoryInterface", func(t *testing.T) {
		buffer := bytes.Buffer{}
		repository := NewCsvTaskRepository(&buffer)
		assert.Implements(t, (*models.TaskRepositoryInterface)(nil), repository)
	})

	t.Run("GetTasks()", func(t *testing.T) {
		csvBuffer := bytes.Buffer{}
		csvBuffer.WriteString(testCsvContent)
		csvTaskRepository := NewCsvTaskRepository(&csvBuffer)
		tasks, err := csvTaskRepository.GetTasks()

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
			{"last-task", true, time.Now().Sub(timeFromString("2019-01-01T10:00:00Z"))},
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
			csvBuffer := bytes.Buffer{}
			csvBuffer.WriteString(testCsvContent)
			csvTaskRepository := NewCsvTaskRepository(&csvBuffer)
			task, err := csvTaskRepository.GetTask("first-task")

			assert.NoError(t, err)
			assert.NotNil(t, task)

			assert.Equal(t, "first-task", task.Identifier, "task.Identifier")
			assert.Equal(t, false, task.IsRunning(), "task.IsRunning()")
			assert.Equal(t, 5*time.Minute, task.Duration(), "task.Duration()")
		})

		t.Run("success (unknown task)", func(t *testing.T) {
			csvBuffer := bytes.Buffer{}
			csvBuffer.WriteString(testCsvContent)
			csvTaskRepository := NewCsvTaskRepository(&csvBuffer)
			task, err := csvTaskRepository.GetTask("unknown-task")

			assert.NoError(t, err)
			assert.Nil(t, task)
		})
	})

	t.Run("StartTask() PauseTask()", func(t *testing.T) {
		csvBuffer := bytes.Buffer{}
		csvTaskRepository := NewCsvTaskRepository(&csvBuffer)

		err := csvTaskRepository.StartTask("first-task")
		assert.NoError(t, err)
		err = csvTaskRepository.PauseTask("first-task")
		assert.NoError(t, err)

		err = csvTaskRepository.StartTask("second-task")
		assert.NoError(t, err)
		err = csvTaskRepository.PauseTask("second-task")
		assert.NoError(t, err)

		err = csvTaskRepository.StartTask("first-task")
		assert.NoError(t, err)
		err = csvTaskRepository.PauseTask("first-task")
		assert.NoError(t, err)

		err = csvTaskRepository.StartTask("last-task")
		assert.NoError(t, err)

		dateRegExpPattern := "\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z"
		testCsvContentRegexpString := regexp.MustCompile(dateRegExpPattern).ReplaceAllString(testCsvContent, dateRegExpPattern)
		testCsvContentRegexpString = "^" + regexp.MustCompile("\\n").ReplaceAllString(testCsvContentRegexpString, "[\\s\\n]*") + "$"

		re := regexp.MustCompile(testCsvContentRegexpString)

		assert.True(t, re.MatchString(csvBuffer.String()), "incorrect csv generation", testCsvContentRegexpString, csvBuffer.String())
	})
}
