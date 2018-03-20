package main

import (
	"os"
	"golang.org/x/crypto/ssh/terminal"
	"fmt"
	"strings"
	"unicode/utf8"
)

type Screen interface {
	Draw(i *ClipboardItem) error
}

type TerminalScreen struct {
	TerminalWidth int
}

func NewTerminalScreen() (*TerminalScreen, error) {
	s := &TerminalScreen{}

	s.TerminalWidth = s.getWindowWidth()

	s.init()

	return s, nil
}

func (s *TerminalScreen) init() {
	sideLength := (s.TerminalWidth - len(AppName)) / 2
	fmt.Printf("%s%s%s\n", strings.Repeat("=", sideLength), "Goppy", strings.Repeat("=", sideLength))
}

func (s *TerminalScreen) Draw(i *ClipboardItem) error {
	if ! utf8.Valid(i.Contents) {
		fmt.Printf("<Cannot display binary data>\n")
	} else {
		fmt.Printf("%s\n", i.Contents)
	}

	fmt.Printf("%s\n", strings.Repeat("=", s.TerminalWidth))

	return nil
}

func (s *TerminalScreen) getWindowWidth() int {
	w, _, err := terminal.GetSize(int(os.Stdin.Fd()))

	if err != nil {
		w = 0
	}

	return MaxInt(w, HeaderLength)
}

type NoScreen struct {}

func (NoScreen) Draw(i *ClipboardItem) error {
	// Pass

	return nil
}

func truncateString(s string, maxLen, maxLines int) string {
	originalLength := len(s)
	splitLines := strings.Split(s, "\n")[:maxLines]
	s = strings.Join(splitLines, "\n")
	maxLen = maxAcceptableLength(s, maxLen)

	// Nothing was truncated
	if len(s) == originalLength {
		return s
	}

	return s[:maxLen] + TruncatedSuffix
}

func maxAcceptableLength(s string, maxLen int) int {
	if len(s) > maxLen {
		return maxLen
	}

	return len(s)
}