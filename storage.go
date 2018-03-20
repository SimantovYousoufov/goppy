package main

import (
	"os"
	"fmt"
	"io/ioutil"
)

//
// Decide which storage to use based on passed flags, prioritizing no storage
//
func ChooseStore(path string, useNullStore, useEncryptedStore bool) (Storage, error) {
	if useNullStore {
		return &NullStore{}, nil
	}

	if useEncryptedStore {
		return NewEncryptedStore(path, CollectPassword())
	}

	return NewFileStore(path)
}

//
// A Storage type provides methods for storing and reading history (such as from a file)
//
type Storage interface {
	Store(*History) error
	Read() (*HistoryItems, error)
	Clear() error
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

	b, err := h.toJson()

	if err != nil {
		return FailedToWriteToFile
	}

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

	b, err := ioutil.ReadAll(fs.file)

	if err != nil {
		return nil, err
	}

	// Empty history, nothing to read
	if len(b) == 0 {
		return h, nil
	}

	return h.fromJson(b)
}

func (fs *FileStore) Clear() error {
	return truncateFile(fs.file)
}

//
// Implement storage which does not store
//
type NullStore struct{}

func (n *NullStore) Clear() error {
	// Pass

	return nil
}

func (n *NullStore) Store(*History) error {
	// Pass

	return nil
}

func (n *NullStore) Read() (*HistoryItems, error) {
	// Pass

	return &HistoryItems{
		Items: make([]*ClipboardItem, 0),
	}, nil
}

//
// Implement an encrypted file store
//
type EncryptedStore struct {
	path string
	key  []byte
	file *os.File
}

func NewEncryptedStore(path string, key string) (*EncryptedStore, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, FailedToReadFile
	}

	return &EncryptedStore{
		path: path,
		file: f,
		key:  []byte(key),
	}, nil
}

func (es *EncryptedStore) Store(h *History) error {
	fmt.Printf("\nWriting encrypted history to %s\n", es.path)

	b, err := h.toJson()

	if err != nil {
		return FailedToWriteToFile
	}

	fmt.Println("Serialized to JSON")

	b, err = encrypt(es.key, b)

	if err != nil {
		return err
	}

	fmt.Println("Encrypted data")

	_, err = es.file.WriteAt(b, 0)

	if err != nil {
		return FailedToWriteToFile
	}

	fmt.Println("Wrote data to file")

	return nil
}

func (es *EncryptedStore) Read() (*HistoryItems, error) {
	h := &HistoryItems{
		Items: make([]*ClipboardItem, 0),
	}

	b, err := ioutil.ReadAll(es.file)

	if err != nil {
		return nil, err
	}

	// Empty history, nothing to read
	if len(b) == 0 {
		return h, nil
	}

	b, err = decrypt(es.key, b)

	return h.fromJson(b)
}

func (es *EncryptedStore) Clear() error {
	return truncateFile(es.file)
}

func checkOrCreateGoppyConfigFolder() error {
	_, err := os.Stat(GoppyConfigFolder)

	if os.IsNotExist(err) {
		return os.Mkdir(GoppyConfigFolder, 0755)
	}

	return err
}

func truncateFile(file *os.File) error {
	return file.Truncate(0)
}