package imageserver

import (
	"bytes"
	"fmt"
	"sort"
)

// Params represents params used in imageserver package
//
// This is a wrapper around map and provides getter and hash methods
//
// Getter methods return an error if the key does not exist or the type does not match
type Params map[string]interface{}

// Set sets the value for the key
func (params Params) Set(key string, value interface{}) {
	params[key] = value
}

// Has returns true if the key exists and false otherwise
func (params Params) Has(key string) bool {
	_, ok := params[key]
	return ok
}

// Len returns the length
func (params Params) Len() int {
	return len(params)
}

// Empty returns true if params is empty and false otherwise
func (params Params) Empty() bool {
	return params.Len() == 0
}

// Keys returns the keys
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

// Get returns the value for the key
//
// It returns an error if the value is not found
func (params Params) Get(key string) (interface{}, error) {
	value, found := params[key]
	if !found {
		return nil, fmt.Errorf("value not found for key %s", key)
	}
	return value, nil
}

// GetString returns the value as a string for the key
//
// It returns an error if the value is not a string
func (params Params) GetString(key string) (string, error) {
	v, err := params.Get(key)
	if err != nil {
		return "", err
	}
	value, ok := v.(string)
	if !ok {
		return value, params.newErrorType(key, v, "string")
	}
	return value, nil
}

// GetInt returns the value as an int for the key
//
// It returns an error if the value is not an int
func (params Params) GetInt(key string) (int, error) {
	v, err := params.Get(key)
	if err != nil {
		return 0, err
	}
	value, ok := v.(int)
	if !ok {
		return value, params.newErrorType(key, v, "int")
	}
	return value, nil
}

// GetBool returns the value as a bool for the key
//
// It returns an error if the value is not a bool
func (params Params) GetBool(key string) (bool, error) {
	v, err := params.Get(key)
	if err != nil {
		return false, err
	}
	value, ok := v.(bool)
	if !ok {
		return value, params.newErrorType(key, v, "bool")
	}
	return value, nil
}

// GetParams returns the value as a Params for the key
//
// It returns an error if the value is not a Params
func (params Params) GetParams(key string) (Params, error) {
	v, err := params.Get(key)
	if err != nil {
		return nil, err
	}
	value, ok := v.(Params)
	if !ok {
		return value, params.newErrorType(key, v, "Params")
	}
	return value, nil
}

func (params Params) newErrorType(key string, value interface{}, expectedType string) error {
	return fmt.Errorf("value %#v for key %s is of type %T instead of %s", value, key, value, expectedType)
}

// String returns the string representation
//
// Keys are sorted alphabetically
func (params Params) String() string {
	keys := params.Keys()
	sort.Strings(keys)

	buffer := new(bytes.Buffer)
	buffer.WriteString("map[")
	for i, key := range keys {
		if i != 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteString(key)
		buffer.WriteString(":")
		buffer.WriteString(fmt.Sprint(params[key]))
	}
	buffer.WriteString("]")

	return buffer.String()
}

// ParamError is an error for a param
type ParamError struct {
	Param   string // Nested param path uses "." as separator
	Message string
}

func (err *ParamError) Error() string {
	return fmt.Sprintf("invalid param \"%s\": %s", err.Param, err.Message)
}
