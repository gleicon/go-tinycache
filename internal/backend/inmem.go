package backend

import (
	"errors"
	"fmt"

	"github.com/gleicon/go-tinycache/internal/datautils"
)

type LRUMemoryBackend struct {
	size int
	data *datautils.LRUCache
}

// NewInmemBackend is the exported constructor for the in-memory backend
func NewLRUMemoryBackend(size int) (*LRUMemoryBackend, error) {
	dd := datautils.NewLRU(size)
	b := LRUMemoryBackend{size: size, data: dd}
	return &b, nil
}

func (be LRUMemoryBackend) Set(key []byte, value []byte) error {
	return be.Put(key, value, false, true)
}

// store data only if the server doesnt holds it yet
func (be LRUMemoryBackend) Add(key []byte, value []byte) error {
	return be.Put(key, value, false, false)
}

// store data only if the server already holds this key
func (be LRUMemoryBackend) Replace(key []byte, value []byte) error {
	return be.Put(key, value, true, false)
}

/*
Incr data, yields error if the represented value doesnt maps to int.
Starts from 0, no negative values
*/
func (be LRUMemoryBackend) Incr(key []byte, value uint) (int, error) {
	return be.Increment(key, int(value), false)
}

/*
Decr data, yields error if the represented value doesnt maps to int.
Stops at 0, no negative values
*/
func (be LRUMemoryBackend) Decr(key []byte, value uint) (int, error) {
	return be.Increment(key, int(value)*-1, false)
}

// Generic get and set for incr/decr tx
func (be LRUMemoryBackend) Increment(key []byte, value int, create_if_not_exists bool) (int, error) {
	return 0, nil
}

func (be LRUMemoryBackend) Put(key []byte, value []byte, replace bool, passthru bool) error {
	be.data.Set(string(key), string(value)) //, time.Now())
	return nil
}

func (be LRUMemoryBackend) Get(key []byte) ([]byte, error) {
	r, ok := be.data.Get(string(key))
	if !ok {
		return nil, errors.New("Error getting value from inmem")
	}
	return []byte(r), nil
}

// returns deleted, error
func (be LRUMemoryBackend) Delete(key []byte, only_if_exists bool) (bool, error) {
	if only_if_exists == true {
		x, err := be.Get(key)
		if err != nil {
			return false, err
		}
		if x == nil {
			return false, nil
		}
	}
	be.data.Delete(string(key))
	return true, nil
}

func (be LRUMemoryBackend) Flush() error {
	return nil
}

func (be LRUMemoryBackend) BucketStats() error {
	return nil
}

func (be LRUMemoryBackend) GetDbPath() string {
	rm := fmt.Sprintf("In memory database: %d bytes", be.data.Size())
	return rm
}

func (be LRUMemoryBackend) SwitchBucket(bucket string) {}
func (be LRUMemoryBackend) Range([]byte, int, []byte, bool) (map[string][]byte, error) {
	return nil, nil
}
func (be LRUMemoryBackend) Close()        {}
func (be LRUMemoryBackend) Stats() string { return "" }
