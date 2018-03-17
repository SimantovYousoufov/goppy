package main

import (
	"os"
	"github.com/gin-gonic/gin/json"
	"fmt"
)

//
// A Storage type provides methods for storing and reading history (such as from a file)
//
type Storage interface {
	Store(*History) error

	Read() (*HistoryItems, error)
}

type HistoryItems struct {
	Items []*ClipboardItem
}

//
// Implement Storage as a JSON file store
//
type FileStore struct {
	path string
	file *os.File
}

func NewFileStore(path string) (*FileStore, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, FailedToReadFile
	}

	return &FileStore{
		path: path,
		file: f,
	}, nil
}

func (fs *FileStore) Store(h *History) error {
	fmt.Printf("\nWriting history to %s\n", fs.path)

	data := &HistoryItems{
		Items: make([]*ClipboardItem, 0),
	}

	ch := h.Iterate()

	for item := range ch {
		data.Items = append(data.Items, item)
	}

	b, err := json.Marshal(data)

	_, err = fs.file.WriteAt(b, 0)

	if err != nil {
		return FailedToWriteToFile
	}

	return nil
}

func (fs *FileStore) Read() (*HistoryItems, error) {
	h := &HistoryItems{
		Items: make([]*ClipboardItem, 0),
	}
	d := json.NewDecoder(fs.file)

	err := d.Decode(h)

	if err != nil {
		return nil, FailedToReadFile
	}

	return h, nil
}
