package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	configFile   string
	loaded       bool
	format       ConfigFormat
	loadedConfig map[string]interface{} // Keep track of loaded config for reload
}

// New creates a new Config instance with optional config file path
// If no file path is provided, it defaults to ".env"
// Supports .env, .json, .yml, .yaml formats
func New(configFile ...string) *Config {
	file := ".env"
	if len(configFile) > 0 {
		file = configFile[0]
	}

	config := &Config{
		configFile:   file,
		loaded:       false,
		format:       detectFormat(file),
		loadedConfig: make(map[string]interface{}),
	}

	// Auto-load the config file
	config.Load()

	return config
}

// Load loads the config file into environment variables
func (c *Config) Load() error {
	if c.loaded {
		return nil // Already loaded
	}

	var err error
	switch c.format {
	case FormatEnv:
		err = c.loadEnvFile(c.configFile)
	case FormatJSON, FormatYAML:
		err = c.loadStructuredFile(c.configFile)
	default:
		err = fmt.Errorf("unsupported config format for file: %s", c.configFile)
	}

	if err == nil {
		c.loaded = true
	}
	return err
}

// MustLoad loads the config file and panics if there's an error
func (c *Config) MustLoad() {
	if err := c.Load(); err != nil {
		panic(fmt.Sprintf("failed to load config file: %v", err))
	}
}

// Str retrieves a string environment variable with optional default value
func (c *Config) Str(key string, defaultValue ...string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// Int retrieves an integer environment variable with optional default value
func (c *Config) Int(key string, defaultValue ...int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

// Bool retrieves a boolean environment variable with optional default value
func (c *Config) Bool(key string, defaultValue ...bool) bool {
	if value := os.Getenv(key); value != "" {
		lowerValue := strings.ToLower(strings.TrimSpace(value))
		switch lowerValue {
		case "true", "1", "yes", "on":
			return true
		case "false", "0", "no", "off":
			return false
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return false
}

// All returns all environment variables as a map
func (c *Config) All() map[string]string {
	return All()
}

// Reload reloads the config file
func (c *Config) Reload() error {
	// Clear previously loaded config
	if len(c.loadedConfig) > 0 {
		clearEnvironmentVariables(c.loadedConfig)
		c.loadedConfig = make(map[string]interface{})
	}

	c.loaded = false
	return c.Load()
}

// SetFile changes the config file path and reloads
func (c *Config) SetFile(configFile string) error {
	// Clear previously loaded config
	if len(c.loadedConfig) > 0 {
		clearEnvironmentVariables(c.loadedConfig)
		c.loadedConfig = make(map[string]interface{})
	}

	c.configFile = configFile
	c.format = detectFormat(configFile)
	c.loaded = false
	return c.Load()
}

// loadStructuredFile loads JSON/YAML config files
func (c *Config) loadStructuredFile(filePath string) error {
	config, err := loadConfigFile(filePath)
	if err != nil {
		return err
	}

	// Store loaded config for reload functionality
	c.loadedConfig = config

	// Set environment variables from config
	setEnvironmentVariables(config)
	return nil
}

// loadEnvFile is the internal method to load .env file (backward compatibility)
func (c *Config) loadEnvFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		// If file doesn't exist, ignore silently
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open env file %s: %w", filePath, err)
	}
	defer file.Close()

	config := make(map[string]interface{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

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

		// Store in config for reload functionality
		config[key] = value

		// Set environment variable
		os.Setenv(key, value)
	}

	// Store loaded config for reload functionality
	c.loadedConfig = config
	return scanner.Err()
}

// Global functions for backward compatibility

// LoadConfigFile loads configuration from various file formats (.env, .json, .yml, .yaml)
func LoadConfigFile(filePath ...string) error {
	configFile := ".env"
	if len(filePath) > 0 {
		configFile = filePath[0]
	}

	format := detectFormat(configFile)

	switch format {
	case FormatEnv:
		return LoadEnvFile(configFile)
	case FormatJSON, FormatYAML:
		config, err := loadConfigFile(configFile)
		if err != nil {
			return err
		}
		setEnvironmentVariables(config)
		return nil
	default:
		return fmt.Errorf("unsupported config format for file: %s", configFile)
	}
}

// MustLoadConfigFile loads configuration file and panics if there's an error
func MustLoadConfigFile(filePath ...string) {
	if err := LoadConfigFile(filePath...); err != nil {
		panic(fmt.Sprintf("failed to load config file: %v", err))
	}
}

// LoadEnvFile loads .env file (backward compatibility)
func LoadEnvFile(filePath ...string) error {
	envFile := ".env"
	if len(filePath) > 0 {
		envFile = filePath[0]
	}

	file, err := os.Open(envFile)
	if err != nil {
		// If file doesn't exist, ignore silently
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open env file %s: %w", envFile, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

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

		// Set environment variable
		os.Setenv(key, value)
	}

	return scanner.Err()
}

// MustLoadEnvFile loads .env file and panics if there's an error (backward compatibility)
func MustLoadEnvFile(filePath ...string) {
	if err := LoadEnvFile(filePath...); err != nil {
		panic(fmt.Sprintf("failed to load env file: %v", err))
	}
}

// Str retrieves a string environment variable with optional default value
func Str(key string, defaultValue ...string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// Int retrieves an integer environment variable with optional default value
func Int(key string, defaultValue ...int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

// Bool retrieves a boolean environment variable with optional default value
func Bool(key string, defaultValue ...bool) bool {
	if value := os.Getenv(key); value != "" {
		lowerValue := strings.ToLower(strings.TrimSpace(value))
		switch lowerValue {
		case "true", "1", "yes", "on":
			return true
		case "false", "0", "no", "off":
			return false
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return false
}

// All returns all environment variables as a map
func All() map[string]string {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) >= 2 {
			envs[parts[0]] = parts[1]
		} else if len(parts) == 1 {
			envs[parts[0]] = ""
		}
	}
	return envs
}
