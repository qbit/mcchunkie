package main

import (
	"bytes"
	"encoding/gob"

	"github.com/matrix-org/gomatrix"
	"github.com/peterbourgon/diskv"
)

// MCStore implements a gomatrix.Storer and exposes a diskv db to be used for
// application storage (account info, config info etc).
type MCStore struct {
	db *diskv.Diskv
}

// NewStore creates a new MCStore instance populated with the proper buckets.
func NewStore(path string) (*MCStore, error) {
	flatTransform := func(s string) []string { return []string{} }
	db := diskv.New(diskv.Options{
		BasePath:     "db",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})

	s := &MCStore{db: db}

	return s, nil
}

func (s *MCStore) set(key string, value string) {
	v := []byte(value)
	s.db.Write(key, v)
}

func (s *MCStore) get(key string) (string, error) {
	b, err := s.db.Read(key)
	return string(b), err
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

// SaveFilterID exposed for gomatrix
func (s *MCStore) SaveFilterID(userID, filterID string) {
	s.set(userID, filterID)

}

// LoadFilterID exposed for gomatrix
func (s *MCStore) LoadFilterID(userID string) string {
	filter, _ := s.get(userID)
	return string(filter)
}

// SaveNextBatch exposed for gomatrix
func (s *MCStore) SaveNextBatch(userID, nextBatchToken string) {
	s.set(userID, nextBatchToken)
}

// LoadNextBatch exposed for gomatrix
func (s *MCStore) LoadNextBatch(userID string) string {
	batch, _ := s.get(userID)
	return string(batch)
}

// SaveRoom exposed for gomatrix
func (s *MCStore) SaveRoom(room *gomatrix.Room) {
	b, _ := s.encodeRoom(room)
	s.set(room.ID, string(b))
}

// LoadRoom exposed for gomatrix
func (s *MCStore) LoadRoom(roomID string) *gomatrix.Room {
	b, _ := s.get(roomID)
	room, _ := s.decodeRoom([]byte(b))
	return room
}
