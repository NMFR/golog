package main

import "io"

// TaskExporterInterface interface is used to export Tasks in a specifict format to a given writable stream,
// Format examples: csv, ical, xml, ...
type TaskExporterInterface interface {
	Export(tasks Tasks, writer io.Writer) error
	GetFileExtension() string
}
