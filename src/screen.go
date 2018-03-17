package main

import (
	"os/exec"
	"os"
	"runtime"
)

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
