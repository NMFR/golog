package main

import (
	"fmt"
	"math"

	taskModel "github.com/mlimaloureiro/golog/models/tasks"
)

// Transformer is a type that has loaded all Tasks entries from storage
type Transformer struct {
	LoadedTasks taskModel.Collection
}

// Transform Transforms all tasks to human readable
func (transformer *Transformer) Transform() map[string]string {
	transformedTasks := map[string]string{}

	tasks := transformer.LoadedTasks
	for _, task := range tasks {
		isActive := task.IsRunning()
		taskSeconds := transformer.TrackingToSeconds(task)
		humanTime := transformer.SecondsToHuman(taskSeconds)

		status := ""
		if isActive {
			status = "(running)"
		}

		transformedTask := fmt.Sprintf("%s    %s %s", humanTime, task.Identifier, status)
		transformedTasks[task.Identifier] = transformedTask
	}

	return transformedTasks
}

// SecondsToHuman returns an human readable string from seconds
func (transformer *Transformer) SecondsToHuman(totalSeconds int) string {
	hours := math.Floor(float64(((totalSeconds % 31536000) % 86400) / 3600))
	minutes := math.Floor(float64((((totalSeconds % 31536000) % 86400) % 3600) / 60))
	seconds := (((totalSeconds % 31536000) % 86400) % 3600) % 60

	return fmt.Sprintf("%dh:%dm:%ds", int(hours), int(minutes), int(seconds))
}

// TrackingToSeconds get entries from storage by identifier and calculate
// time between each start/stop for a single identifier
func (transformer *Transformer) TrackingToSeconds(task taskModel.Task) int {
	return (int)(task.GetDuration().Seconds())
}
