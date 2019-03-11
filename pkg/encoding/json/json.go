// +build !jsoniter

package json

import (
	"encoding/json"
	"io"
)

func NewDecoder(r io.Reader) *json.Decoder {
	return json.NewDecoder(r)
}

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}
