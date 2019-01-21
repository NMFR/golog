package tasks // import github.com/mlimaloureiro/golog/services/tasks

import (
	"testing"

	"github.com/mlimaloureiro/golog/repositories/tasks/memory"

	"github.com/stretchr/testify/assert"
)

func TestTaskService(t *testing.T) {
	t.Run("StartTask() PauseTask()", func(t *testing.T) {
		repository := memory.New()
		service := New(&repository)

		err := service.StartTask("first-task")
		assert.NoError(t, err)
		repositoryTasks, err := repository.GetTasks()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(repositoryTasks))
		repositoryTask, err := repository.GetTask("first-task")
		assert.NoError(t, err)
		assert.NotNil(t, repositoryTask)
		assert.True(t, repositoryTask.IsRunning())

		err = service.PauseTask("first-task")
		assert.NoError(t, err)
		repositoryTasks, err = repository.GetTasks()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(repositoryTasks))
		repositoryTask, err = repository.GetTask("first-task")
		assert.NoError(t, err)
		assert.NotNil(t, repositoryTask)
		assert.False(t, repositoryTask.IsRunning())
	})

	t.Run("DeleteTask() DeleteTasks()", func(t *testing.T) {
		repository := memory.New()
		service := New(&repository)

		err := service.StartTask("first-task")
		assert.NoError(t, err)
		err = service.StartTask("second-task")
		assert.NoError(t, err)
		err = service.StartTask("last-task")
		assert.NoError(t, err)
		repositoryTasks, err := repository.GetTasks()
		assert.NoError(t, err)
		assert.Equal(t, 3, len(repositoryTasks))

		repositoryTask, err := repository.GetTask("second-task")
		assert.NoError(t, err)
		assert.NotNil(t, repositoryTask)

		err = service.DeleteTask("second-task")
		assert.NoError(t, err)
		repositoryTasks, err = repository.GetTasks()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(repositoryTasks))

		repositoryTask, err = repository.GetTask("second-task")
		assert.NoError(t, err)
		assert.Nil(t, repositoryTask)

		err = service.DeleteTasks()
		assert.NoError(t, err)
		repositoryTasks, err = repository.GetTasks()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(repositoryTasks))
	})
}
