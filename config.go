// Package conf позволяет удобно пользоваться конфигурационными файлами, написанными в формате JSON
package conf

import (
	"errors"
	"github.com/buger/jsonparser"
	"github.com/dgryski/go-farm"
	"github.com/orcaman/concurrent-map"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var configs cmap.ConcurrentMap // Список конфигов
var configNames []string       // Имена конфигов из списка, требуется для сортировки списка по ключу
var cache cmap.ConcurrentMap   // Кеш значений конфига
var appName string             // Имя приложения

// getKeysHash возвращает farm hash для произвольного набора строк, составляющих ключ
func getKeysHash(keys ...string) string {
	return strconv.FormatUint(farm.Hash64([]byte(strings.Join(keys, "|"))), 36)
}

// Get по списку ключей возвращает слайс байтов с требуемым куском JSON
func Get(keys ...string) ([]byte, bool) {
	// jsonparser может упасть и тут, к сожалению
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error: Recovered", r)
		}
	}()
	hash := getKeysHash(keys...)
	// Проверим сначала кеш
	if val, ok := cache.Get(hash); ok {
		switch v := val.(type) {
		case []byte:
			return v, true
		default: // nil - значит ключа нет в конфиге
			return nil, false
		}
	}
	// Перебираем имена конфигов (они отсортированы по возрастанию)
	for _, name := range configNames {
		// Получаем сам конфиг из кеша
		if data, ok := configs.Get(name); ok {
			// Это должен быть срез байт
			if config, ok := data.([]byte); ok {
				// Ищем в конфиге нужное значение
				if val, _, _, err := jsonparser.Get(config, keys...); err == nil {
					// Запишем в кеш и вернём
					cache.Set(hash, val)
					return val, true
				}
			}
		}
		// Если не смогли найти в этом конфиге, ищем в следующем
	}
	// А вот если совсем ничего не смогли найти, то запишем это в кеш
	cache.Set(hash, nil)
	return nil, false
}

// CheckJson проверяет валидность конфига
// Каждый конфиг должен иметь поле "app": appName для проверки json, и чтобы нельзя было случайно подсунуть другой json
func CheckJson(data []byte, appName string) bool {
	// Библиотека может упасть, но надеюсь, все падения можно поймать на таком типе проверок
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error: Recovered", r)
		}
	}()
	app, err := jsonparser.GetString(data, "app")
	if err != nil || app != appName {
		return false
	}
	return true
}

// fileExists проверяет наличие файла в path
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// stringInSlice проверяет, есть ли строка в слайсе
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// SetConfig добавляет в библиотеку новый конфиг или перезаписывает существующий
// Конфиги будут применяться в порядке возрастания имени ключа. Ex: 01-file, 02-add, 99-default
func SetConfig(config []byte, name string) error {
	if !CheckJson(config, appName) {
		return errors.New("Config: " + name + " - json check failed or field \"app\": \"" + appName + "\" not found")
	}
	// Добавляем конфиг
	configs.Set(name, config)
	// Список имён должен оставаться уникальным
	if !stringInSlice(name, configNames) {
		configNames = append(configNames, name)
	}
	// Сортируем список, чтобы применять конфиги в правильном порядке
	sort.Strings(configNames)
	// Сбрасываем кеш
	cache = cmap.New()
	stringArrayMutex.Lock()
	stringArrayCache = make(map[string][]string)
	stringArrayMutex.Unlock()
	return nil
}

// SetConfigFromFile добавляет в библиотеку новый или перезаписывает существующий конфиг из файла
func SetConfigFromFile(fileName, name string) error {
	if !fileExists(fileName) {
		return errors.New("Config file: " + fileName + " - doesn't exist")
	}
	config, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	return SetConfig(config, name)
}

// Init должен быть первым вызовом библиоткеи в приложении,
// applicationName должен соответствовать полю "app": appName в используемых конфигах
func Init(applicationName string) {
	cache = cmap.New()
	configs = cmap.New()
	stringArrayCache = make(map[string][]string)
	appName = applicationName
}
