package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

// FStore is the path to a directory which will contain our data.
type FStore string

// NewStore creates a new instance of FStore
func NewStore(s string) (*FStore, error) {
	fi, err := os.Lstat(s)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("not a directory")
	}
	fstore := FStore(s)
	return &fstore, nil
}

// Set dumps value into a file named key
func (s FStore) Set(key string, value string) {
	err := ioutil.WriteFile(path.Join(string(s), key), []byte(value), 0600)
	if err != nil {
		log.Println(err)
	}
}

// Get pulls value from a file named key
func (s FStore) Get(key string) (string, error) {
	data, err := ioutil.ReadFile(path.Join(string(s), key))
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(data)), nil
}
