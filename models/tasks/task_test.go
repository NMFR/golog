package tasks // import github.com/mlimaloureiro/golog/models/tasks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	t.Run("IsRunning()", func(t *testing.T) {
		t.Run("check on a running Task", func(t *testing.T) {
			task := Task{Identifier: "identifier-1", Activity: []TaskActivity{
				{StartDate: timeFromString("2017-01-01T19:06:02Z"), EndDate: timeFromString("2017-01-01T19:06:03Z")},
				{StartDate: timeFromString("2016-01-02T15:04:00Z")},
			}}

			assert.True(t, task.IsRunning())
		})

		t.Run("check on a Task that is not running", func(t *testing.T) {
			task := Task{Identifier: "identifier-1", Activity: []TaskActivity{
				{StartDate: timeFromString("2017-01-01T19:06:02Z"), EndDate: timeFromString("2017-01-01T19:06:03Z")},
				{StartDate: timeFromString("2016-01-02T15:04:00Z"), EndDate: timeFromString("2016-01-02T17:04:02Z")},
			}}

			assert.False(t, task.IsRunning())
		})
	})

	t.Run("GetDuration()", func(t *testing.T) {
		task := Task{Identifier: "identifier-1", Activity: []TaskActivity{
			{StartDate: timeFromString("2017-01-01T19:06:02Z"), EndDate: timeFromString("2017-01-01T19:06:03Z")},
			{StartDate: timeFromString("2016-01-02T15:04:00Z"), EndDate: timeFromString("2016-01-02T15:04:02Z")},
		}}

		assert.Equal(t, 3*time.Second, task.GetDuration())
	})

	t.Run("GetRunningTaskActivity()", func(t *testing.T) {
		t.Run("check on a running Task", func(t *testing.T) {
			task := Task{Identifier: "identifier-1", Activity: []TaskActivity{
				{StartDate: timeFromString("2017-01-01T19:06:02Z"), EndDate: timeFromString("2017-01-01T19:06:03Z")},
				{StartDate: timeFromString("2016-01-02T15:04:00Z")},
			}}

			assert.Equal(t, &task.Activity[1], task.GetRunningTaskActivity())
		})

		t.Run("check on a Task that is not running", func(t *testing.T) {
			task := Task{Identifier: "identifier-1", Activity: []TaskActivity{
				{StartDate: timeFromString("2017-01-01T19:06:02Z"), EndDate: timeFromString("2017-01-01T19:06:03Z")},
				{StartDate: timeFromString("2016-01-02T15:04:00Z"), EndDate: timeFromString("2016-01-02T17:04:02Z")},
			}}

			assert.Nil(t, task.GetRunningTaskActivity())
		})
	})
}
