package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/matrix-org/gomatrix"
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

func (s *FStore) encodeRoom(room *gomatrix.Room) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(room)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *FStore) decodeRoom(room []byte) (*gomatrix.Room, error) {
	var r *gomatrix.Room
	buf := bytes.NewBuffer(room)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&r)
	if err != nil {
		return nil, err
	}
	return r, nil
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
	return string(data), nil
}

// SaveFilterID exposed for gomatrix
func (s *FStore) SaveFilterID(userID, filterID string) {
	s.Set(fmt.Sprintf("filter_%s", userID), filterID)
}

// LoadFilterID exposed for gomatrix
func (s *FStore) LoadFilterID(userID string) string {
	filter, _ := s.Get(fmt.Sprintf("filter_%s", userID))
	return filter
}

func (s *FStore) SaveNextBatch(userID, nextBatchToken string) {
	s.Set(fmt.Sprintf("batch_%s", userID), nextBatchToken)
}

// LoadNextBatch exposed for gomatrix
func (s *FStore) LoadNextBatch(userID string) string {
	batch, _ := s.Get(fmt.Sprintf("batch_%s", userID))
	return batch
}

// SaveRoom exposed for gomatrix
func (s *FStore) SaveRoom(room *gomatrix.Room) {
	b, _ := s.encodeRoom(room)
	s.Set(fmt.Sprintf("room_%s", room.ID), string(b))
}

// LoadRoom exposed for gomatrix
func (s *FStore) LoadRoom(roomID string) *gomatrix.Room {
	b, _ := s.Get(fmt.Sprintf("room_%s", roomID))
	room, _ := s.decodeRoom([]byte(b))
	return room
}
