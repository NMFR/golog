package main

import "testing"

func TestCsvRepositoryLoad(t *testing.T) {
	taskCsvRepository := TaskCsvRepository{Path: "fixtures/test_load.csv"}
	tasks, _ := taskCsvRepository.load()
	if !(len(tasks.Items) == 2) {
		t.Error("Expected to have 2 items in csv file.")
	}

	expectedTasks := Tasks{}
	expectedTasks.addItem(Task{Identifier: "track-search", Action: TaskStart, At: "2016-01-02T15:04:00Z"})
	expectedTasks.addItem(Task{Identifier: "live-demo", Action: TaskStop, At: "2016-01-02T15:04:00Z"})

	if expectedTasks.Items[0] != tasks.Items[0] || expectedTasks.Items[1] != tasks.Items[1] {
		t.Error("Tasks loaded are different from tasks expected.")
	}
}

func TestCsvRepositorySave(t *testing.T) {
	taskCsvRepository := TaskCsvRepository{Path: "fixtures/test_save.csv"}
	tasks, err := taskCsvRepository.load()
	if !(len(tasks.Items) == 0) || err != nil {
		t.Error("Expected to have 0 items in csv file.")
	}

	taskCsvRepository.save(Task{Identifier: "identifier-1", Action: TaskStart, At: "2016-01-02T15:04:00Z"})
	taskCsvRepository.save(Task{Identifier: "identifier-2", Action: TaskStop, At: "2016-01-02T15:04:00Z"})
	tasks, _ = taskCsvRepository.load()
	if !(len(tasks.Items) == 2) {
		t.Error("Expected to have 2 items in csv file.")
	}
	if tasks.Items[1].getIdentifier() != "identifier-2" {
		t.Error("Last line should be with with identifier-2 task.")
	}
	taskCsvRepository.clear()
}
