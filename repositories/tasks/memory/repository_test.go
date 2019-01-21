package memory // import github.com/mlimaloureiro/golog/repositories/tasks/memory

import (
	"reflect"
	"testing"
	"time"

	tasksModel "github.com/mlimaloureiro/golog/models/tasks"
	tasksRepositories "github.com/mlimaloureiro/golog/repositories/tasks"

	"github.com/stretchr/testify/assert"
)

func timeFromString(str string) time.Time {
	date, _ := time.Parse(time.RFC3339, str)
	return date
}

func TestTaskRepository(t *testing.T) {
	t.Run("implements TaskRepositoryInterface", func(t *testing.T) {
		repository := New()
		assert.Implements(t, (*tasksRepositories.TaskRepositoryInterface)(nil), &repository)
	})

	t.Run("SetTasks() GetTasks()", func(t *testing.T) {
		caseTasks := tasksModel.Collection{
			tasksModel.Task{Identifier: "first-task", Activity: []tasksModel.TaskActivity{
				{StartDate: timeFromString("2017-01-01T19:06:02Z"), EndDate: timeFromString("2017-01-01T19:06:03Z")},
				{StartDate: timeFromString("2016-01-02T15:04:00Z")},
			}},
			tasksModel.Task{Identifier: "second-task", Activity: []tasksModel.TaskActivity{
				{StartDate: timeFromString("2016-01-02T15:04:00Z")},
			}},
		}

		repository := New()
		repository.SetTasks(caseTasks)
		repositoryTasks, err := repository.GetTasks()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(repositoryTasks))
		assert.True(t, reflect.DeepEqual(caseTasks, repositoryTasks))

		repositoryTasks[0].Identifier = "test"
		repositoryTasks, err = repository.GetTasks()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(repositoryTasks))
		assert.True(t, reflect.DeepEqual(caseTasks, repositoryTasks))
	})

	t.Run("SetTask() GetTask()", func(t *testing.T) {
		caseTasks := tasksModel.Collection{
			tasksModel.Task{Identifier: "first-task", Activity: []tasksModel.TaskActivity{
				{StartDate: timeFromString("2017-01-01T19:06:02Z"), EndDate: timeFromString("2017-01-01T19:06:03Z")},
				{StartDate: timeFromString("2016-01-02T15:04:00Z")},
			}},
			tasksModel.Task{Identifier: "second-task", Activity: []tasksModel.TaskActivity{
				{StartDate: timeFromString("2016-01-02T15:04:00Z")},
			}},
		}

		repository := New()
		repositoryTasks, err := repository.GetTasks()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(repositoryTasks))

		repositoryTask, err := repository.GetTask("unknown")
		assert.NoError(t, err)
		assert.Nil(t, repositoryTask)

		repository.SetTask(caseTasks[0])
		repositoryTasks, err = repository.GetTasks()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(repositoryTasks))

		repositoryTask, err = repository.GetTask("unknown")
		assert.NoError(t, err)
		assert.Nil(t, repositoryTask)

		repositoryTask, err = repository.GetTask(caseTasks[0].Identifier)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(caseTasks[0], *repositoryTask))

		repositoryTask.Identifier = "test"
		repositoryTask, err = repository.GetTask(caseTasks[0].Identifier)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(caseTasks[0], *repositoryTask))
	})
}
