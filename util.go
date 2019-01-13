package framework

import (
	"math/rand"
	"os"
	"time"
)

// PathExist check wether file/dir is exist
func PathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

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
