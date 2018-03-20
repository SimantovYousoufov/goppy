package main

import (
	"flag"
	"os"
	"path/filepath"
	"fmt"
	"os/signal"
	"syscall"
	"golang.org/x/crypto/ssh/terminal"
	"strings"
	"os/exec"
	"runtime"
)

func main() {
	help := flag.Bool("h", false, "Get usage information.")
	nullStore := flag.Bool("s", false, "Use the NullStore to prevent clipboard history from being written to disk.")
	encryptedStore := flag.Bool("e", false, "Use the encrypted file storage format.")
	limit := flag.Int("n", 50, "Number of items to keep in history.")
	file := flag.String("f", DefaultHistoryFile, "Path to history file.")

	flag.Parse()

	err := checkOrCreateGoppyConfigFolder()

	if err != nil {
		fmt.Println("Could not create Goppy config directory.")
		panic(err)
	}

	if *help {
		flag.PrintDefaults()

		os.Exit(0)
	}

	store, err := ChooseStore(*file, *nullStore, *encryptedStore)

	if err != nil {
		fmt.Println("Could not init storage engine.")
		panic(err)
	}

	screen, err := NewTerminalScreen()

	if err != nil {
		fmt.Println("Could not init terminal screen.")
		panic(err)
	}

	app, err := NewApplication(store, screen, *limit)

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
