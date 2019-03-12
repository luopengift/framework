// +build !jsoniter

package json

import (
	"encoding/json"
	"io"
)

// NewDecoder returns a new decoder that reads from r.
// The decoder introduces its own buffering and may read data from r beyond the JSON values requested.
func NewDecoder(r io.Reader) *json.Decoder {
	return json.NewDecoder(r)
}

// Marshal returns the JSON encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalIndent is like Marshal but applies Indent to format the output. Each JSON element in the output will begin on a new line beginning with prefix followed by one or more copies of indent according to the indentation nesting.
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}
