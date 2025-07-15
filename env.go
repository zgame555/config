package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	envFile string
	loaded  bool
}

// New creates a new Env instance with optional .env file path
// If no file path is provided, it defaults to ".env"
func New(envFile ...string) *Config {
	file := ".env"
	if len(envFile) > 0 {
		file = envFile[0]
	}

	env := &Config{
		envFile: file,
		loaded:  false,
	}

	// Auto-load the .env file
	env.Load()

	return env
}

// Load loads the .env file into environment variables
func (e *Config) Load() error {
	if e.loaded {
		return nil // Already loaded
	}

	err := e.loadEnvFile(e.envFile)
	if err == nil {
		e.loaded = true
	}
	return err
}

// MustLoad loads the .env file and panics if there's an error
func (e *Config) MustLoad() {
	if err := e.Load(); err != nil {
		panic(fmt.Sprintf("failed to load env file: %v", err))
	}
}

// Str retrieves a string environment variable with optional default value
func (e *Config) Str(key string, defaultValue ...string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// Int retrieves an integer environment variable with optional default value
func (e *Config) Int(key string, defaultValue ...int) int {
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
func (e *Config) Bool(key string, defaultValue ...bool) bool {
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
func (e *Config) All() map[string]string {
	return All()
}

// Reload reloads the .env file
func (e *Config) Reload() error {
	e.loaded = false
	return e.Load()
}

// SetFile changes the .env file path and reloads
func (e *Config) SetFile(envFile string) error {
	e.envFile = envFile
	e.loaded = false
	return e.Load()
}

// loadEnvFile is the internal method to load .env file
func (e *Config) loadEnvFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		// If file doesn't exist, ignore silently
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open env file %s: %w", filePath, err)
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

		// Set environment variable only if it's not already set
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

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

		// Set environment variable only if it's not already set
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

func MustLoadEnvFile(filePath ...string) {
	if err := LoadEnvFile(filePath...); err != nil {
		panic(fmt.Sprintf("failed to load env file: %v", err))
	}
}

func Str(key string, defaultValue ...string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

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
