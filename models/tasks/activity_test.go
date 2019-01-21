package tasks // import github.com/mlimaloureiro/golog/models/tasks

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
