package env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

//Errors
var (
	ErrBadSyntax = errors.New("bad syntax")
)

//Unmarshaler is the interface for load from environment
type Unmarshaler interface {
	UnmarshalField(field, name, value string) (interface{}, error)
}

//Load loads values from environment variables
func Load(v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("load(non-pointer %v)", v)
	}

	if value.IsNil() {
		return fmt.Errorf("load pointer nil %v", v)
	}

	return walk(value, nil, nil)
}

var bits = map[reflect.Kind]int{
	reflect.Uint:   64,
	reflect.Uint8:  8,
	reflect.Uint16: 16,
	reflect.Uint32: 32,
	reflect.Uint64: 64,

	reflect.Int:   64,
	reflect.Int8:  8,
	reflect.Int16: 16,
	reflect.Int32: 32,
	reflect.Int64: 64,
}

var castInts = map[reflect.Kind]func(interface{}) reflect.Value{
	reflect.Uint: func(v interface{}) reflect.Value {
		out := uint(v.(uint64))
		return reflect.ValueOf(&out)
	},
	reflect.Uint8: func(v interface{}) reflect.Value {
		out := uint8(v.(uint64))
		return reflect.ValueOf(&out)
	},
	reflect.Uint16: func(v interface{}) reflect.Value {
		out := uint16(v.(uint64))
		return reflect.ValueOf(&out)
	},
	reflect.Uint32: func(v interface{}) reflect.Value {
		out := uint32(v.(uint64))
		return reflect.ValueOf(&out)
	},
	reflect.Uint64: func(v interface{}) reflect.Value {
		out := v.(uint64)
		return reflect.ValueOf(&out)
	},
	reflect.Int: func(v interface{}) reflect.Value {
		out := int(v.(int64))
		return reflect.ValueOf(&out)
	},
	reflect.Int8: func(v interface{}) reflect.Value {
		out := int8(v.(int64))
		return reflect.ValueOf(&out)
	},
	reflect.Int16: func(v interface{}) reflect.Value {
		out := int16(v.(int64))
		return reflect.ValueOf(&out)
	},
	reflect.Int32: func(v interface{}) reflect.Value {
		out := int32(v.(int64))
		return reflect.ValueOf(&out)
	},
	reflect.Int64: func(v interface{}) reflect.Value {
		out := v.(int64)
		return reflect.ValueOf(&out)
	},
}

func unmarshalFromValue(envName string, strField *reflect.StructField, str *reflect.Value) (interface{}, error) {
	if str != nil && str.IsValid() && str.CanInterface() {
		if unmarshaler, ok := str.Interface().(Unmarshaler); ok {
			valUnmarshal, err := unmarshaler.UnmarshalField(strField.Name, envName, os.Getenv(envName))
			if err != nil {
				return nil, err
			}
			if valUnmarshal != nil {
				return valUnmarshal, nil
			}
		}
	}
	return nil, nil
}

//nolint: gocyclo
func loadValueEnv(value reflect.Value, strField *reflect.StructField, str *reflect.Value) error {
	if strField == nil {
		return nil
	}

	envName, ok := strField.Tag.Lookup("env")
	if !ok {
		return nil
	}

	var (
		envValueAsIs interface{}
		envValue     string
	)

	valUnmarshal, err := unmarshalFromValue(envName, strField, str)
	if err != nil {
		return err
	}
	if valUnmarshal != nil {
		envValueAsIs = valUnmarshal
	}

	kind := strField.Type.Kind()
	ptr := false
	if kind == reflect.Ptr {
		kind = strField.Type.Elem().Kind()
		ptr = true
	}

	if envValueAsIs == nil {
		envValueStr, ok := os.LookupEnv(envName)
		if !ok {
			return nil
		}
		envValue = envValueStr
	} else {
		if ptr {
			str.FieldByName(strField.Name).Set(reflect.ValueOf(envValueAsIs))
			return nil
		}
		value.Set(reflect.ValueOf(envValueAsIs).Convert(value.Type()))
		return nil
	}

	if err := setValue(ptr, str, strField, envValue, value, kind); err != nil {
		return err
	}

	return nil
}

func setValue(ptr bool, str *reflect.Value, strField *reflect.StructField, envValue string, value reflect.Value, kind reflect.Kind) error {
	switch kind {
	case reflect.Bool:
		setBool(envValue, ptr, str, strField, value)
	case reflect.String:
		setString(ptr, str, strField, envValue, value)
	case reflect.Float64:
		if err := setFloat64(ptr, str, strField, envValue, value); err != nil {
			return err
		}
	case reflect.Float32:
		if err := setFloat32(ptr, str, strField, envValue, value); err != nil {
			return err
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if err := setInt(ptr, str, strField, envValue, value, kind); err != nil {
			return err
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if err := setUInt(ptr, str, strField, envValue, value, kind); err != nil {
			return err
		}
	}
	return nil
}

func setUInt(ptr bool, str *reflect.Value, strField *reflect.StructField, envValue string, value reflect.Value, kind reflect.Kind) error {
	v, err := strconv.ParseUint(envValue, 10, bits[kind])
	if err != nil {
		return nil
	}
	if ptr {
		vCast := castInts[kind](v)
		str.FieldByName(strField.Name).Set(vCast)
	} else {
		value.SetUint(v)
	}
	return nil
}
func setInt(ptr bool, str *reflect.Value, strField *reflect.StructField, envValue string, value reflect.Value, kind reflect.Kind) error {
	v, err := strconv.ParseInt(envValue, 10, bits[kind])
	if err != nil {
		return nil
	}
	if ptr {
		vCast := castInts[kind](v)
		str.FieldByName(strField.Name).Set(vCast)
	} else {
		value.SetInt(v)
	}
	return nil
}

func setFloat32(ptr bool, str *reflect.Value, strField *reflect.StructField, envValue string, value reflect.Value) error {
	v, err := strconv.ParseFloat(envValue, 32)
	v32 := float32(v)
	if err != nil {
		return err
	}
	if ptr {
		str.FieldByName(strField.Name).Set(reflect.ValueOf(&v32))
	} else {
		value.SetFloat(v)
	}
	return nil
}

func setFloat64(ptr bool, str *reflect.Value, strField *reflect.StructField, envValue string, value reflect.Value) error {
	v, err := strconv.ParseFloat(envValue, 64)
	if err != nil {
		return err
	}
	if ptr {
		str.FieldByName(strField.Name).Set(reflect.ValueOf(&v))
	} else {
		value.SetFloat(v)
	}
	return nil
}

func setString(ptr bool, str *reflect.Value, strField *reflect.StructField, envValue string, value reflect.Value) {
	if ptr {
		str.FieldByName(strField.Name).Set(reflect.ValueOf(&envValue))
	} else {
		value.SetString(envValue)
	}
}

func setBool(envValue string, ptr bool, str *reflect.Value, strField *reflect.StructField, value reflect.Value) {
	v := envValue != "" && envValue != "0" && strings.ToLower(envValue) != "false"
	if ptr {
		str.FieldByName(strField.Name).Set(reflect.ValueOf(&v))
	} else {
		value.SetBool(v)
	}
}

func walk(value reflect.Value, strField *reflect.StructField, str *reflect.Value) error {
	value = reflect.Indirect(value)

	if err := loadValueEnv(value, strField, str); err != nil {
		return nil
	}

	if value.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < value.NumField(); i++ {

		cType := value.Type()
		strFieldCur := cType.Field(i)

		if err := walk(value.Field(i), &strFieldCur, &value); err != nil {
			return err
		}
	}

	return nil
}
