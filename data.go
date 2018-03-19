package main

import (
	"github.com/atotto/clipboard"
	"time"
	"encoding/json"
)

type Application struct {
	Screen  Screen
	History *History
	Storage Storage
}

func NewApplication(storage Storage, screen Screen, limit int) (*Application, error) {
	a := &Application{
		Storage: storage,
		History: NewHistory(limit),
		Screen:  screen,
	}

	history, err := a.Storage.Read()

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

		a.Screen.Draw(a.History.First())

		time.Sleep(SleepTime)
	}

	return nil
}

func (a *Application) SaveHistory() error {
	return a.Storage.Store(a.History)
}

func (a *Application) ClearHistory() {
	a.History.Clear()

	a.SaveHistory()
}

func (a *Application) Hydrate(items *HistoryItems) error {
	for _, item := range items.Items {
		err := a.History.PushClipboardItem(item)

		if err != nil {
			return err
		}

		a.Screen.Draw(item)
	}

	return nil
}

type ClipboardItem struct {
	Contents  string
	Timestamp time.Time
	next      *ClipboardItem
}

type HistoryItems struct {
	Items []*ClipboardItem
}

type History struct {
	Limit int
	Size  int
	Head  *ClipboardItem
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
		h.DropLast()
	}

	i.next = h.Head
	h.Head = i
	h.Size += 1

	return nil
}

func (h *History) DropLast() (*ClipboardItem, error) {
	if h.Size == 0 {
		return nil, EmptyHistoryError
	}

	current := h.Head

	for current.next.next != nil {
		current = current.next
	}

	last := current.next

	// DropLast last item
	current.next = nil
	h.Size -= 1

	return last, nil
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

func (h *History) toJson() ([]byte, error) {
	data := &HistoryItems{
		Items: make([]*ClipboardItem, 0),
	}

	ch := h.Iterate()

	for item := range ch {
		data.Items = append(data.Items, item)
	}

	b, err := json.Marshal(data)

	return b, err
}

func (h *History) Clear() {
	h.Head = nil
	h.Size = 0
}

func (h *HistoryItems) fromJson(data []byte) (*HistoryItems, error) {
	err := json.Unmarshal(data, h)

	return h, err
}
