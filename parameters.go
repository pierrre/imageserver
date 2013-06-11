package imageserver

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

type Parameters map[string]interface{}

func (parameters Parameters) Set(key string, value interface{}) {
	parameters[key] = value
}

func (parameters Parameters) Get(key string) (value interface{}, err error) {
	value, ok := parameters[key]
	if !ok {
		err = fmt.Errorf("Value not found")
	}
	return
}

func (parameters Parameters) GetString(key string) (value string, err error) {
	v, err := parameters.Get(key)
	if err != nil {
		return
	}
	value, ok := v.(string)
	if !ok {
		err = fmt.Errorf("Not a string")
	}
	return
}

func (parameters Parameters) GetInt(key string) (value int, err error) {
	v, err := parameters.Get(key)
	if err != nil {
		return
	}
	value, ok := v.(int)
	if !ok {
		err = fmt.Errorf("Not an int")
	}
	return
}

func (parameters Parameters) GetBool(key string) (value bool, err error) {
	v, err := parameters.Get(key)
	if err != nil {
		return
	}
	value, ok := v.(bool)
	if !ok {
		err = fmt.Errorf("Not a bool")
	}
	return
}

func (parameters Parameters) Hash() string {
	hash := sha256.New()
	io.WriteString(hash, fmt.Sprint(parameters))
	data := hash.Sum(nil)
	hexaData := hex.EncodeToString(data)
	return hexaData
}
