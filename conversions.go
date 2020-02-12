package conf

import (
	"github.com/buger/jsonparser"
	"log"
	"strconv"
	"sync"
)

// String возвращает строку из конфига по набору ключей
func String(keys ...string) (string, bool) {
	val, ok := Get(keys...)
	if ok {
		return string(val), true
	}
	return "", false
}

// FirstString для списка наборов ключей возвращает первое значение найденное в конфиге
// Функции с префиксом First- придуманы для удобства имплементации дефолтных значений
func FirstString(keysList ...[]string) (string, bool) {
	for _, key := range keysList {
		if value, ok := String(key...); ok {
			return value, true
		}
	}
	return "", false
}

var stringArrayCache map[string][]string
var stringArrayMutex = sync.Mutex{}

// StringArray возвращает массив строк из конфига по набору ключей
func StringArray(keys ...string) ([]string, bool) {
	hash := getKeysHash(keys...)

	stringArrayMutex.Lock()
	vals, ok := stringArrayCache[hash]
	stringArrayMutex.Unlock()
	if ok {
		return vals, true
	}

	val, ok := Get(keys...)
	if ok {
		vals := []string{}
		// Ошибки тут не важны
		_, _ = jsonparser.ArrayEach(val, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			vals = append(vals, string(value))
		})

		stringArrayMutex.Lock()
		stringArrayCache[hash] = vals
		stringArrayMutex.Unlock()

		return vals, true
	}
	return []string{}, false
}

// StringArray возвращает первый найденный в конфиге массив строк по набору ключей
func FirstStringArray(keysList ...[]string) ([]string, bool) {
	for _, key := range keysList {
		if value, ok := StringArray(key...); ok {
			return value, true
		}
	}
	return []string{}, false
}

// StringMap возвращает словарь строк по набору ключей
func StringMap(keys ...string) (map[string]string, bool) {
	// todo: cache?
	vals := map[string]string{}
	val, ok := Get(keys...)
	if !ok {
		return vals, false
	}

	// Ошибки тут можно игнорировать
	_ = jsonparser.ObjectEach(val, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		vals[string(key)] = string(value)
		return nil
	})

	return vals, true
}

// Uint64 возвращает число из конфига по набору ключей
func Uint64(keys ...string) (uint64, bool) {
	str, ok := String(keys...)
	if !ok {
		return 0, false
	}
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		log.Println("Error: Value", keys, "is not int")
		return 0, false
	}
	return val, true
}

// FirstUint64 возвращает первое найденное число из конфига по списку наборов ключей
func FirstUint64(keysList ...[]string) (uint64, bool) {
	for _, key := range keysList {
		if value, ok := Uint64(key...); ok {
			return value, true
		}
	}
	return 0, false
}

// Uint32 возвращает число из конфига по набору ключей
func Uint32(keys ...string) (uint32, bool) {
	val, ok := Uint64(keys...)
	if ok {
		return uint32(val), true
	}
	return 0, false
}

// FirstUint32 возвращает первое найденное число из конфига по списку наборов ключей
func FirstUint32(keysList ...[]string) (uint32, bool) {
	for _, key := range keysList {
		if value, ok := Uint32(key...); ok {
			return value, true
		}
	}
	return 0, false
}

// Int64 возвращает число из конфига по набору ключей
func Int64(keys ...string) (int64, bool) {
	str, ok := String(keys...)
	if !ok {
		return 0, false
	}
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Println("Error: Value", keys, "is not int")
		return 0, false
	}
	return val, true
}

// FirstInt64 возвращает первое найденное число из конфига по списку наборов ключей
func FirstInt64(keysList ...[]string) (int64, bool) {
	for _, key := range keysList {
		if value, ok := Int64(key...); ok {
			return value, true
		}
	}
	return 0, false
}

// Int32 возвращает число из конфига по набору ключей
func Int32(keys ...string) (int32, bool) {
	val, ok := Int64(keys...)
	if ok {
		return int32(val), true
	}
	return 0, false
}

// FirstInt32 возвращает первое найденное число из конфига по списку наборов ключей
func FirstInt32(keysList ...[]string) (int32, bool) {
	for _, key := range keysList {
		if value, ok := Int32(key...); ok {
			return value, true
		}
	}
	return 0, false
}

// Int8 возвращает число из конфига по набору ключей
func Int8(keys ...string) (int8, bool) {
	val, ok := Int64(keys...)
	if ok {
		return int8(val), true
	}
	return 0, false
}

// FirstInt8 возвращает первое найденное число из конфига по списку наборов ключей
func FirstInt8(keysList ...[]string) (int8, bool) {
	for _, key := range keysList {
		if value, ok := Int8(key...); ok {
			return value, true
		}
	}
	return 0, false
}

// Int возвращает число из конфига по набору ключей
func Int(keys ...string) (int, bool) {
	val, ok := Int64(keys...)
	if ok {
		return int(val), true
	}
	return 0, false
}

// FirstInt возвращает первое найденное число из конфига по списку наборов ключей
func FirstInt(keysList ...[]string) (int, bool) {
	for _, key := range keysList {
		if value, ok := Int(key...); ok {
			return value, true
		}
	}
	return 0, false
}

// Float64 возвращает число из конфига по набору ключей
func Float64(keys ...string) (float64, bool) {
	str, ok := String(keys...)
	if !ok {
		return 0, false
	}
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Println("Error: Value", keys, "is not float")
		return 0, false
	}
	return val, true
}

// FirstFloat64 возвращает первое найденное число из конфига по списку наборов ключей
func FirstFloat64(keysList ...[]string) (float64, bool) {
	for _, key := range keysList {
		if value, ok := Float64(key...); ok {
			return value, true
		}
	}
	return 0, false
}

// Float32 возвращает число из конфига по набору ключей
func Float32(keys ...string) (float32, bool) {
	val, ok := Float64(keys...)
	if ok {
		return float32(val), true
	}
	return 0, false
}

// FirstFloat32 возвращает первое найденное число из конфига по списку наборов ключей
func FirstFloat32(keysList ...[]string) (float32, bool) {
	for _, key := range keysList {
		if value, ok := Float32(key...); ok {
			return value, true
		}
	}
	return 0, false
}
