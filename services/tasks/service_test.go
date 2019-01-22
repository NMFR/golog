package tasks // import github.com/mlimaloureiro/golog/services/tasks

import (
	"testing"

	tasksModel "github.com/mlimaloureiro/golog/models/tasks"
	"github.com/mlimaloureiro/golog/repositories/tasks/memory"

	"github.com/stretchr/testify/assert"
)

func countRunningTasks(tasks tasksModel.Collection) int {
	counter := 0
	for i := range tasks {
		if tasks[i].IsRunning() {
			counter++
		}
	}
	return counter
}

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

	t.Run("SwitchTask()", func(t *testing.T) {
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
		runningTasksCounter := countRunningTasks(repositoryTasks)
		assert.Equal(t, 3, runningTasksCounter)

		err = service.SwitchTask("switched-task")
		assert.NoError(t, err)

		repositoryTasks, err = repository.GetTasks()
		assert.NoError(t, err)
		runningTasksCounter = countRunningTasks(repositoryTasks)
		assert.Equal(t, 1, runningTasksCounter)

		err = service.StartTask("first-task")
		assert.NoError(t, err)
		err = service.StartTask("second-task")
		assert.NoError(t, err)
		err = service.StartTask("last-task")
		assert.NoError(t, err)

		repositoryTasks, err = repository.GetTasks()
		assert.NoError(t, err)
		runningTasksCounter = countRunningTasks(repositoryTasks)
		assert.Equal(t, 4, runningTasksCounter)

		err = service.SwitchTask("first-task")
		assert.NoError(t, err)

		repositoryTasks, err = repository.GetTasks()
		assert.NoError(t, err)
		runningTasksCounter = countRunningTasks(repositoryTasks)
		assert.Equal(t, 1, runningTasksCounter)
	})
}
