package main

import (
	"flag"
	"os"
	"path/filepath"
	"fmt"
	"os/signal"
	"syscall"
)

//
// @todos
// - Implement a Screen interface for drawing
// - Implement a null store to not store history
// - Implement way to clear history
// - Implement file storage with encryption
//

func main() {
	limit := flag.Int("n", 50, "Number of items to keep in history.")
	file := flag.String("f", fmt.Sprintf("%s/%s", GetCwd(), HistoryFilename), "Path to history file, defaults to $PWD/goppy_history.json")

	flag.Parse()

	fs, err := NewFileStore(*file)

	if err != nil {
		fmt.Println("Encountered error while opening file.")
		panic(err)
	}

	app, err := NewApplication(fs, *limit)

	if err != nil {
		panic(err)
	}

	CatchSignals(app)

	err = app.Watch()

	if err != nil {
		panic(err)
	}
}

//
// Find this applications current working directory
//
func GetCwd() string {
	ex, err := os.Executable()

	if err != nil {
		panic(err)
	}

	return filepath.Dir(ex)
}

//
// Catch relevant signals and force a history save before we exit
//
func CatchSignals(a *Application) {
	c := make(chan os.Signal, 3)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		<-c

		err := a.SaveHistory()

		if err != nil {
			fmt.Printf("Encountered error while writing to file: %v.\n", err)
			panic(err)
		}

		fmt.Println("Successfully wrote history to file.")

		os.Exit(0)
	}()
}

//
// Utility function for getting the max of two ints
//
func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}
