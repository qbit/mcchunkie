package main

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/matrix-org/gomatrix"
	bolt "go.etcd.io/bbolt"
)

// MCStore implements a gomatrix.Storer and exposes a bbolt db to be used for
// application storage (account info, config info etc).
type MCStore struct {
	db *bolt.DB
}

// NewStore creates a new MCStore instance populated with the proper buckets.
func NewStore(path string) (*MCStore, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	s := &MCStore{db: db}

	err = s.db.Update(func(tx *bolt.Tx) error {
		buckets := []string{"filter", "batch", "room", "account", "config", "errata"}
		for _, b := range buckets {
			if _, err := tx.CreateBucketIfNotExists([]byte(b)); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		s.db.Close()
		return nil, err
	}
	return s, nil
}

func (s *MCStore) set(bucket, key, value string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		err = bkt.Put([]byte(key), []byte(value))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (s *MCStore) get(bucket, key string) (string, error) {
	var result string
	return result, s.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bucket))
		data := bkt.Get([]byte(key))

		if data == nil {
			return nil
		}

		result = string(data)

		return nil
	})
}

func (s *MCStore) getBytes(bucket, key string) ([]byte, error) {
	var result []byte
	return result, s.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bucket))
		data := bkt.Get([]byte(key))

		if data == nil {
			return nil
		}

		result = data
		return nil
	})
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
	_ = s.set("filter", userID, filterID)

}

// LoadFilterID exposed for gomatrix
func (s *MCStore) LoadFilterID(userID string) string {
	filter, _ := s.get("filter", userID)
	return string(filter)
}

// SaveNextBatch exposed for gomatrix
func (s *MCStore) SaveNextBatch(userID, nextBatchToken string) {
	_ = s.set("batch", userID, nextBatchToken)
}

// LoadNextBatch exposed for gomatrix
func (s *MCStore) LoadNextBatch(userID string) string {
	batch, _ := s.get("batch", userID)
	return string(batch)
}

// SaveRoom exposed for gomatrix
func (s *MCStore) SaveRoom(room *gomatrix.Room) {
	b, _ := s.encodeRoom(room)
	_ = s.set("room", room.ID, string(b))
}

// LoadRoom exposed for gomatrix
func (s *MCStore) LoadRoom(roomID string) *gomatrix.Room {
	b, _ := s.getBytes("room", roomID)
	room, _ := s.decodeRoom(b)
	return room
}
