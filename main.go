package main

import (
	"flag"
	"os"
	"fmt"
	"os/signal"
	"syscall"
	"golang.org/x/crypto/ssh/terminal"
	"strings"
	"os/exec"
	"runtime"
)

func main() {
	err := checkOrCreateGoppyConfigFolder()
	AssertNilError(err)

	// Watch command setup
	watchCommand := flag.NewFlagSet("watch", flag.ExitOnError)
	nullStore := watchCommand.Bool("s", false, "Use the NullStore to prevent clipboard history from being written to disk.")
	encryptedStore := watchCommand.Bool("e", false, "Use the encrypted file storage format.")
	limit := watchCommand.Int("n", 50, "Number of items to keep in history.")
	file := watchCommand.String("f", DefaultHistoryFile, "Path to history file.")

	// Clear command setup
	clearCommand := flag.NewFlagSet("clear", flag.ExitOnError)
	clearFile := clearCommand.String("f", DefaultHistoryFile, "Path to history file.")

	// Help command setup
	helpCommand := flag.NewFlagSet("help", flag.ExitOnError)

	commands := map[string]*flag.FlagSet{
		"watch": watchCommand,
		"clear": clearCommand,
		"help":  helpCommand,
	}

	if len(os.Args) < 2 {
		fmt.Println("A subcommand is required.")
		fmt.Println("Usage:")
		fmt.Println("goppy command [arg1] [arg2] ...")

		printDefaults(commands)
		os.Exit(1)
	}

	subCommand := os.Args[1]

	switch subCommand {
	case "watch":
		watchCommand.Parse(os.Args[2:])

		store, err := ChooseStore(*file, *nullStore, *encryptedStore)

		AssertNilError(err)

		screen, err := NewTerminalScreen()

		AssertNilError(err)

		app, err := NewApplication(store, screen, *limit)

		CatchSignals(app)

		err = app.Watch()

		AssertNilError(err)
	case "clear":
		store, err := NewFileStore(*clearFile)

		AssertNilError(err)

		store.Clear()
	case "help":
		printDefaults(commands)
	default:
		fmt.Printf("%s not supported.\n", subCommand)
		os.Exit(1)
	}
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

//
// Provide a password input with no echo
//
func CollectPassword() string {
	fmt.Print("Enter password: ")
	b, err := terminal.ReadPassword(int(syscall.Stdin))

	if err != nil {
		panic("Failed to read password.")
	}

	fmt.Println("")
	return strings.TrimSpace(string(b))
}

var clear map[string]func()

func init() {
	clear = make(map[string]func())

	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func ClearScreen() {
	value, ok := clear[runtime.GOOS]

	if ! ok {
		panic("Platform is unsupported.")
	}

	value()
}

// AssertNilError is intended to be called with which error return to simplify error handling
// Usage:
// foo, err := GetFoo()
// AssertNilError(err)
// DoSomethingBecauseNoError()
func AssertNilError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\nEncountered an error: %s\n", err.Error())

	os.Exit(1)
}

func printDefaults(commands map[string]*flag.FlagSet) {
	for cmd, set := range commands {
		fmt.Printf("Usage for `%s`:\n", cmd)
		set.PrintDefaults()
		fmt.Println()
	}
}
