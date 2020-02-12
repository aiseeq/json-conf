package conf

import (
	"testing"
)

func initTestConfig(t *testing.T) {
	Init("test")
	if err := SetConfig([]byte(`{
  "app": "test",
  "key": {
    "subkey": "val",
    "float": 0.5,
    "int": 41
  },
  "strings": [
    "one",
    "two"
  ]
}`), "02-test"); err != nil {
		t.Error(err)
	}
	if err := SetConfig([]byte(`{
  "app": "test",
  "key": {
    "subkey": "val2",
    "int": 42
  },
  "strings": [
    "one",
    "two",
    "three"
  ]
}`), "01-test-rewrite"); err != nil {
		t.Error(err)
	}
}

func TestInt(t *testing.T) {
	initTestConfig(t)
	val, ok := Uint32("key", "int")
	if !ok {
		t.Error("Value not found")
	}
	if val != 42 {
		t.Error("Incorrect value")
	}
}

func TestFloat(t *testing.T) {
	initTestConfig(t)
	val, ok := Float32("key", "float")
	if !ok {
		t.Error("Value not found")
	}
	if val != 0.5 {
		t.Error("Incorrect value")
	}
}

func TestGetFirst(t *testing.T) {
	initTestConfig(t)
	val, ok := FirstString([]string{"nil"}, []string{"key", "int"}, []string{"app"})
	if !ok {
		t.Error("Value not found")
	}
	if val != "42" {
		t.Error("Incorrect value")
	}

	vals, ok := FirstStringArray([]string{"nil"}, []string{"strings"}, []string{"app"})
	if !ok {
		t.Error("Value not found")
	}
	for k, v := range []string{"one", "two", "three"} {
		if vals[k] != v {
			t.Errorf("Incorrect value: %s, expected: %s", vals[k], v)
		}
	}

	val2, ok := FirstUint64([]string{"nil"}, []string{"key", "int"}, []string{"app"})
	if !ok {
		t.Error("Value not found")
	}
	if val2 != 42 {
		t.Error("Incorrect value")
	}
	val3, ok := FirstUint32([]string{"nil"}, []string{"key", "int"}, []string{"app"})
	if !ok {
		t.Error("Value not found")
	}
	if val3 != 42 {
		t.Error("Incorrect value")
	}
	val4, ok := FirstInt64([]string{"nil"}, []string{"key", "int"}, []string{"app"})
	if !ok {
		t.Error("Value not found")
	}
	if val4 != 42 {
		t.Error("Incorrect value")
	}
	val5, ok := FirstInt32([]string{"nil"}, []string{"key", "int"}, []string{"app"})
	if !ok {
		t.Error("Value not found")
	}
	if val5 != 42 {
		t.Error("Incorrect value")
	}
	val6, ok := FirstInt8([]string{"nil"}, []string{"key", "int"}, []string{"app"})
	if !ok {
		t.Error("Value not found")
	}
	if val6 != 42 {
		t.Error("Incorrect value")
	}
	val7, ok := FirstInt([]string{"nil"}, []string{"key", "int"}, []string{"app"})
	if !ok {
		t.Error("Value not found")
	}
	if val7 != 42 {
		t.Error("Incorrect value")
	}

	val8, ok := FirstFloat64([]string{"nil"}, []string{"key", "int"}, []string{"app"})
	if !ok {
		t.Error("Value not found")
	}
	if val8 != 42 {
		t.Error("Incorrect value")
	}
	val9, ok := FirstFloat32([]string{"nil"}, []string{"key", "int"}, []string{"app"})
	if !ok {
		t.Error("Value not found")
	}
	if val9 != 42 {
		t.Error("Incorrect value")
	}
}

func TestStringMap(t *testing.T) {
	initTestConfig(t)
	vals, ok := StringMap("key")
	if !ok {
		t.Error("Value not found")
	}
	for k, v := range map[string]string{"subkey": "val2", "int": "42"} {
		if vals[k] != v {
			t.Errorf("Incorrect value: %s, expected: %s", vals[k], v)
		}
	}
}

func TestNoValue(t *testing.T) {
	initTestConfig(t)
	_, ok := String("nope")
	if ok {
		t.Error("Value found")
	}
	_, ok = String("nope")
	if ok {
		t.Error("Value found in cache")
	}
}
