package imageserver

import (
	"bytes"
	"fmt"
	"sort"
	"sync"
)

// Params are params used in imageserver.
//
// This is a wrapper around map[string]interface{} and provides utility methods.
// It should only contains basic Go types values (string, int float64, ...) or nested Params.
//
// Getter methods return a *ParamError if the key does not exist or the type does not match.
type Params map[string]interface{}

// Set sets the value for the key.
func (params Params) Set(key string, value interface{}) {
	params[key] = value
}

// Get returns the value for the key.
func (params Params) Get(key string) (interface{}, error) {
	v, ok := params[key]
	if !ok {
		return nil, &ParamError{Param: key, Message: "not set"}
	}
	return v, nil
}

// GetString returns the string value for the key.
func (params Params) GetString(key string) (string, error) {
	v, err := params.Get(key)
	if err != nil {
		return "", err
	}
	vt, ok := v.(string)
	if !ok {
		return vt, newErrorType(key, v, "string")
	}
	return vt, nil
}

// GetInt returns the int value for the key.
func (params Params) GetInt(key string) (int, error) {
	v, err := params.Get(key)
	if err != nil {
		return 0, err
	}
	vt, ok := v.(int)
	if !ok {
		return 0, newErrorType(key, v, "int")
	}
	return vt, nil
}

// GetFloat returns the float64 value for the key.
func (params Params) GetFloat(key string) (float64, error) {
	v, err := params.Get(key)
	if err != nil {
		return 0, err
	}
	vt, ok := v.(float64)
	if !ok {
		return 0, newErrorType(key, v, "float")
	}
	return vt, nil
}

// GetBool returns the bool value for the key.
func (params Params) GetBool(key string) (bool, error) {
	v, err := params.Get(key)
	if err != nil {
		return false, err
	}
	vt, ok := v.(bool)
	if !ok {
		return false, newErrorType(key, v, "bool")
	}
	return vt, nil
}

// GetParams returns the Params value for the key.
func (params Params) GetParams(key string) (Params, error) {
	v, err := params.Get(key)
	if err != nil {
		return nil, err
	}
	vt, ok := v.(Params)
	if !ok {
		return nil, newErrorType(key, v, "Params")
	}
	return vt, nil
}

func newErrorType(key string, value interface{}, expectedType string) error {
	return &ParamError{Param: key, Message: fmt.Sprintf("contains a value of type %T instead of %s", value, expectedType)}
}

// Has returns true if the key exists and false otherwise.
func (params Params) Has(key string) bool {
	_, ok := params[key]
	return ok
}

// Len returns the length.
func (params Params) Len() int {
	return len(params)
}

// Empty returns true if params is empty and false otherwise.
func (params Params) Empty() bool {
	return params.Len() == 0
}

// Keys returns the keys.
func (params Params) Keys() []string {
	length := params.Len()
	keys := make([]string, length)
	i := 0
	for key := range params {
		keys[i] = key
		i++
	}
	return keys
}

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// String returns the string representation.
//
// Keys are sorted alphabetically.
func (params Params) String() string {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)
	params.toBuffer(buf)
	return buf.String()
}

func (params Params) toBuffer(buf *bytes.Buffer) {
	keys := params.Keys()
	sort.Strings(keys)
	buf.WriteString("map[")
	for i, key := range keys {
		if i != 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(key)
		buf.WriteString(":")
		switch value := params[key].(type) {
		case Params:
			value.toBuffer(buf)
		default:
			fmt.Fprint(buf, value)
		}

	}
	buf.WriteString("]")
}

// ParamError is an error for a param.
type ParamError struct {
	Param   string // Nested param path uses "." as separator
	Message string
}

func (err *ParamError) Error() string {
	return fmt.Sprintf("invalid param \"%s\": %s", err.Param, err.Message)
}
