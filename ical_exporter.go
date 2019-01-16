package main

import (
	"io"

	"github.com/mlimaloureiro/golog/models"
)

const prodID = "-//mlimaloureiro/golog"

// ICalTaskExporter is an implementation of TaskExporterInterface that exports Tasks in the ical (ics) format
type ICalTaskExporter struct {
	Version    string
	CALSCALE   string
	TimeFormat string // UTC time format string (layout string passed to time.Format)
}

func writeStrings(writer io.Writer, strings ...string) error {
	for _, str := range strings {
		_, err := io.WriteString(writer, str)
		if err != nil {
			return err
		}
	}
	return nil
}

func (iCalTaskExporter ICalTaskExporter) getVersion() string {
	if iCalTaskExporter.Version == "" {
		return "2.0"
	}
	return iCalTaskExporter.Version
}

func (iCalTaskExporter ICalTaskExporter) getCALSCALE() string {
	if iCalTaskExporter.CALSCALE == "" {
		return "GREGORIAN"
	}
	return iCalTaskExporter.CALSCALE
}

func (iCalTaskExporter ICalTaskExporter) getTimeFormat() string {
	if iCalTaskExporter.TimeFormat == "" {
		return "20060102T150405Z"
	}
	return iCalTaskExporter.TimeFormat
}

func (iCalTaskExporter ICalTaskExporter) writeICalHeader(writer io.Writer) error {
	err := writeStrings(
		writer,
		"BEGIN:VCALENDAR\n",
		"SUMMARY:", iCalTaskExporter.getVersion(), "\n",
		"PRODID:", prodID, "\n",
		"CALSCALE:", iCalTaskExporter.getCALSCALE(), "\n",
	)
	return err
}

func (iCalTaskExporter ICalTaskExporter) writeICalFooter(writer io.Writer) error {
	_, err := io.WriteString(writer, "END:VCALENDAR")
	return err
}

func (iCalTaskExporter ICalTaskExporter) writeICalEvent(task models.Task, taskActivity models.TaskActivity, writer io.Writer) error {
	if err := writeStrings(
		writer,
		"BEGIN:VEVENT\n",
		"SUMMARY:", task.Identifier, "\n",
		"DTSTART:", taskActivity.StartDate.UTC().Format(iCalTaskExporter.getTimeFormat()), "\n",
	); err != nil {
		return err
	}

	if taskActivity.IsRunning() == false {
		if err := writeStrings(
			writer,
			"DTEND:", taskActivity.EndDate.UTC().Format(iCalTaskExporter.getTimeFormat()), "\n",
		); err != nil {
			return err
		}
	}

	if err := writeStrings(writer, "END:VEVENT\n"); err != nil {
		return err
	}

	return nil
}

// Export tasks in the ical format to the writer
func (iCalTaskExporter ICalTaskExporter) Export(tasks models.Tasks, writer io.Writer) error {
	var err error

	if err = iCalTaskExporter.writeICalHeader(writer); err != nil {
		return err
	}

	for _, task := range tasks {
		for _, taskActivity := range task.Activity {
			if err = iCalTaskExporter.writeICalEvent(task, taskActivity, writer); err != nil {
				return err
			}
		}
	}

	err = iCalTaskExporter.writeICalFooter(writer)
	return err
}

// GetFileExtension returns the default file extension for ical formated files
func (iCalTaskExporter ICalTaskExporter) GetFileExtension() string {
	return "ics"
}
