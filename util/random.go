package util

import (
	"math/rand"
	"time"
)

var bytes = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

// Random random string
func Random(n int) string {
	result := make([]byte, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		result[i] = bytes[r.Intn(len(bytes))]
	}
	return string(result)
}
