# Config Package

Go package สำหรับจัดการ environment variables และโหลดไฟล์ .env อย่างง่ายดาย

## ฟีเจอร์

- โหลดไฟล์ .env อัตโนมัติ
- รองรับ default values สำหรับทุก data type
- Type-safe methods สำหรับ string, int, และ boolean
- สามารถใช้งานแบบ instance-based หรือ global functions
- จัดการ comments และ empty lines ในไฟล์ .env
- รองรับ quoted values
- ไม่ override environment variables ที่มีอยู่แล้ว

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
    // สร้าง config instance (จะโหลด .env อัตโนมัติ)
    env := config.New()

    // หรือระบุไฟล์ .env เอง
    env := config.New("custom.env")

    // อ่านค่า environment variables
    dbHost := env.Str("DB_HOST", "localhost")
    dbPort := env.Int("DB_PORT", 5432)
    debug := env.Bool("DEBUG", false)

    fmt.Printf("DB Host: %s\n", dbHost)
    fmt.Printf("DB Port: %d\n", dbPort)
    fmt.Printf("Debug: %t\n", debug)
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
    // โหลดไฟล์ .env
    config.LoadEnvFile() // โหลด .env
    // หรือ
    config.LoadEnvFile("custom.env") // โหลดไฟล์ที่ระบุ

    // อ่านค่า environment variables
    dbHost := config.Str("DB_HOST", "localhost")
    dbPort := config.Int("DB_PORT", 5432)
    debug := config.Bool("DEBUG", false)

    fmt.Printf("DB Host: %s\n", dbHost)
    fmt.Printf("DB Port: %d\n", dbPort)
    fmt.Printf("Debug: %t\n", debug)
}
```

### ตัวอย่างไฟล์ .env

```env
# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME="myapp"
DB_USER='admin'

# Application settings
DEBUG=true
PORT=8080
API_KEY=your-secret-key

# Boolean values
ENABLE_FEATURE=1
DISABLE_CACHE=false
```

## API Reference

### Instance Methods

#### `New(envFile ...string) *Config`

สร้าง config instance ใหม่ และโหลดไฟล์ .env อัตโนมัติ

#### `Load() error`

โหลดไฟล์ .env (จะไม่โหลดซ้ำถ้าโหลดแล้ว)

#### `MustLoad()`

โหลดไฟล์ .env และ panic ถ้าเกิดข้อผิดพลาด

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

โหลดไฟล์ .env ใหม่

#### `SetFile(envFile string) error`

เปลี่ยนไฟล์ .env และโหลดใหม่

### Global Functions

#### `LoadEnvFile(filePath ...string) error`

โหลดไฟล์ .env (default: ".env")

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

## ข้อดี

- **ง่ายต่อการใช้งาน**: API ที่เรียบง่ายและใช้งานง่าย
- **Type Safety**: Methods ที่ type-safe สำหรับ data types ต่างๆ
- **Flexible**: รองรับทั้ง instance-based และ global functions
- **Default Values**: รองรับ default values สำหรับทุก method
- **Error Handling**: จัดการ error อย่างเหมาะสม
- **Performance**: ไม่โหลดไฟล์ .env ซ้ำ

## ข้อควรรู้

- ไฟล์ .env ที่ไม่มีจะไม่ทำให้เกิด error
- Environment variables ที่มีอยู่แล้วจะไม่ถูก override
- รองรับ comments (บรรทัดที่ขึ้นต้นด้วย #)
- รองรับ quoted values (single และ double quotes)
- Empty lines จะถูกข้าม

## License

MIT License
