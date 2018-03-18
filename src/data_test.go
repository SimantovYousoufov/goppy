package main

import (
	"testing"
	"time"
)

func TestApplicationHydratesFromHistory(t *testing.T) {
	h := &HistoryItems{
		Items: make([]*ClipboardItem, 0),
	}

	h.Items = append(h.Items, &ClipboardItem{
		Contents:  "first item",
		Timestamp: time.Now().Add(5 * time.Second),
	})

	h.Items = append(h.Items, &ClipboardItem{
		Contents:  "second item",
		Timestamp: time.Now().Add(10 * time.Second),
	})

	h.Items = append(h.Items, &ClipboardItem{
		Contents:  "third item",
		Timestamp: time.Now().Add(15 * time.Second),
	})

	f := &FakeStore{
		H: h,
	}

	a := &Application{
		Storage: f,
		Screen: &NoScreen{},
		History: NewHistory(5),
	}

	err := a.Hydrate(h)

	assertNilError(t, err)

	if a.History.Size != 3 {
		t.Fatal("History count mismatch")
	}

	ch := a.History.Iterate()

	i := 2
	for item := range ch {
		if item.Contents != h.Items[i].Contents {
			t.Fatal("History items were not hydrated properly")
		}

		i--
	}
}

func TestHistoryCanPushNewStringIntoStackAndIncrementSize(t *testing.T) {
	h := NewHistory(5)

	if h.Size != 0 {
		t.Fatal("History size mismatch")
	}

	h.Push("first")

	if h.Size != 1 {
		t.Fatal("History size mismatch")
	}

	if h.First().Contents != "first" {
		t.Fatal("History did not set first item")
	}

	h.Push("second")

	if h.Size != 2 {
		t.Fatal("History size mismatch")
	}

	if h.First().Contents != "second" {
		t.Fatal("History did not set second item")
	}

	h.Push("third")

	if h.Size != 3 {
		t.Fatal("History size mismatch")
	}

	if h.First().Contents != "third" {
		t.Fatal("History did not set third item")
	}
}

func TestHistoryCanDropOldestItemFromStackAndDecrementSize(t *testing.T) {
	h := NewHistory(5)
	h.Push("first")
	h.Push("second")
	h.Push("third")

	if h.Size != 3 {
		t.Fatal("History size mismatch")
	}

	oldest, err := h.DropLast()

	assertNilError(t, err)

	if oldest.Contents != "first" {
		t.Fatal("History failed to drop oldest item from stack")
	}
}

func TestHistoryDropsOldestIfAtSizeLimit(t *testing.T) {
	h := NewHistory(5)
	h.Push("first")
	h.Push("second")
	h.Push("third")
	h.Push("fourth")
	h.Push("fifth")

	if h.Size != 5 {
		t.Fatal("History size mismatch")
	}

	h.Push("sixth")

	if h.Size != 5 {
		t.Fatal("History size mismatch")
	}

	oldest, err := h.DropLast()

	assertNilError(t, err)

	if oldest.Contents != "second" {
		t.Fatal("History failed to drop items in expected order")
	}
}

func TestHistoryDoesNotTryToDropOldestWhenSize0(t *testing.T) {
	h := NewHistory(5)

	oldest, err := h.DropLast()

	if oldest != nil {
		t.Fatal("History item mismatch")
	}

	if err == nil || err != EmptyHistoryError {
		t.Fatal("History did not return expected error")
	}
}

func TestHistoryCanSerializeToAndFromJson(t *testing.T) {
	h := NewHistory(5)
	h.Push("first")
	h.Push("second")
	h.Push("third")

	data, err := h.toJson()

	assertNilError(t, err)

	newH := &HistoryItems{}

	newH.fromJson(data)

	if len(newH.Items) != 3 {
		t.Fatal("History size mismatch")
	}

	if newH.Items[0].Contents != "third" {
		t.Fatal("History order mismatch")
	}

	if newH.Items[1].Contents != "second" {
		t.Fatal("History order mismatch")
	}

	if newH.Items[2].Contents != "first" {
		t.Fatal("History order mismatch")
	}
}

type FakeStore struct {
	H *HistoryItems
}

func (f *FakeStore) Store(*History) error {
	panic("implement me")
}

func (f *FakeStore) Read() (*HistoryItems, error) {
	return f.H, nil
}

func (f *FakeStore) Clear() error {
	panic("implement me")
}
