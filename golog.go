package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	tasksModel "github.com/mlimaloureiro/golog/models/tasks"
	"github.com/mlimaloureiro/golog/repositories/tasks/file"
	"github.com/mlimaloureiro/golog/repositories/tasks/file/csv"
	tasksServices "github.com/mlimaloureiro/golog/services/tasks"

	"github.com/codegangsta/cli"
	homedir "github.com/mitchellh/go-homedir"
)

const alphanumericRegex = "^[a-zA-Z0-9_-]*$"
const dbFile = "~/.golog"

var dbPath, _ = homedir.Expand(dbFile)
var taskService tasksServices.TaskService
var transformer = Transformer{}
var commands = []cli.Command{
	{
		Name:         "start",
		Usage:        "Start tracking a given task",
		Action:       Start,
		BashComplete: AutocompleteTasks,
	},
	{
		Name:         "stop",
		Usage:        "Stop tracking a given task",
		Action:       Stop,
		BashComplete: AutocompleteTasks,
	},
	{
		Name:         "switch",
		Usage:        "Switch to a given task, stops all running tasks and starts a given task",
		ArgsUsage:    "[task_name]",
		Action:       Switch,
		BashComplete: AutocompleteTasks,
	},
	{
		Name:         "status",
		Usage:        "Give status of all tasks",
		Action:       Status,
		BashComplete: AutocompleteTasks,
	},
	{
		Name:   "clear",
		Usage:  "Clear all data",
		Action: Clear,
	},
	{
		Name:   "list",
		Usage:  "List all tasks",
		Action: List,
	},
	{
		Name:         "export",
		Usage:        "Export all tasks in a specific format",
		ArgsUsage:    fmt.Sprintf("[%s] [filePath]", strings.Join(file.GetFormatNames(), " | ")),
		Action:       Export,
		BashComplete: AutocompleteExport,
	},
}

// Start a given task
func Start(context *cli.Context) error {
	identifier := context.Args().First()
	if !IsValidIdentifier(identifier) {
		return invalidIdentifier(identifier)
	}

	err := taskService.StartTask(identifier)

	if err == nil {
		fmt.Println("Started tracking ", identifier)
	}
	return err
}

// Stop a given task
func Stop(context *cli.Context) error {
	identifier := context.Args().First()
	if !IsValidIdentifier(identifier) {
		return invalidIdentifier(identifier)
	}

	err := taskService.PauseTask(identifier)

	if err == nil {
		fmt.Println("Stopped tracking ", identifier)
	}
	return err
}

// Switch to a given task
func Switch(context *cli.Context) error {
	identifier := context.Args().First()
	if !IsValidIdentifier(identifier) {
		return invalidIdentifier(identifier)
	}

	err := taskService.SwitchTask(identifier)

	if err == nil {
		fmt.Println("Switched to task ", identifier)
	}
	return err
}

// Status display tasks being tracked
func Status(context *cli.Context) error {
	identifier := context.Args().First()
	if !IsValidIdentifier(identifier) {
		return invalidIdentifier(identifier)
	}

	task, err := taskService.GetTask(identifier)
	if err != nil {
		return err
	}
	transformer.LoadedTasks = tasksModel.Collection{*task}
	fmt.Println(transformer.Transform()[identifier])
	return nil
}

// List lists all tasks
func List(context *cli.Context) error {
	var err error
	transformer.LoadedTasks, err = taskService.GetTasks()
	if err != nil {
		return err
	}

	for _, task := range transformer.Transform() {
		fmt.Println(task)
	}
	return nil
}

// Export all tasks to a specific format
func Export(context *cli.Context) error {
	format, filePath := file.Format(strings.ToLower(context.Args().Get(0))), context.Args().Get(1)

	if format == "" {
		format = file.CSV
	}

	if filePath == "" {
		filePath = fmt.Sprintf("./golog%s", file.GetFormatFileExtension(format))
	}

	filePath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	err = taskService.Export(format, filePath)
	if err == nil {
		fmt.Println("Exported tasks to ", filePath)
	}
	return err
}

// Clear all data
func Clear(context *cli.Context) error {
	err := taskService.DeleteTasks()
	if err == nil {
		fmt.Println("All tasks deleted")
	}
	return err
}

// AutocompleteTasks loads tasks from repository and show them for completion
func AutocompleteTasks(context *cli.Context) {
	var err error
	transformer.LoadedTasks, err = taskService.GetTasks()
	// This will complete if no args are passed
	//   or there is problem with tasks repo
	if len(context.Args()) > 0 || err != nil {
		return
	}

	for _, task := range transformer.LoadedTasks {
		fmt.Println(task.Identifier)
	}
}

// AutocompleteExport shows the list of available export formats
func AutocompleteExport(context *cli.Context) {
	if len(context.Args()) > 0 {
		return
	}

	for _, exporter := range file.GetFormatNames() {
		fmt.Println(exporter)
	}
}

// IsValidIdentifier checks if the string passed is a valid task identifier
func IsValidIdentifier(identifier string) bool {
	re := regexp.MustCompile(alphanumericRegex)
	return len(identifier) > 0 && re.MatchString(identifier)
}

func checkInitialDbFile() {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		os.Create(dbPath)
	}
}

func runCliApp() (err error) {
	// @todo remove this from here, should be in file repo implementation
	checkInitialDbFile()

	taskService = tasksServices.New(csv.New(dbPath))

	app := cli.NewApp()
	app.Name = "Golog"
	app.Usage = "Easy CLI time tracker for your tasks"
	app.Version = "0.1"
	app.EnableBashCompletion = true
	app.Commands = commands

	err = app.Run(os.Args)

	return err
}

func main() {
	err := runCliApp()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func invalidIdentifier(identifier string) error {
	return fmt.Errorf("identifier %q is invalid", identifier)
}
