package csv // import github.com/mlimaloureiro/golog/repositories/tasks/file/csv

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"testing"
	"time"

	tasksModel "github.com/mlimaloureiro/golog/models/tasks"
	tasksRepositories "github.com/mlimaloureiro/golog/repositories/tasks"

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
		setFileContent(t, "fixtures/test.csv", testCsvContent)
		repository := New("fixtures/test.csv")
		assert.Implements(t, (*tasksRepositories.TaskRepositoryInterface)(nil), repository)
	})

	t.Run("implements StarterPauserTaskRepositoryInterface", func(t *testing.T) {
		setFileContent(t, "fixtures/test.csv", testCsvContent)
		repository := New("fixtures/test.csv")
		assert.Implements(t, (*tasksRepositories.StarterPauserTaskRepositoryInterface)(nil), repository)
	})

	t.Run("GetTasks()", func(t *testing.T) {
		setFileContent(t, "fixtures/test.csv", testCsvContent)
		taskRepository := New("fixtures/test.csv")
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
			setFileContent(t, "fixtures/test.csv", testCsvContent)
			taskRepository := New("fixtures/test.csv")
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
			taskRepository := New("fixtures/test.csv")
			task, err := taskRepository.GetTask("unknown-task")

			assert.NoError(t, err)
			assert.Nil(t, task)
		})
	})

	t.Run("StartTask() PauseTask()", func(t *testing.T) {
		taskRepository := New("fixtures/write_test.csv")
		defer func() {
			os.Remove("fixtures/write_test.csv")
		}()

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

		fileContent := getFileContent(t, "fixtures/write_test.csv")
		assert.True(t, re.MatchString(fileContent), "incorrect csv generation", testCsvContentRegexpString, fileContent)
	})

	t.Run("SetTask()", func(t *testing.T) {
		setFileContent(t, "fixtures/write_test.csv", testCsvContent)
		taskRepository := New("fixtures/write_test.csv")
		defer func() {
			os.Remove("fixtures/write_test.csv")
		}()

		err := taskRepository.SetTask(tasksModel.Task{
			Identifier: "first-task",
			Activity: []tasksModel.TaskActivity{
				{StartDate: tryParseTime("2010-06-01T15:00:00Z")},
			},
		})
		assert.NoError(t, err)
		fileContent := getFileContent(t, "fixtures/write_test.csv")
		assert.Equal(t, otherTestCsvContent, fileContent, "incorrect csv generation", otherTestCsvContent, fileContent)
	})

	t.Run("SetTask()", func(t *testing.T) {
		t.Run("3 tasks", func(t *testing.T) {
			taskRepository := New("fixtures/write_test.csv")
			defer func() {
				os.Remove("fixtures/write_test.csv")
			}()

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
			fileContent := getFileContent(t, "fixtures/write_test.csv")
			assert.Equal(t, otherTestCsvContent, fileContent, "incorrect csv generation", otherTestCsvContent, fileContent)
		})

		t.Run("0 tasks", func(t *testing.T) {
			taskRepository := New("fixtures/write_test.csv")
			defer func() {
				os.Remove("fixtures/write_test.csv")
			}()

			tasks := tasksModel.Collection{}

			err := taskRepository.SetTasks(tasks)
			assert.NoError(t, err)
			fileContent := getFileContent(t, "fixtures/write_test.csv")
			assert.Equal(t, "", fileContent, "incorrect csv generation", "", fileContent)
		})
	})
}
