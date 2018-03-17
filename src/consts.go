package main

import (
	"errors"
	"time"
)

var (
	FailedToGetTerminalWidth = errors.New("failed to get terminal width")
	EmptyHistoryError        = errors.New("history is empty")
	FailedToWriteToFile      = errors.New("failed to write history to file")
	FailedToReadFile         = errors.New("failed to read history from file")
)

const (
	HeaderLength    = 100
	AppName         = "Goppy"
	SleepTime       = 1 * time.Second
	HistoryFilename = "goppy_history.json"
)
