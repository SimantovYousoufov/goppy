package main

import (
	"github.com/atotto/clipboard"
	"fmt"
	"time"
	"strings"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

type Application struct {
	TerminalWidth int
	History       *History
	Storage       Storage
}

func NewApplication(s Storage, limit int) (*Application, error) {
	a := &Application{
		Storage: s,
		History: NewHistory(limit),
	}

	width, err := a.GetWindowWidth()

	if err != nil {
		return nil, err
	}

	a.TerminalWidth = width

	history, err := a.Storage.Read()

	sideLength := (a.TerminalWidth - len(AppName)) / 2
	fmt.Printf("%s%s%s\n", strings.Repeat("=", sideLength), "Goppy", strings.Repeat("=", sideLength))

	if err == nil {
		a.Hydrate(history)
	}

	return a, nil
}

func (a *Application) Watch() error {
	ClearScreen()

	var last string

	for {
		s, err := clipboard.ReadAll()

		if err != nil {
			return err
		}

		if s == last {
			time.Sleep(SleepTime)

			continue
		}

		a.History.Push(s)
		last = s

		a.Draw(a.History.First())

		time.Sleep(SleepTime)
	}

	return nil
}

type ClipboardItem struct {
	Contents  string
	Timestamp time.Time
	next      *ClipboardItem
}

type History struct {
	Limit int
	Size  int
	Head  *ClipboardItem
}

//
// @todo implement draw separately with a Screen interface
//
func (a *Application) Draw(i *ClipboardItem) {
	fmt.Printf("%s\n", i.Contents)

	fmt.Printf("%s\n", strings.Repeat("=", a.TerminalWidth))
}

func (a *Application) GetWindowWidth() (int, error) {
	w, _, err := terminal.GetSize(int(os.Stdin.Fd()))

	if err != nil {
		return 0, FailedToGetTerminalWidth
	}

	return MaxInt(w, HeaderLength), nil
}

func (a *Application) SaveHistory() error {
	return a.Storage.Store(a.History)
}

func (a *Application) Hydrate(items *HistoryItems) error {
	for _, item := range items.Items {
		err := a.History.PushClipboardItem(item)

		if err != nil {
			return err
		}

		a.Draw(item)
	}

	return nil
}

func NewHistory(limit int) *History {
	return &History{
		Size:  0,
		Limit: limit,
	}
}

func (h *History) First() *ClipboardItem {
	return h.Head
}

func (h *History) Push(s string) error {
	h.PushClipboardItem(&ClipboardItem{
		Contents:  s,
		Timestamp: time.Now(),
	})

	return nil
}

func (h *History) PushClipboardItem(i *ClipboardItem) error {
	if h.Size == h.Limit {
		h.Pop()
	}

	i.next = h.Head
	h.Head = i
	h.Size += 1

	return nil
}

func (h *History) Pop() (*ClipboardItem, error) {
	if h.Size == 0 {
		return nil, EmptyHistoryError
	}

	current := h.Head

	for current.next.next != nil {
		current = current.next
	}

	// Pop last item
	current.next = nil
	h.Size -= 1

	return current, nil
}

func (h *History) Iterate() <-chan *ClipboardItem {
	ch := make(chan *ClipboardItem)

	go func() {
		defer close(ch)

		current := h.Head

		for current != nil {
			ch <- current

			current = current.next
		}
	}()

	return ch
}
