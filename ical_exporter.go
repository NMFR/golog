package main

import "io"

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

func (iCalTaskExporter ICalTaskExporter) writeICalEvent(startTask Task, stopTask *Task, writer io.Writer) error {
	startTime := parseTime(startTask.getAt())

	if err := writeStrings(
		writer,
		"BEGIN:VEVENT\n",
		"SUMMARY:", startTask.getIdentifier(), "\n",
		"DTSTART:", startTime.UTC().Format(iCalTaskExporter.getTimeFormat()), "\n",
	); err != nil {
		return err
	}

	if stopTask != nil {
		stopTime := parseTime(stopTask.getAt())
		if err := writeStrings(
			writer,
			"DTEND:", stopTime.UTC().Format(iCalTaskExporter.getTimeFormat()), "\n",
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
func (iCalTaskExporter ICalTaskExporter) Export(tasks Tasks, writer io.Writer) error {
	var err error

	if err = iCalTaskExporter.writeICalHeader(writer); err != nil {
		return err
	}

	taskMap := map[string]Task{}
	for _, task := range tasks.Items {
		switch task.getAction() {
		case TaskStart:
			taskMap[task.getIdentifier()] = task
		case TaskStop:
			if startTask, inMap := taskMap[task.getIdentifier()]; inMap {
				err = iCalTaskExporter.writeICalEvent(startTask, &task, writer)
				if err != nil {
					return err
				}
				delete(taskMap, task.getIdentifier())
			}
		}
	}

	// Iterate running tasks:
	for _, startTask := range taskMap {
		err = iCalTaskExporter.writeICalEvent(startTask, nil, writer)
		if err != nil {
			return err
		}
	}

	err = iCalTaskExporter.writeICalFooter(writer)
	return err
}

// GetFileExtension returns the default file extension for ical formated files
func (iCalTaskExporter ICalTaskExporter) GetFileExtension() string {
	return "ics"
}
