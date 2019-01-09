package main

import (
	"fmt"
	"math"
	"os"
	"time"
)

// Transformer is a type that has loaded all Tasks entries from storage
type Transformer struct {
	LoadedTasks Tasks
}

// Transform Transforms all tasks to human readable
func (transformer *Transformer) Transform() map[string]string {
	transformedTasks := map[string]string{}
	tasks := transformer.LoadedTasks.Items
	for _, task := range tasks {
		if _, inMap := transformedTasks[task.getIdentifier()]; inMap {
			continue
		}
		taskSeconds, isActive := transformer.TrackingToSeconds(task.getIdentifier())
		humanTime := transformer.SecondsToHuman(taskSeconds)

		status := ""
		if isActive {
			status = "(running)"
		}

		transformedTask := fmt.Sprintf("%s    %s %s", humanTime, task.getIdentifier(), status)
		transformedTasks[task.getIdentifier()] = transformedTask
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
func (transformer *Transformer) TrackingToSeconds(identifier string) (int, bool) {
	nextAction := TaskStart
	var durationInSeconds float64
	var startTime, stopTime time.Time

	tasks := transformer.LoadedTasks.getByIdentifier(identifier)
	for _, task := range tasks.Items {
		if task.getAction() == TaskStart && nextAction == TaskStart {
			nextAction = TaskStop
			startTime = parseTime(task.getAt())
		}
		if task.getAction() == TaskStop && nextAction == TaskStop {
			nextAction = TaskStart
			stopTime = parseTime(task.getAt())
			durationInSeconds += stopTime.Sub(startTime).Seconds()
		}
	}

	if isActive(nextAction) {
		durationInSeconds += time.Since(startTime).Seconds()
	}

	return int(durationInSeconds), isActive(nextAction)
}

// we can check if a task is active if we reach the end of the loop
// without finding the last stop action
func isActive(nextAction string) bool {
	return nextAction == TaskStop
}

func parseTime(at string) time.Time {
	then, err := time.Parse(time.RFC3339, at)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return then
}
