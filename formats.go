package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ConfigFormat represents the configuration file format
type ConfigFormat int

const (
	FormatEnv ConfigFormat = iota
	FormatJSON
	FormatYAML
)

// detectFormat detects the configuration file format based on file extension
func detectFormat(filePath string) ConfigFormat {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".json":
		return FormatJSON
	case ".yml", ".yaml":
		return FormatYAML
	default:
		return FormatEnv
	}
}

// loadConfigFile loads configuration from various file formats
func loadConfigFile(filePath string) (map[string]interface{}, error) {
	format := detectFormat(filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		// If file doesn't exist, return empty config
		if os.IsNotExist(err) {
			return make(map[string]interface{}), nil
		}
		return nil, fmt.Errorf("failed to read config file %s: %w", filePath, err)
	}

	switch format {
	case FormatJSON:
		return loadJSONConfig(data)
	case FormatYAML:
		return loadYAMLConfig(data)
	case FormatEnv:
		return loadEnvConfig(data)
	default:
		return nil, fmt.Errorf("unsupported config format for file: %s", filePath)
	}
}

// loadJSONConfig loads configuration from JSON data
func loadJSONConfig(data []byte) (map[string]interface{}, error) {
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse JSON config: %w", err)
	}
	return flattenConfig(config, ""), nil
}

// loadYAMLConfig loads configuration from YAML data
func loadYAMLConfig(data []byte) (map[string]interface{}, error) {
	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}
	return flattenConfig(config, ""), nil
}

// loadEnvConfig loads configuration from ENV data
func loadEnvConfig(data []byte) (map[string]interface{}, error) {
	config := make(map[string]interface{})
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		config[key] = value
	}

	return config, nil
}

// flattenConfig flattens nested configuration into dot notation
func flattenConfig(config map[string]interface{}, prefix string) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range config {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			// Recursively flatten nested maps
			nested := flattenConfig(v, fullKey)
			for nestedKey, nestedValue := range nested {
				result[nestedKey] = nestedValue
			}
		case []interface{}:
			// Convert arrays to comma-separated strings
			var strValues []string
			for _, item := range v {
				strValues = append(strValues, fmt.Sprintf("%v", item))
			}
			result[fullKey] = strings.Join(strValues, ",")
		default:
			result[fullKey] = fmt.Sprintf("%v", value)
		}
	}

	return result
}

// setEnvironmentVariables sets environment variables from config map
func setEnvironmentVariables(config map[string]interface{}) {
	for key, value := range config {
		// Convert key to uppercase for environment variables
		envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))

		// Always set the environment variable (allow override for reload)
		os.Setenv(envKey, fmt.Sprintf("%v", value))
	}
}

// clearEnvironmentVariables clears environment variables that were set from config
func clearEnvironmentVariables(config map[string]interface{}) {
	for key := range config {
		// Convert key to uppercase for environment variables
		envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
		os.Unsetenv(envKey)
	}
}
