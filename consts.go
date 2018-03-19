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
	FailedToEncrypt          = errors.New("failed to encrypt")
	FailedToDecrypt          = errors.New("failed to decrypt")
)

const (
	HeaderLength       = 100
	AppName            = "Goppy"
	SleepTime          = 1 * time.Second
	HistoryFilename    = "goppy_history.dat"
	GoppyConfigFolder  = "/usr/local/etc/goppy/"
	DefaultHistoryFile = GoppyConfigFolder + HistoryFilename
	SaltBytes          = 32
	Pbkdf2Iters        = 4096
	KeySize            = 32
	TruncatedSuffix    = "..."
)
