package main

import (
	"testing"
	"time"

	"github.com/mlimaloureiro/golog/models"
)

const hourInSeconds = 3600
const hourInMinutes = 60

func timeFromString(str string) time.Time {
	date, _ := time.Parse(time.RFC3339, str)
	return date
}

func TestTransform(t *testing.T) {
	tasks := models.Tasks{
		{Identifier: "identifier-1", Activity: []models.TaskActivity{
			{StartDate: timeFromString("2016-01-02T15:04:00Z"), EndDate: timeFromString("2016-01-02T17:04:02Z")},
			{StartDate: timeFromString("2016-12-29T19:04:00Z"), EndDate: timeFromString("2016-12-29T19:06:02Z")},
		}},
		{Identifier: "identifier-2", Activity: []models.TaskActivity{
			{StartDate: timeFromString("2016-01-02T15:04:00Z"), EndDate: timeFromString("2016-01-02T16:04:00Z")},
		}},
	}
	transformer := Transformer{LoadedTasks: tasks}
	transformedTasks := transformer.Transform()
	expectedString1 := "2h:2m:4s    identifier-1 "
	expectedString2 := "1h:0m:0s    identifier-2 "
	if transformedTasks["identifier-1"] != expectedString1 {
		t.Errorf("Expected %s, got %s.", expectedString1, transformedTasks["identifier-1"])
	}
	if transformedTasks["identifier-2"] != expectedString2 {
		t.Errorf("Expected %s, got %s.", expectedString2, transformedTasks["identifier-2"])
	}
}

func TestSecondsToHuman(t *testing.T) {
	transformer := Transformer{}
	secondsCase1 := 1432
	secondsCase2 := 4432
	if transformer.SecondsToHuman(secondsCase1) != "0h:23m:52s" {
		t.Errorf(
			"1432 Seconds to human should be 0h:23m:52s, got %s.",
			transformer.SecondsToHuman(secondsCase1),
		)
	}
	if transformer.SecondsToHuman(secondsCase2) != "1h:13m:52s" {
		t.Errorf(
			"4432 Seconds to human should be 1h:13m:52s, got %s.",
			transformer.SecondsToHuman(secondsCase2),
		)
	}
}

func TestTrackingToSeconds(t *testing.T) {
	tasks := models.Tasks{
		{Identifier: "identifier-1", Activity: []models.TaskActivity{
			{StartDate: timeFromString("2016-01-02T15:04:00Z"), EndDate: timeFromString("2016-01-02T17:04:02Z")},
			{StartDate: timeFromString("2016-12-29T19:04:00Z"), EndDate: timeFromString("2016-12-29T19:06:02Z")},
			{StartDate: timeFromString("2017-01-01T19:06:02Z"), EndDate: timeFromString("2017-01-01T19:06:03Z")},
		}},
		{Identifier: "identifier-2", Activity: []models.TaskActivity{
			{StartDate: timeFromString("2016-01-02T15:04:00Z"), EndDate: timeFromString("2016-01-02T16:04:00Z")},
		}},
	}
	transformer := Transformer{LoadedTasks: tasks}
	// @todo test status
	seconds := transformer.TrackingToSeconds(tasks[0])
	if seconds != hourInSeconds*2+hourInMinutes*2+5 {
		t.Errorf(
			"Transformation for identifier-1 should be 7325 seconds, got %d.",
			seconds,
		)
	}
	seconds = transformer.TrackingToSeconds(tasks[1])
	if seconds != hourInSeconds*1 {
		t.Errorf(
			"Transformation for identifier-1 should be 3600 seconds, got %d.",
			seconds,
		)
	}
}
