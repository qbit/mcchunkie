package chats

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/matrix-org/gomatrix"
)

// MStore is the path to a directory which will contain our data.
type MStore struct {
	ParentStore interface{}
}

// NewStore creates a new instance of MStore
func NewStore() *MStore {
	fstore := MStore{}
	return &fstore
}

func (s *MStore) encodeRoom(room *gomatrix.Room) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(room)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *MStore) decodeRoom(room []byte) (*gomatrix.Room, error) {
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
func (s MStore) Set(key string, value string) {
	if ds, ok := s.ParentStore.(MStore); ok {
		ds.Set(key, value)
	}
}

// Get pulls value from a file named key
func (s MStore) Get(key string) (string, error) {
	if ds, ok := s.ParentStore.(MStore); ok {
		return ds.Get(key)
	}
	return "", errors.New("can't initialize data store")
}

// SaveFilterID exposed for gomatrix
func (s *MStore) SaveFilterID(userID, filterID string) {
	s.Set(fmt.Sprintf("filter_%s", userID), filterID)
}

// LoadFilterID exposed for gomatrix
func (s *MStore) LoadFilterID(userID string) string {
	filter, _ := s.Get(fmt.Sprintf("filter_%s", userID))
	return filter
}

// SaveNextBatch exposed for gomatrix
func (s *MStore) SaveNextBatch(userID, nextBatchToken string) {
	s.Set(fmt.Sprintf("batch_%s", userID), nextBatchToken)
}

// LoadNextBatch exposed for gomatrix
func (s *MStore) LoadNextBatch(userID string) string {
	batch, _ := s.Get(fmt.Sprintf("batch_%s", userID))
	return batch
}

// SaveRoom exposed for gomatrix
func (s *MStore) SaveRoom(room *gomatrix.Room) {
	b, _ := s.encodeRoom(room)
	s.Set(fmt.Sprintf("room_%s", room.ID), string(b))
}

// LoadRoom exposed for gomatrix
func (s *MStore) LoadRoom(roomID string) *gomatrix.Room {
	b, _ := s.Get(fmt.Sprintf("room_%s", roomID))
	room, _ := s.decodeRoom([]byte(b))
	return room
}
