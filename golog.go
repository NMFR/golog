package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mlimaloureiro/golog/models"
	"github.com/mlimaloureiro/golog/repositories"

	"github.com/codegangsta/cli"
	homedir "github.com/mitchellh/go-homedir"
)

const (
	exportCsv  = "csv"
	exportICal = "ical"
	exportIcs  = "ics"
)

const alphanumericRegex = "^[a-zA-Z0-9_-]*$"
const dbFile = "~/.golog"

var dbPath, _ = homedir.Expand(dbFile)
var repositoryFile *os.File
var taskRepository models.TaskRepositoryInterface
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
		ArgsUsage:    fmt.Sprintf("[%s] [filePath]", strings.Join(getExportFormatNames(), " | ")),
		Action:       Export,
		BashComplete: AutocompleteExport,
	},
}

func getExportFormatNames() []string {
	return []string{exportCsv, exportICal, exportIcs}
}

// Start a given task
func Start(context *cli.Context) error {
	identifier := context.Args().First()
	if !IsValidIdentifier(identifier) {
		return invalidIdentifier(identifier)
	}

	err := taskRepository.StartTask(identifier)

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

	err := taskRepository.PauseTask(identifier)

	if err == nil {
		fmt.Println("Stopped tracking ", identifier)
	}
	return err
}

// Status display tasks being tracked
func Status(context *cli.Context) error {
	identifier := context.Args().First()
	if !IsValidIdentifier(identifier) {
		return invalidIdentifier(identifier)
	}

	task, err := taskRepository.GetTask(identifier)
	if err != nil {
		return err
	}
	transformer.LoadedTasks = models.Tasks{*task}
	fmt.Println(transformer.Transform()[identifier])
	return nil
}

// List lists all tasks
func List(context *cli.Context) error {
	var err error
	transformer.LoadedTasks, err = taskRepository.GetTasks()
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
	formatDefaultExtensionMap := map[string]string{
		"csv":  ".csv",
		"ical": ".ics",
		"ics":  ".ics",
	}
	format, outputFilePath := context.Args().Get(0), context.Args().Get(1)

	if format == "" {
		format = exportCsv
	} else {
		format = strings.ToLower(format)
	}

	if outputFilePath == "" {
		outputFilePath = fmt.Sprintf("./golog%s", formatDefaultExtensionMap[format])
	}

	outputFilePath, err := filepath.Abs(outputFilePath)
	if err != nil {
		return err
	}

	tasks, err := taskRepository.GetTasks()
	if err != nil {
		return err
	}

	var file *os.File
	file, err = os.OpenFile(outputFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	var exporter models.TaskRepositoryInterface
	switch format {
	case exportCsv:
		exporter = repositories.NewCsvTaskRepository(file)
	case exportICal:
		fallthrough
	case exportIcs:
		exporter = repositories.NewICalTaskRepository(file)
	}

	if exporter == nil {
		return fmt.Errorf("invalid format \"%s\"", format)
	}

	err = exporter.SetTasks(tasks)
	if err == nil {
		fmt.Println("Exported tasks to ", outputFilePath)
	}

	return err
}

// Clear all data
func Clear(context *cli.Context) error {
	err := repositoryFile.Truncate(0)
	if err == nil {
		fmt.Println("All tasks deleted")
	}
	return err
}

// AutocompleteTasks loads tasks from repository and show them for completion
func AutocompleteTasks(context *cli.Context) {
	var err error
	transformer.LoadedTasks, err = taskRepository.GetTasks()
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

	for _, exporter := range getExportFormatNames() {
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

	repositoryFile, err := os.OpenFile(dbPath, os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := repositoryFile.Close(); err == nil {
			err = closeErr
		}
	}()

	taskRepository = repositories.NewCsvTaskRepository(repositoryFile)

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
