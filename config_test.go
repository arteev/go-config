package config

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFromEnv(t *testing.T) {
	v := struct {
		Addr       string `json:"addr" env:"test_env1"`
		ConfigMode Mode
	}{}

	//error Ptr
	err := LoadFromEnv(struct{}{})
	assert.EqualError(t, err, "LoadFromFile(non-pointer {})")

	//Mode
	os.Setenv("test_env1", "server_addr")
	err = LoadFromEnv(&v)
	assert.Equal(t, ModeEnvironment, v.ConfigMode)
	assert.Equal(t, "server_addr", v.Addr)
	assert.NoError(t, err)

}

func TestLoadFromFile(t *testing.T) {
	v := struct {
		Addr       string `json:"addr"`
		ConfigMode Mode
	}{}

	vWoMode := struct {
		Addr       string `json:"addr"`
		ConfigMode int
	}{}

	readerFile = func(fileName string) ([]byte, error) {
		assert.Equal(t, fileName, "file.json")
		return []byte(`{
"addr":"test"
}`), nil
	}

	err := LoadFromFile(&v, "file.json", UnmarshalerFunc(json.Unmarshal))
	assert.NoError(t, err)
	assert.Equal(t, v.Addr, "test")
	assert.Equal(t, ModeFile, v.ConfigMode)

	//w/o Mode
	err = LoadFromFile(&vWoMode, "file.json", UnmarshalerFunc(json.Unmarshal))
	assert.NoError(t, err)
	assert.Equal(t, 0, vWoMode.ConfigMode)

	//error Ptr
	err = LoadFromFile(struct{}{}, "file.json", UnmarshalerFunc(json.Unmarshal))
	assert.EqualError(t, err, "LoadFromFile(non-pointer {})")

	//error file
	wantErr := errors.New("reader err")
	readerFile = func(_ string) ([]byte, error) {
		return nil, wantErr
	}
	err = LoadFromFile(&v, "file.json", UnmarshalerFunc(json.Unmarshal))
	assert.EqualError(t, err, wantErr.Error())

	//error json
	readerFile = func(_ string) ([]byte, error) {
		return []byte(`{`), nil
	}
	err = LoadFromFile(&v, "file.json", UnmarshalerFunc(json.Unmarshal))
	assert.EqualError(t, err, "unexpected end of JSON input")
}

type testChecker struct {
	MustErr error
	Field   string `json:"field" check:"required"`
}

func (t testChecker) Check() error {
	return t.MustErr
}

func TestModeString(t *testing.T) {
	cases := map[Mode]string{
		ModeUnknown:     "unknown",
		ModeFile:        "file",
		ModeEnvironment: "environment",
		Mode(100):       "unknown",
	}

	for k, v := range cases {
		got := k.String()
		assert.Equal(t, v, got)
	}
}
