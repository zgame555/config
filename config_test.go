package config

import (
	"os"
	"testing"
)

// Helper function to create test files
func createTestFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

// Helper function to cleanup test files
func cleanupTestFile(filename string) {
	os.Remove(filename)
}

func TestEnvFormat(t *testing.T) {
	// Create test .env file
	envContent := `# Test environment variables
DB_HOST=localhost
DB_PORT=5432
DEBUG=true
API_KEY="secret-key"
FEATURE_FLAG='enabled'
`
	err := createTestFile("test.env", envContent)
	if err != nil {
		t.Fatalf("Failed to create test env file: %v", err)
	}
	defer cleanupTestFile("test.env")

	// Test loading
	config := New("test.env")

	// Test string values
	if config.Str("DB_HOST") != "localhost" {
		t.Errorf("Expected DB_HOST=localhost, got %s", config.Str("DB_HOST"))
	}

	// Test int values
	if config.Int("DB_PORT") != 5432 {
		t.Errorf("Expected DB_PORT=5432, got %d", config.Int("DB_PORT"))
	}

	// Test bool values
	if !config.Bool("DEBUG") {
		t.Errorf("Expected DEBUG=true, got %v", config.Bool("DEBUG"))
	}

	// Test quoted values
	if config.Str("API_KEY") != "secret-key" {
		t.Errorf("Expected API_KEY=secret-key, got %s", config.Str("API_KEY"))
	}

	if config.Str("FEATURE_FLAG") != "enabled" {
		t.Errorf("Expected FEATURE_FLAG=enabled, got %s", config.Str("FEATURE_FLAG"))
	}
}

func TestJSONFormat(t *testing.T) {
	// Create test JSON file
	jsonContent := `{
  "database": {
    "host": "localhost",
    "port": 5432,
    "name": "testdb"
  },
  "app": {
    "debug": true,
    "features": ["auth", "logging", "metrics"]
  },
  "api_key": "json-secret-key"
}`
	err := createTestFile("test.json", jsonContent)
	if err != nil {
		t.Fatalf("Failed to create test json file: %v", err)
	}
	defer cleanupTestFile("test.json")

	// Test loading
	config := New("test.json")

	// Test nested values (flattened with dots)
	if config.Str("DATABASE_HOST") != "localhost" {
		t.Errorf("Expected DATABASE_HOST=localhost, got %s", config.Str("DATABASE_HOST"))
	}

	if config.Int("DATABASE_PORT") != 5432 {
		t.Errorf("Expected DATABASE_PORT=5432, got %d", config.Int("DATABASE_PORT"))
	}

	if config.Str("DATABASE_NAME") != "testdb" {
		t.Errorf("Expected DATABASE_NAME=testdb, got %s", config.Str("DATABASE_NAME"))
	}

	if !config.Bool("APP_DEBUG") {
		t.Errorf("Expected APP_DEBUG=true, got %v", config.Bool("APP_DEBUG"))
	}

	// Test array values (converted to comma-separated string)
	features := config.Str("APP_FEATURES")
	expected := "auth,logging,metrics"
	if features != expected {
		t.Errorf("Expected APP_FEATURES=%s, got %s", expected, features)
	}

	if config.Str("API_KEY") != "json-secret-key" {
		t.Errorf("Expected API_KEY=json-secret-key, got %s", config.Str("API_KEY"))
	}
}

func TestYAMLFormat(t *testing.T) {
	// Create test YAML file
	yamlContent := `database:
  host: localhost
  port: 5432
  name: testdb
  ssl: false

app:
  debug: true
  name: "Test App"
  features:
    - auth
    - logging
    - metrics

api_key: yaml-secret-key
`
	err := createTestFile("test.yaml", yamlContent)
	if err != nil {
		t.Fatalf("Failed to create test yaml file: %v", err)
	}
	defer cleanupTestFile("test.yaml")

	// Test loading
	config := New("test.yaml")

	// Test nested values
	if config.Str("DATABASE_HOST") != "localhost" {
		t.Errorf("Expected DATABASE_HOST=localhost, got %s", config.Str("DATABASE_HOST"))
	}

	if config.Int("DATABASE_PORT") != 5432 {
		t.Errorf("Expected DATABASE_PORT=5432, got %d", config.Int("DATABASE_PORT"))
	}

	if config.Bool("DATABASE_SSL") {
		t.Errorf("Expected DATABASE_SSL=false, got %v", config.Bool("DATABASE_SSL"))
	}

	if !config.Bool("APP_DEBUG") {
		t.Errorf("Expected APP_DEBUG=true, got %v", config.Bool("APP_DEBUG"))
	}

	if config.Str("APP_NAME") != "Test App" {
		t.Errorf("Expected APP_NAME=Test App, got %s", config.Str("APP_NAME"))
	}

	// Test array values
	features := config.Str("APP_FEATURES")
	expected := "auth,logging,metrics"
	if features != expected {
		t.Errorf("Expected APP_FEATURES=%s, got %s", expected, features)
	}

	if config.Str("API_KEY") != "yaml-secret-key" {
		t.Errorf("Expected API_KEY=yaml-secret-key, got %s", config.Str("API_KEY"))
	}
}

func TestYMLFormat(t *testing.T) {
	// Create test .yml file
	ymlContent := `server:
  host: 0.0.0.0
  port: 8080
  timeout: 30

logging:
  level: info
  format: json
  enabled: true
`
	err := createTestFile("test.yml", ymlContent)
	if err != nil {
		t.Fatalf("Failed to create test yml file: %v", err)
	}
	defer cleanupTestFile("test.yml")

	// Test loading
	config := New("test.yml")

	// Test nested values
	if config.Str("SERVER_HOST") != "0.0.0.0" {
		t.Errorf("Expected SERVER_HOST=0.0.0.0, got %s", config.Str("SERVER_HOST"))
	}

	if config.Int("SERVER_PORT") != 8080 {
		t.Errorf("Expected SERVER_PORT=8080, got %d", config.Int("SERVER_PORT"))
	}

	if config.Int("SERVER_TIMEOUT") != 30 {
		t.Errorf("Expected SERVER_TIMEOUT=30, got %d", config.Int("SERVER_TIMEOUT"))
	}

	if config.Str("LOGGING_LEVEL") != "info" {
		t.Errorf("Expected LOGGING_LEVEL=info, got %s", config.Str("LOGGING_LEVEL"))
	}

	if !config.Bool("LOGGING_ENABLED") {
		t.Errorf("Expected LOGGING_ENABLED=true, got %v", config.Bool("LOGGING_ENABLED"))
	}
}

func TestDefaultValues(t *testing.T) {
	// Test with non-existent file
	config := New("nonexistent.env")

	// Test default values
	if config.Str("NONEXISTENT_KEY", "default") != "default" {
		t.Errorf("Expected default value for string")
	}

	if config.Int("NONEXISTENT_KEY", 42) != 42 {
		t.Errorf("Expected default value for int")
	}

	if config.Bool("NONEXISTENT_KEY", true) != true {
		t.Errorf("Expected default value for bool")
	}
}

func TestGlobalFunctions(t *testing.T) {
	// Create test JSON file
	jsonContent := `{
  "global_test": "success",
  "global_port": 3000,
  "global_debug": false
}`
	err := createTestFile("global_test.json", jsonContent)
	if err != nil {
		t.Fatalf("Failed to create test json file: %v", err)
	}
	defer cleanupTestFile("global_test.json")

	// Test global LoadConfigFile function
	err = LoadConfigFile("global_test.json")
	if err != nil {
		t.Fatalf("Failed to load config file: %v", err)
	}

	// Test global accessor functions
	if Str("GLOBAL_TEST") != "success" {
		t.Errorf("Expected GLOBAL_TEST=success, got %s", Str("GLOBAL_TEST"))
	}

	if Int("GLOBAL_PORT") != 3000 {
		t.Errorf("Expected GLOBAL_PORT=3000, got %d", Int("GLOBAL_PORT"))
	}

	if Bool("GLOBAL_DEBUG") {
		t.Errorf("Expected GLOBAL_DEBUG=false, got %v", Bool("GLOBAL_DEBUG"))
	}
}

func TestFormatDetection(t *testing.T) {
	tests := []struct {
		filename string
		expected ConfigFormat
	}{
		{".env", FormatEnv},
		{"config.env", FormatEnv},
		{"app.json", FormatJSON},
		{"config.yml", FormatYAML},
		{"app.yaml", FormatYAML},
		{"unknown.txt", FormatEnv}, // Default to env
	}

	for _, test := range tests {
		format := detectFormat(test.filename)
		if format != test.expected {
			t.Errorf("Expected format %d for %s, got %d", test.expected, test.filename, format)
		}
	}
}

func TestReload(t *testing.T) {
	// Create initial config file
	initialContent := `TEST_VALUE=initial`
	err := createTestFile("reload_test.env", initialContent)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer cleanupTestFile("reload_test.env")

	config := New("reload_test.env")

	// Check initial value
	if config.Str("TEST_VALUE") != "initial" {
		t.Errorf("Expected TEST_VALUE=initial, got %s", config.Str("TEST_VALUE"))
	}

	// Update file content
	updatedContent := `TEST_VALUE=updated`
	err = createTestFile("reload_test.env", updatedContent)
	if err != nil {
		t.Fatalf("Failed to update test file: %v", err)
	}

	// Reload and check updated value
	err = config.Reload()
	if err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	if config.Str("TEST_VALUE") != "updated" {
		t.Errorf("Expected TEST_VALUE=updated after reload, got %s", config.Str("TEST_VALUE"))
	}
}
