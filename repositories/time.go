package repositories // import github.com/mlimaloureiro/golog/repositories

import "time"

func formatTime(at time.Time) string {
	return at.Format(time.RFC3339)
}

func parseTime(at string) (time.Time, error) {
	then, err := time.Parse(time.RFC3339, at)
	return then, err
}

func tryParseTime(str string) time.Time {
	date, _ := parseTime(str)
	return date
}
