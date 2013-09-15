package imageserver

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

// Parameters used for provider, processor, cache, ...
//
// This is a wrapper around map and provides getter and hash methods
type Parameters map[string]interface{}

func (parameters Parameters) Set(key string, value interface{}) {
	parameters[key] = value
}

func (parameters Parameters) Get(key string) (interface{}, error) {
	value, found := parameters[key]
	if !found {
		err := fmt.Errorf("Value not found")
		return nil, err
	}
	return value, nil
}

func (parameters Parameters) GetString(key string) (string, error) {
	v, err := parameters.Get(key)
	if err != nil {
		return "", err
	}
	value, ok := v.(string)
	if !ok {
		err = fmt.Errorf("Not a string")
		return "", err
	}
	return value, nil
}

func (parameters Parameters) GetInt(key string) (int, error) {
	v, err := parameters.Get(key)
	if err != nil {
		return 0, err
	}
	value, ok := v.(int)
	if !ok {
		err = fmt.Errorf("Not an int")
		return 0, err
	}
	return value, nil
}

func (parameters Parameters) GetBool(key string) (bool, error) {
	v, err := parameters.Get(key)
	if err != nil {
		return false, err
	}
	value, ok := v.(bool)
	if !ok {
		err = fmt.Errorf("Not a bool")
		return false, err
	}
	return value, nil
}

// Hash content with sha256 algorithm and returns a string
func (parameters Parameters) Hash() string {
	hash := sha256.New()
	io.WriteString(hash, fmt.Sprint(parameters))
	data := hash.Sum(nil)
	return hex.EncodeToString(data)
}
