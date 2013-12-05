package imageserver

import (
	"bytes"
	"fmt"
	"sort"
)

// Parameters represents parameters used in imageserver package
//
// This is a wrapper around map and provides getter and hash methods
//
// Getter methods return an error if the key does not exist or the type does not match
type Parameters map[string]interface{}

// Set sets the value for the key
func (parameters Parameters) Set(key string, value interface{}) {
	parameters[key] = value
}

// Has returns true if the key exists and false otherwise
func (parameters Parameters) Has(key string) bool {
	_, ok := parameters[key]
	return ok
}

// Len returns the length
func (parameters Parameters) Len() int {
	return len(parameters)
}

// Empty returns true if parameters is empty and false otherwise
func (parameters Parameters) Empty() bool {
	return parameters.Len() == 0
}

// Keys returns the keys
func (parameters Parameters) Keys() []string {
	length := parameters.Len()
	keys := make([]string, length)
	i := 0
	for key := range parameters {
		keys[i] = key
		i++
	}
	return keys
}

// Get returns the value for the key
//
// It returns an error if the value is not found
func (parameters Parameters) Get(key string) (interface{}, error) {
	value, found := parameters[key]
	if !found {
		return nil, fmt.Errorf("value not found for key %s", key)
	}
	return value, nil
}

// GetString returns the value as a string for the key
//
// It returns an error if the value is not a string
func (parameters Parameters) GetString(key string) (string, error) {
	v, err := parameters.Get(key)
	if err != nil {
		return "", err
	}
	value, ok := v.(string)
	if !ok {
		return value, parameters.newErrorType(key, v, "string")
	}
	return value, nil
}

// GetInt returns the value as an int for the key
//
// It returns an error if the value is not an int
func (parameters Parameters) GetInt(key string) (int, error) {
	v, err := parameters.Get(key)
	if err != nil {
		return 0, err
	}
	value, ok := v.(int)
	if !ok {
		return value, parameters.newErrorType(key, v, "int")
	}
	return value, nil
}

// GetBool returns the value as a bool for the key
//
// It returns an error if the value is not a bool
func (parameters Parameters) GetBool(key string) (bool, error) {
	v, err := parameters.Get(key)
	if err != nil {
		return false, err
	}
	value, ok := v.(bool)
	if !ok {
		return value, parameters.newErrorType(key, v, "bool")
	}
	return value, nil
}

// GetParameters returns the value as a Parameters for the key
//
// It returns an error if the value is not a Parameters
func (parameters Parameters) GetParameters(key string) (Parameters, error) {
	v, err := parameters.Get(key)
	if err != nil {
		return nil, err
	}
	value, ok := v.(Parameters)
	if !ok {
		return value, parameters.newErrorType(key, v, "Parameters")
	}
	return value, nil
}

func (parameters Parameters) newErrorType(key string, value interface{}, expectedType string) error {
	return fmt.Errorf("value %#v for key %s is of type %T instead of %s", value, key, value, expectedType)
}

// String returns the string representation
//
// Keys are sorted alphabetically
func (parameters Parameters) String() string {
	keys := parameters.Keys()
	sort.Strings(keys)

	buffer := new(bytes.Buffer)
	buffer.WriteString("map[")
	for i, key := range keys {
		if i != 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteString(key)
		buffer.WriteString(":")
		buffer.WriteString(fmt.Sprint(parameters[key]))
	}
	buffer.WriteString("]")

	return buffer.String()
}
