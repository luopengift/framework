package json

// Format implement Formatter interface
// map[string]interface{} => struct{}
// eg: Format(map[string]interface{...}, &Struct{})
func Format(in, out interface{}) error {
	b, err := Marshal(in)
	if err != nil {
		return err
	}
	return Unmarshal(b, out)
}
