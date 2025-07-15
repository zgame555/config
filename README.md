# Config Package

Go package สำหรับจัดการ configuration และโหลดไฟล์ config หลายรูปแบบอย่างง่ายดาย

## ฟีเจอร์

- **รองรับหลายรูปแบบ**: .env, .json, .yml, .yaml
- **โหลดไฟล์ config อัตโนมัติ** ตามนามสกุลไฟล์
- **รองรับ default values** สำหรับทุก data type
- **Type-safe methods** สำหรับ string, int, และ boolean
- **สามารถใช้งานแบบ instance-based หรือ global functions**
- **Nested configuration support** สำหรับ JSON/YAML (แปลงเป็น dot notation)
- **Array support** สำหรับ JSON/YAML (แปลงเป็น comma-separated string)
- **จัดการ comments และ empty lines** ในไฟล์ .env
- **รองรับ quoted values** ในไฟล์ .env
- **Hot reload** สำหรับการเปลี่ยนแปลง config
- **Format detection** อัตโนมัติตามนามสกุลไฟล์

## การติดตั้ง

```bash
go get github.com/zgame555/config
```

## การใช้งาน

### Instance-based Usage

```go
package main

import (
    "fmt"
    "github.com/zgame555/config"
)

func main() {
    // โหลดไฟล์ .env (default)
    env := config.New()

    // หรือระบุไฟล์ config เอง
    jsonConfig := config.New("config.json")
    yamlConfig := config.New("config.yaml")
    ymlConfig := config.New("app.yml")

    // อ่านค่า configuration
    appName := env.Str("APP_NAME")
    port := env.Int("PORT", 8080)        // default: 8080
    debug := env.Bool("DEBUG", false)    // default: false

    fmt.Printf("App: %s, Port: %d, Debug: %t\n", appName, port, debug)
}
```

### Global Functions

```go
package main

import (
    "fmt"
    "github.com/zgame555/config"
)

func main() {
    // โหลดไฟล์ config (รองรับทุก format)
    config.LoadConfigFile("config.json")
    // หรือ
    config.LoadConfigFile("config.yaml")
    // หรือ
    config.LoadConfigFile() // default: .env

    // อ่านค่า configuration
    dbHost := config.Str("DATABASE_HOST", "localhost")
    dbPort := config.Int("DATABASE_PORT", 5432)
    debug := config.Bool("DEBUG", false)

    fmt.Printf("DB: %s:%d, Debug: %t\n", dbHost, dbPort, debug)
}
```

## รูปแบบไฟล์ที่รองรับ

### 1. ไฟล์ .env

```env
# Application Configuration
APP_NAME=My Go App
APP_PORT=8080
APP_DEBUG=true

# Database Configuration
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME="myapp_db"
DATABASE_USER='postgres'
```

### 2. ไฟล์ JSON

```json
{
  "app": {
    "name": "My Go App",
    "port": 8080,
    "debug": true
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "name": "myapp_db"
  },
  "features": ["auth", "logging", "metrics"]
}
```

**การเข้าถึง:**

```go
config.Str("APP_NAME")           // "My Go App"
config.Int("APP_PORT")           // 8080
config.Bool("APP_DEBUG")         // true
config.Str("DATABASE_HOST")      // "localhost"
config.Str("FEATURES")           // "auth,logging,metrics"
```

### 3. ไฟล์ YAML/YML

```yaml
app:
  name: "My Go App"
  port: 8080
  debug: true

database:
  host: localhost
  port: 5432
  name: myapp_db

features:
  - auth
  - logging
  - metrics
```

**การเข้าถึง:**

```go
config.Str("APP_NAME")           // "My Go App"
config.Int("APP_PORT")           // 8080
config.Bool("APP_DEBUG")         // true
config.Str("DATABASE_HOST")      // "localhost"
config.Str("FEATURES")           // "auth,logging,metrics"
```

## API Reference

### Instance Methods

#### `New(configFile ...string) *Config`

สร้าง config instance ใหม่ และโหลดไฟล์ config อัตโนมัติ

- รองรับ .env, .json, .yml, .yaml
- Default: ".env"

#### `Load() error`

โหลดไฟล์ config (จะไม่โหลดซ้ำถ้าโหลดแล้ว)

#### `MustLoad()`

โหลดไฟล์ config และ panic ถ้าเกิดข้อผิดพลาด

#### `Str(key string, defaultValue ...string) string`

อ่านค่า string จาก environment variable

#### `Int(key string, defaultValue ...int) int`

อ่านค่า integer จาก environment variable

#### `Bool(key string, defaultValue ...bool) bool`

อ่านค่า boolean จาก environment variable
รองรับค่า: `true`, `1`, `yes`, `on` (true) และ `false`, `0`, `no`, `off` (false)

#### `All() map[string]string`

คืนค่า environment variables ทั้งหมดเป็น map

#### `Reload() error`

โหลดไฟล์ config ใหม่ (hot reload)

#### `SetFile(configFile string) error`

เปลี่ยนไฟล์ config และโหลดใหม่

### Global Functions

#### `LoadConfigFile(filePath ...string) error`

โหลดไฟล์ config (รองรับทุก format)

#### `MustLoadConfigFile(filePath ...string)`

โหลดไฟล์ config และ panic ถ้าเกิดข้อผิดพลาด

#### `LoadEnvFile(filePath ...string) error`

โหลดไฟล์ .env (backward compatibility)

#### `MustLoadEnvFile(filePath ...string)`

โหลดไฟล์ .env และ panic ถ้าเกิดข้อผิดพลาด

#### `Str(key string, defaultValue ...string) string`

อ่านค่า string จาก environment variable

#### `Int(key string, defaultValue ...int) int`

อ่านค่า integer จาก environment variable

#### `Bool(key string, defaultValue ...bool) bool`

อ่านค่า boolean จาก environment variable

#### `All() map[string]string`

คืนค่า environment variables ทั้งหมดเป็น map

## การทำงานกับ Nested Configuration

สำหรับไฟล์ JSON และ YAML nested objects จะถูกแปลงเป็น environment variables โดยใช้ dot notation และแปลงเป็น uppercase พร้อม underscore:

```yaml
database:
  connection:
    host: localhost
    port: 5432
    ssl:
      enabled: true
```

จะกลายเป็น:

- `DATABASE_CONNECTION_HOST` = "localhost"
- `DATABASE_CONNECTION_PORT` = "5432"
- `DATABASE_CONNECTION_SSL_ENABLED` = "true"

## การทำงานกับ Arrays

Arrays ใน JSON/YAML จะถูกแปลงเป็น comma-separated string:

```json
{
  "allowed_origins": ["http://localhost:3000", "https://myapp.com"],
  "features": ["auth", "logging", "metrics"]
}
```

การเข้าถึง:

```go
origins := config.Str("ALLOWED_ORIGINS")  // "http://localhost:3000,https://myapp.com"
features := config.Str("FEATURES")        // "auth,logging,metrics"

// แปลงกลับเป็น slice
originsList := strings.Split(origins, ",")
```

## ตัวอย่างการใช้งาน

ดูตัวอย่างการใช้งานใน [examples/usage.go](examples/usage.go)

### รันตัวอย่าง

```bash
cd examples
go run usage.go
```

## ข้อดี

- **ง่ายต่อการใช้งาน**: API ที่เรียบง่ายและใช้งานง่าย
- **รองรับหลายรูปแบบ**: .env, JSON, YAML ในแพ็คเกจเดียว
- **Type Safety**: Methods ที่ type-safe สำหรับ data types ต่างๆ
- **Flexible**: รองรับทั้ง instance-based และ global functions
- **Default Values**: รองรับ default values สำหรับทุก method
- **Hot Reload**: รองรับการโหลดใหม่โดยไม่ต้อง restart
- **Nested Support**: รองรับ nested configuration สำหรับ JSON/YAML
- **Error Handling**: จัดการ error อย่างเหมาะสม
- **Performance**: ไม่โหลดไฟล์ config ซ้ำ

## ข้อควรรู้

- ไฟล์ config ที่ไม่มีจะไม่ทำให้เกิด error
- Environment variables ที่มีอยู่แล้วจะถูก override เมื่อ reload
- รองรับ comments (บรรทัดที่ขึ้นต้นด้วย #) ในไฟล์ .env
- รองรับ quoted values (single และ double quotes) ในไฟล์ .env
- Empty lines จะถูกข้าม
- JSON/YAML nested objects จะถูกแปลงเป็น uppercase environment variables พร้อม underscore
- Arrays จะถูกแปลงเป็น comma-separated strings

## Dependencies

- [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3) สำหรับ YAML parsing

## License

MIT License
