package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	err := Load(struct{}{})
	assert.EqualError(t, err, "load(non-pointer {})")

	var vPtr *uint
	err = Load(vPtr)
	assert.EqualError(t, err, "load pointer nil <nil>")

	type test struct {
		Field    string `env:""`
		FieldInt int    `env:"test_int"`
	}
	err = Load(&test{})
	assert.NoError(t, err)

	os.Setenv("test_int", "GGG")
	err = Load(&test{})
	assert.NoError(t, err)

}

func TestLoadString(t *testing.T) {
	type test struct {
		Field        string  `env:"test_f1"`
		FieldPtr     *string `env:"test_f1"`
		UnknownField string
	}

	v := &test{}
	err := Load(v)
	assert.NoError(t, err)
	assert.Empty(t, v.Field)

	wantEnv := "ENV_VALUE_F1"
	os.Setenv("test_f1", wantEnv)
	err = Load(v)
	assert.NoError(t, err)
	assert.Equal(t, wantEnv, v.Field)
	assert.Equal(t, wantEnv, *v.FieldPtr)

	wantEnv = "ENV_VALUE_F2"
	os.Setenv("test_f1", wantEnv)
	err = Load(v)
	assert.NoError(t, err)
	assert.Equal(t, wantEnv, *v.FieldPtr)
}

func TestLoadUint(t *testing.T) {
	type test struct {
		Field    uint  `env:"test_f1"`
		FieldPtr *uint `env:"test_f1"`

		Field8    uint8  `env:"test_f1"`
		Field8Ptr *uint8 `env:"test_f1"`

		Field16    uint16  `env:"test_f1"`
		Field16Ptr *uint16 `env:"test_f1"`

		Field32    uint32  `env:"test_f1"`
		Field32Ptr *uint32 `env:"test_f1"`

		Field64    uint64  `env:"test_f1"`
		Field64Ptr *uint64 `env:"test_f1"`
	}

	v := &test{}
	err := Load(v)
	assert.NoError(t, err)
	assert.Empty(t, v.Field)

	wantEnv := "250"
	os.Setenv("test_f1", wantEnv)
	err = Load(v)
	assert.NoError(t, err)
	assert.Equal(t, uint(250), v.Field)
	assert.Equal(t, uint(250), *v.FieldPtr)

	assert.Equal(t, uint8(250), v.Field8)
	assert.Equal(t, uint8(250), *v.Field8Ptr)

	assert.Equal(t, uint16(250), v.Field16)
	assert.Equal(t, uint16(250), *v.Field16Ptr)

	assert.Equal(t, uint32(250), v.Field32)
	assert.Equal(t, uint32(250), *v.Field32Ptr)

	assert.Equal(t, uint64(250), v.Field64)
	assert.Equal(t, uint64(250), *v.Field64Ptr)

}

func TestLoadInt(t *testing.T) {
	type test struct {
		Field    int  `env:"test_f1"`
		FieldPtr *int `env:"test_f1"`

		Field8    int8  `env:"test_f1"`
		Field8Ptr *int8 `env:"test_f1"`

		Field16    int16  `env:"test_f1"`
		Field16Ptr *int16 `env:"test_f1"`

		Field32    int32  `env:"test_f1"`
		Field32Ptr *int32 `env:"test_f1"`

		Field64    int64  `env:"test_f1"`
		Field64Ptr *int64 `env:"test_f1"`
	}

	v := &test{}
	os.Setenv("test_f1", "")
	err := Load(v)
	assert.NoError(t, err)
	assert.Empty(t, v.Field)

	wantEnv := "102"
	os.Setenv("test_f1", wantEnv)
	err = Load(v)
	assert.NoError(t, err)
	assert.Equal(t, int(102), v.Field)
	assert.Equal(t, int(102), *v.FieldPtr)

	assert.Equal(t, int8(102), v.Field8)
	assert.Equal(t, int8(102), *v.Field8Ptr)

	assert.Equal(t, int16(102), v.Field16)
	assert.Equal(t, int16(102), *v.Field16Ptr)

	assert.Equal(t, int32(102), v.Field32)
	assert.Equal(t, int32(102), *v.Field32Ptr)

	assert.Equal(t, int64(102), v.Field64)
	assert.Equal(t, int64(102), *v.Field64Ptr)

}

func TestLoadBool(t *testing.T) {
	type test struct {
		Field    bool  `env:"test_f1"`
		FieldPtr *bool `env:"test_f1"`
	}

	os.Setenv("test_f1", "")
	v := &test{}
	err := Load(v)
	assert.NoError(t, err)
	assert.Empty(t, v.Field)

	for _, testCase := range []string{"0", "false", ""} {
		os.Setenv("test_f1", testCase)
		err = Load(v)
		assert.NoError(t, err)
		assert.Equal(t, false, v.Field)
		assert.Equal(t, false, *v.FieldPtr)
	}

	for _, testCase := range []string{"1", "true", "zzz", " "} {
		os.Setenv("test_f1", testCase)
		err = Load(v)
		assert.NoError(t, err)
		assert.Equal(t, true, v.Field)
		assert.Equal(t, true, *v.FieldPtr)
	}

}

func TestLoadFloat(t *testing.T) {
	type test struct {
		Field32    float32  `env:"test_f1"`
		FieldPtr32 *float32 `env:"test_f1"`

		Field64    float64  `env:"test_f1"`
		FieldPtr64 *float64 `env:"test_f1"`
	}

	v := &test{}
	err := Load(v)
	assert.NoError(t, err)
	assert.Zero(t, v.Field32)
	assert.Nil(t, v.FieldPtr32)
	assert.Zero(t, v.Field64)
	assert.Nil(t, v.FieldPtr64)

	wantEnv := "3.14"
	os.Setenv("test_f1", wantEnv)
	err = Load(v)
	assert.NoError(t, err)
	assert.Equal(t, float32(3.14), v.Field32)
	assert.Equal(t, float32(3.14), *v.FieldPtr32)

	assert.Equal(t, float64(3.14), v.Field64)
	assert.Equal(t, float64(3.14), *v.FieldPtr64)
}

type testUnmarshalerType int
type testEnvUnmarshaler struct {
	T      *testing.T
	Field  string  `env:"test"`
	Field2 *string `env:"test2"`
	Field3 *uint   `env:"test3"`
	Field4 string
	Field5 testUnmarshalerType  `env:"test_um_1"`
	Field6 *testUnmarshalerType `env:"test_um_2"`
}

func (t testEnvUnmarshaler) UnmarshalField(field string, name, value string) (interface{}, error) {
	if field == "Field" {
		return name, nil
	}
	if field == "Field2" {
		s := name
		return &s, nil
	}
	if field == "Field3" {
		return nil, nil
	}

	if name == "test_um_1" {
		return testUnmarshalerType(1), nil
	}
	if name == "test_um_2" {
		v := testUnmarshalerType(2)
		return &v, nil
	}

	assert.Failf(t.T, "Expected call %s", field)
	return nil, nil
}

func TestLoadWithUnmarshalerFields(t *testing.T) {
	v := &testEnvUnmarshaler{T: t}
	err := Load(v)
	assert.NoError(t, err)
	assert.Equal(t, "test", v.Field)
	assert.Equal(t, "test2", *v.Field2)

	assert.Equal(t, (*uint)(nil), v.Field3)
	assert.Equal(t, testUnmarshalerType(1), v.Field5)
	assert.Equal(t, testUnmarshalerType(2), *v.Field6)

	assert.Nil(t, v.Field3)
}
