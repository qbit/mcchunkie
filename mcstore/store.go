package mcstore

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/matrix-org/gomatrix"
)

// MCStore is the path to a directory which will contain our data.
type MCStore string

// NewStore creates a new instance of FStore
func NewStore(s string) (*MCStore, error) {
	fi, err := os.Lstat(s)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("not a directory")
	}
	fstore := MCStore(s)
	return &fstore, nil
}

func (s *MCStore) encodeRoom(room *gomatrix.Room) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(room)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *MCStore) decodeRoom(room []byte) (*gomatrix.Room, error) {
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
func (s MCStore) Set(key string, value string) {
	err := os.WriteFile(path.Join(string(s), key), []byte(value), 0600)
	if err != nil {
		log.Println(err)
	}
}

// Get pulls value from a file named key
func (s MCStore) Get(key string) (string, error) {
	data, err := os.ReadFile(path.Join(string(s), key))
	if err != nil {
		return "", fmt.Errorf("no entry for %q: %q", key, err)
	}
	return strings.TrimSpace(string(data)), nil
}

// SaveFilterID exposed for gomatrix
func (s *MCStore) SaveFilterID(userID, filterID string) {
	s.Set(fmt.Sprintf("filter_%s", userID), filterID)
}

// LoadFilterID exposed for gomatrix
func (s *MCStore) LoadFilterID(userID string) string {
	filter, _ := s.Get(fmt.Sprintf("filter_%s", userID))
	return filter
}

func (s *MCStore) SaveNextBatch(userID, nextBatchToken string) {
	s.Set(fmt.Sprintf("batch_%s", userID), nextBatchToken)
}

// LoadNextBatch exposed for gomatrix
func (s *MCStore) LoadNextBatch(userID string) string {
	batch, _ := s.Get(fmt.Sprintf("batch_%s", userID))
	return batch
}

// SaveRoom exposed for gomatrix
func (s *MCStore) SaveRoom(room *gomatrix.Room) {
	b, _ := s.encodeRoom(room)
	s.Set(fmt.Sprintf("room_%s", room.ID), string(b))
}

// LoadRoom exposed for gomatrix
func (s *MCStore) LoadRoom(roomID string) *gomatrix.Room {
	b, _ := s.Get(fmt.Sprintf("room_%s", roomID))
	room, _ := s.decodeRoom([]byte(b))
	return room
}
