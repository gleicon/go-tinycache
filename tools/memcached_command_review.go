package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	server := "127.0.0.1:11211"
	mc := memcache.New(server)

	tests := []struct {
		name string
		test func(*memcache.Client) error
	}{
		{"set/get", performSetGet},
		{"add", performAdd},
		{"replace", performReplace},
		{"delete", performDelete},
	}

	var wg sync.WaitGroup
	for _, t := range tests {
		t := t
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := t.test(mc); err != nil {
				log.Printf("%s: FAIL: %v", t.name, err)
			} else {
				log.Printf("%s: PASS", t.name)
			}
		}()
	}
	wg.Wait()
}

func performSetGet(mc *memcache.Client) error {
	key := "foo"
	val := "bar"
	if err := mc.Set(&memcache.Item{Key: key, Value: []byte(val)}); err != nil {
		return fmt.Errorf("set: %w", err)
	}
	item, err := mc.Get(key)
	if err != nil {
		return fmt.Errorf("get: %w", err)
	}
	if string(item.Value) != val {
		return fmt.Errorf("expected %q, got %q", val, string(item.Value))
	}
	return nil
}

func performAdd(mc *memcache.Client) error {
	key := "addkey"
	mc.Delete(key)
	if err := mc.Add(&memcache.Item{Key: key, Value: []byte("val1")}); err != nil {
		return fmt.Errorf("add: %w", err)
	}
	if err := mc.Add(&memcache.Item{Key: key, Value: []byte("val2")}); err == nil {
		return fmt.Errorf("add should fail on existing key")
	}
	return nil
}

func performReplace(mc *memcache.Client) error {
	key := "replacekey"
	mc.Delete(key)
	if err := mc.Replace(&memcache.Item{Key: key, Value: []byte("val1")}); err == nil {
		return fmt.Errorf("replace should fail on missing key")
	}
	mc.Set(&memcache.Item{Key: key, Value: []byte("val1")})
	if err := mc.Replace(&memcache.Item{Key: key, Value: []byte("val2")}); err != nil {
		return fmt.Errorf("replace: %w", err)
	}
	item, err := mc.Get(key)
	if err != nil || string(item.Value) != "val2" {
		return fmt.Errorf("replace did not update value")
	}
	return nil
}

func performDelete(mc *memcache.Client) error {
	key := "delkey"
	mc.Set(&memcache.Item{Key: key, Value: []byte("gone")})
	if err := mc.Delete(key); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	if _, err := mc.Get(key); err == nil {
		return fmt.Errorf("delete did not remove key")
	}
	return nil
}
