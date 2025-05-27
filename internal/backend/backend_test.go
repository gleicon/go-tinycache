package backend

import (
	"fmt"
	"math/rand"
)

var vboltdb *KVBoltDBBackend

func init() {
	vboltdb, _ = NewKVBoltDBBackend("bolt.db", "memcached", 10000)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func errUnexpected(msg interface{}) string {
	return fmt.Sprintf("Unexpected response: %#v\n", msg)
}
