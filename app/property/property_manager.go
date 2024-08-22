package property

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"strings"
)

type Manager struct {
	properties      Properties
	overrideHelper  OverrideHelper
	keysInOrder     []string
	DefaultHTTPPort int
}

type Properties map[string]string

func NewManager() *Manager {
	return &Manager{properties: make(Properties)}
}

func (m *Manager) EnableLocalPropertyOverride(appName string) {
	m.overrideHelper.App = appName
}

func (m *Manager) LoadPropertiesByPath(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	m.LoadProperties(file)
}

func (m *Manager) LoadProperties(file fs.File) {
	scanner := bufio.NewScanner(file)
	var key, value string
	var isMultiLine bool

	addProperty := func() {
		if key != "" {
			m.add(strings.TrimSpace(key), strings.TrimSpace(value))
			key, value = "", ""
		}
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if isMultiLine {
			if strings.HasSuffix(line, "\\") {
				value += line[:len(line)-1]
			} else {
				value += line
				isMultiLine = false
				addProperty()
			}
			continue
		}

		if line == "" || line[0] == '#' || strings.HasPrefix(line, "//") {
			continue
		}

		index := strings.IndexByte(line, '=')
		if index == -1 {
			continue
		}

		key = line[:index]
		value = line[index+1:]

		if strings.HasSuffix(value, "\\") {
			isMultiLine = true
			value = value[:len(value)-1]
		} else {
			addProperty()
		}
	}

	addProperty()

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading properties:", err)
	}
}

func (m *Manager) add(key, value string) {
	m.properties[key] = value
	m.keysInOrder = append(m.keysInOrder, key)
}

func (m *Manager) GetKeysInOrder() []string {
	return m.keysInOrder
}

func (m *Manager) Get(key string, required ...bool) string {
	value, ok := m.properties[key]
	if required != nil && len(required) > 0 && required[0] {
		if !ok {
			log.Panic("required property not found! key:" + key)
		}
	}

	if !ok {
		slog.Debug(fmt.Sprintf("property not found! key=%s", key))
		return ""
	}

	overrideValue := m.overrideHelper.Get(key)
	if overrideValue != "" {
		return overrideValue
	}

	envVarName := EnvVarName(key)
	envVarValue := os.Getenv(envVarName)
	if envVarValue != "" {
		slog.Warn(fmt.Sprintf("found overridden property by env var %s, key=%s, value=%s", envVarName, key, MaskValue(key, envVarValue)))
		return envVarValue
	}

	return value
}

func (m *Manager) Keys() []string {
	var keys = make([]string, len(m.keysInOrder))
	copy(keys, m.keysInOrder)
	return keys
}

func MaskValue(key, value string) string {
	lowerCaseKey := strings.ToLower(key)
	if strings.Contains(lowerCaseKey, "password") || strings.Contains(lowerCaseKey, "secret") {
		return "******"
	}
	return value
}
