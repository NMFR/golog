package models // import github.com/mlimaloureiro/golog/models

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func timeFromString(str string) time.Time {
	date, _ := time.Parse(time.RFC3339, str)
	return date
}

func TestTaskActivity(t *testing.T) {
	t.Run("IsRunning()", func(t *testing.T) {
		t.Run("running", func(t *testing.T) {
			taskActivity := TaskActivity{StartDate: timeFromString("2019-01-01T10:00:00Z")}
			assert.True(t, taskActivity.IsRunning())
		})

		t.Run("not running", func(t *testing.T) {
			taskActivity := TaskActivity{StartDate: timeFromString("2019-01-01T10:00:00Z"), EndDate: timeFromString("2019-01-02T10:00:00Z")}
			assert.False(t, taskActivity.IsRunning())
		})
	})

	t.Run("Duration()", func(t *testing.T) {
		t.Run("running", func(t *testing.T) {
			taskActivity := TaskActivity{StartDate: timeFromString("2019-01-01T10:00:01Z")}
			timeNow := time.Now()
			taskDuration := taskActivity.Duration()
			expectedTaskDuration := timeNow.Sub(timeFromString("2019-01-01T10:00:01Z"))

			if diff := time.Duration(math.Abs(float64(expectedTaskDuration - taskDuration))); diff < (1 * time.Second) {
				expectedTaskDuration = taskDuration
			}
			assert.Equal(t, expectedTaskDuration, taskDuration)
		})

		t.Run("not running", func(t *testing.T) {
			taskActivity := TaskActivity{StartDate: timeFromString("2019-01-01T10:00:01Z"), EndDate: timeFromString("2019-01-02T10:00:00Z")}
			assert.Equal(t, 24*time.Hour-1*time.Second, taskActivity.Duration())
		})
	})
}

func TestTask(t *testing.T) {
	t.Run("IsRunning()", func(t *testing.T) {
		t.Run("running", func(t *testing.T) {
			task := Task{Identifier: "identifier-1", Activity: []TaskActivity{
				{StartDate: timeFromString("2017-01-01T19:06:02Z"), EndDate: timeFromString("2017-01-01T19:06:03Z")},
				{StartDate: timeFromString("2016-01-02T15:04:00Z")},
			}}

			assert.True(t, task.IsRunning())
		})

		t.Run("not running", func(t *testing.T) {
			task := Task{Identifier: "identifier-1", Activity: []TaskActivity{
				{StartDate: timeFromString("2017-01-01T19:06:02Z"), EndDate: timeFromString("2017-01-01T19:06:03Z")},
				{StartDate: timeFromString("2016-01-02T15:04:00Z"), EndDate: timeFromString("2016-01-02T17:04:02Z")},
			}}

			assert.False(t, task.IsRunning())
		})
	})

	t.Run("Duration()", func(t *testing.T) {
		task := Task{Identifier: "identifier-1", Activity: []TaskActivity{
			{StartDate: timeFromString("2017-01-01T19:06:02Z"), EndDate: timeFromString("2017-01-01T19:06:03Z")},
			{StartDate: timeFromString("2016-01-02T15:04:00Z"), EndDate: timeFromString("2016-01-02T15:04:02Z")},
		}}

		assert.Equal(t, 3*time.Second, task.Duration())
	})

	t.Run("GetRunningTaskActivity()", func(t *testing.T) {
		t.Run("running", func(t *testing.T) {
			task := Task{Identifier: "identifier-1", Activity: []TaskActivity{
				{StartDate: timeFromString("2017-01-01T19:06:02Z"), EndDate: timeFromString("2017-01-01T19:06:03Z")},
				{StartDate: timeFromString("2016-01-02T15:04:00Z")},
			}}

			assert.Equal(t, &task.Activity[1], task.GetRunningTaskActivity())
		})

		t.Run("not running", func(t *testing.T) {
			task := Task{Identifier: "identifier-1", Activity: []TaskActivity{
				{StartDate: timeFromString("2017-01-01T19:06:02Z"), EndDate: timeFromString("2017-01-01T19:06:03Z")},
				{StartDate: timeFromString("2016-01-02T15:04:00Z"), EndDate: timeFromString("2016-01-02T17:04:02Z")},
			}}

			assert.Nil(t, task.GetRunningTaskActivity())
		})
	})
}

func TestTasks(t *testing.T) {
	t.Run("GetByIdentifier()", func(t *testing.T) {
		tasks := Tasks{
			{Identifier: "identifier-1", Activity: []TaskActivity{
				{StartDate: timeFromString("2016-01-02T15:04:00Z"), EndDate: timeFromString("2016-01-02T17:04:02Z")},
			}},
			{Identifier: "identifier-2", Activity: []TaskActivity{
				{StartDate: timeFromString("2016-01-02T15:04:00Z")},
			}},
		}

		assert.NotNil(t, tasks.GetByIdentifier("identifier-1"))
		assert.NotNil(t, tasks.GetByIdentifier("identifier-2"))
		assert.Nil(t, tasks.GetByIdentifier("unknow"))
	})
}
