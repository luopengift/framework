// +build !jsoniter

package json

import (
	"encoding/json"
	"io"
)

// MarshalToString convenient method to write as string instead of []byte
func MarshalToString(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}

// Marshal returns the JSON encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalIndent is like Marshal but applies Indent to format the output.
// Each JSON element in the output will begin on a new line beginning with prefix followed by one or more copies of indent according to the indentation nesting.
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// UnmarshalFromString convenient method to read from string instead of []byte
func UnmarshalFromString(str string, v interface{}) error {
	return json.Unmarshal([]byte(str), v)
}

// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
// If v is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError.
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// NewEncoder returns a new encoder that writes to writer.
func NewEncoder(writer io.Writer) *json.Encoder {
	return json.NewEncoder(writer)
}

// NewDecoder returns a new decoder that reads from r.
// The decoder introduces its own buffering and may read data from r beyond the JSON values requested.
func NewDecoder(r io.Reader) *json.Decoder {
	return json.NewDecoder(r)
}

// Valid reports whether data is a valid JSON encoding.
func Valid(data []byte) bool {
	return json.Valid(data)
}
