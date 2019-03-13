package util

import (
	"math/rand"
	"time"
)

// Random random string
func Random(n int) string {
	bytes := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	result := make([]byte, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		result[i] = bytes[r.Intn(len(bytes))]
	}
	return string(result)
}
