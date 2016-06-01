package main

import (
	"io/ioutil"
	// "fmt"
	// "os"
)

type Store interface {
	Write([]byte) error
	Read() ([]byte, error)
}

type FileStore struct {
	filename string
}

func (store *FileStore) Write(contents []byte) error {
	return ioutil.WriteFile(store.filename, contents, 0660)
}

func (store *FileStore) Read() ([]byte, error) {
	return ioutil.ReadFile(store.filename)
}
