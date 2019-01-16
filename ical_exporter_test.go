package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/mlimaloureiro/golog/models"
)

const (
	TestICalExportExpectedResult = `BEGIN:VCALENDAR
SUMMARY:2.0
PRODID:-//mlimaloureiro/golog
CALSCALE:GREGORIAN
BEGIN:VEVENT
SUMMARY:identifier-1
DTSTART:20160102T150400Z
DTEND:20160102T170402Z
END:VEVENT
BEGIN:VEVENT
SUMMARY:identifier-1
DTSTART:20170101T190602Z
DTEND:20170101T190603Z
END:VEVENT
BEGIN:VEVENT
SUMMARY:identifier-2
DTSTART:20160102T150400Z
DTEND:20160102T160400Z
END:VEVENT
END:VCALENDAR`
	TestICalExportRunningTasksExpectedResult = `BEGIN:VCALENDAR
SUMMARY:2.0
PRODID:-//mlimaloureiro/golog
CALSCALE:GREGORIAN
BEGIN:VEVENT
SUMMARY:identifier-1
DTSTART:20160102T150400Z
DTEND:20160102T170402Z
END:VEVENT
BEGIN:VEVENT
SUMMARY:identifier-2
DTSTART:20160102T150400Z
END:VEVENT
END:VCALENDAR`
)

func timeFromString(str string) time.Time {
	date, _ := time.Parse(time.RFC3339, str)
	return date
}

func TestICalExport(t *testing.T) {
	tasks := models.Tasks{
		{Identifier: "identifier-1", Activity: []models.TaskActivity{
			{StartDate: timeFromString("2016-01-02T15:04:00Z"), EndDate: timeFromString("2016-01-02T17:04:02Z")},
			{StartDate: timeFromString("2017-01-01T19:06:02Z"), EndDate: timeFromString("2017-01-01T19:06:03Z")},
		}},
		{Identifier: "identifier-2", Activity: []models.TaskActivity{
			{StartDate: timeFromString("2016-01-02T15:04:00Z"), EndDate: timeFromString("2016-01-02T16:04:00Z")},
		}},
	}
	buffer := bytes.Buffer{}
	exporter := ICalTaskExporter{}

	err := exporter.Export(tasks, &buffer)

	if err != nil {
		t.Error("Failed export.")
	}

	if buffer.String() != TestICalExportExpectedResult {
		t.Errorf("Exported result did not match expected result;\nExpected:\n%s\n\nReturned:\n%s", TestICalExportExpectedResult, buffer.String())
	}
}

func TestICalExportRunningTasks(t *testing.T) {
	tasks := models.Tasks{
		{Identifier: "identifier-1", Activity: []models.TaskActivity{
			{StartDate: timeFromString("2016-01-02T15:04:00Z"), EndDate: timeFromString("2016-01-02T17:04:02Z")},
		}},
		{Identifier: "identifier-2", Activity: []models.TaskActivity{
			{StartDate: timeFromString("2016-01-02T15:04:00Z")},
		}},
	}
	buffer := bytes.Buffer{}
	exporter := ICalTaskExporter{}

	err := exporter.Export(tasks, &buffer)

	if err != nil {
		t.Error("Failed export.")
	}

	if buffer.String() != TestICalExportRunningTasksExpectedResult {
		t.Errorf("Exported result did not match expected result;\nExpected:\n%s\n\nReturned:\n%s", TestICalExportRunningTasksExpectedResult, buffer.String())
	}
}
