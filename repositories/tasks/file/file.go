package file // import github.com/mlimaloureiro/golog/repositories/tasks/file

import (
	"fmt"

	tasksRepositories "github.com/mlimaloureiro/golog/repositories/tasks"
	"github.com/mlimaloureiro/golog/repositories/tasks/file/csv"
	"github.com/mlimaloureiro/golog/repositories/tasks/file/ical"
)

// Type of file formats
const (
	CSV  Format = "csv"
	ICAL Format = "ical"
	ICS  Format = "ics"
)

var formatDefaultExtensionMap = map[Format]string{
	CSV:  ".csv",
	ICAL: ".ics",
	ICS:  ".ics",
}

// Format is a enum that contains the supported task repository file formats
type Format string

// GetFormats returns all supported Formats
func GetFormats() []Format {
	return []Format{CSV, ICAL, ICS}
}

// GetFormatNames returns the name of all supported Formats
func GetFormatNames() []string {
	formats := GetFormats()
	strs := make([]string, len(formats))
	for i := range formats {
		strs[i] = string(formats[i])
	}
	return strs
}

// GetTaskFileRepository returns the repository of the format that will presist in the file present at filePath
func GetTaskFileRepository(format Format, filePath string) (tasksRepositories.TaskRepositoryInterface, error) {
	var repository tasksRepositories.TaskRepositoryInterface
	switch format {
	case CSV:
		repository = csv.New(filePath)
	case ICAL:
		fallthrough
	case ICS:
		repository = ical.New(filePath)
	}

	if repository == nil {
		return nil, fmt.Errorf("invalid format \"%s\"", format)
	}
	return repository, nil
}

// GetFormatFileExtension returns the file extension corresponding to the supplied format
func GetFormatFileExtension(format Format) string {
	return formatDefaultExtensionMap[format]
}
