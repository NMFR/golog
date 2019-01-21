package tasks // import github.com/mlimaloureiro/golog/models/tasks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTasks(t *testing.T) {
	t.Run("GetByIdentifier()", func(t *testing.T) {
		tasks := Collection{
			{Identifier: "identifier-1", Activity: []TaskActivity{
				{StartDate: timeFromString("2016-01-02T15:04:00Z"), EndDate: timeFromString("2016-01-02T17:04:02Z")},
			}},
			{Identifier: "identifier-2", Activity: []TaskActivity{
				{StartDate: timeFromString("2016-01-02T15:04:00Z")},
			}},
		}

		assert.NotNil(t, tasks.GetByIdentifier("identifier-1"))
		assert.NotNil(t, tasks.GetByIdentifier("identifier-2"))
		assert.Nil(t, tasks.GetByIdentifier("unknown"))
	})
}
