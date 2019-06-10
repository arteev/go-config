package config

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/arteev/go-config/env"
)

//Unmarshaler is the interface parses the encoded data and stores the result
type Unmarshaler interface {
	Unmarshal(data []byte, v interface{}) error
}

//UnmarshalerFunc use for cast func to Unmarshaler
type UnmarshalerFunc func(data []byte, v interface{}) error

//Unmarshal implements Unmarshaler
func (f UnmarshalerFunc) Unmarshal(data []byte, v interface{}) error {
	return f(data, v)
}

var (
	readerFile = ioutil.ReadFile
)

//Mode use field in config: ConfigMode Mode for auto setup
type Mode int

//Modes
const (
	ModeUnknown Mode = iota
	ModeFile
	ModeEnvironment
)

//LoadFromFile load config from file using unmarshaler
func LoadFromFile(v interface{}, fileName string, unmarshaler Unmarshaler) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("LoadFromFile(non-pointer %v)", v)
	}

	data, err := readerFile(fileName)
	if err != nil {
		return err
	}

	err = unmarshaler.Unmarshal(data, v)
	if err != nil {
		return err
	}
	setMode(value, ModeFile)
	return nil
}

//LoadFromEnv load from environment variables
func LoadFromEnv(v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("LoadFromFile(non-pointer %v)", v)
	}
	err := env.Load(v)
	if err != nil {
		return nil
	}
	setMode(value, ModeEnvironment)
	return nil
}

func setMode(v reflect.Value, mode Mode) {
	v = reflect.Indirect(v)
	kind := v.Kind()
	if kind == reflect.Struct {
		vMode := v.FieldByName("ConfigMode")
		if !vMode.IsValid() || vMode.Type() != reflect.TypeOf(mode) {
			return
		}
		vMode.Set(reflect.ValueOf(mode))
	}
}

//SetReaderFile replaces readerFile for tests
func SetReaderFile(fn func(string) ([]byte, error)) {
	readerFile = fn
}

func (m Mode) String() string {
	switch m {
	case ModeFile:
		return "file"
	case ModeEnvironment:
		return "environment"
	default:
		return "unknown"
	}
}
