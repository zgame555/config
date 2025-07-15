# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2024-12-19

### Added

- **รองรับไฟล์ JSON**: สามารถโหลดไฟล์ .json และแปลง nested objects เป็น environment variables
- **รองรับไฟล์ YAML**: สามารถโหลดไฟล์ .yml และ .yaml
- **Nested Configuration Support**: รองรับ nested configuration โดยแปลงเป็น dot notation
- **Array Support**: รองรับ arrays โดยแปลงเป็น comma-separated strings
- **Format Auto-detection**: ตรวจจับรูปแบบไฟล์อัตโนมัติจากนามสกุล
- **Hot Reload**: รองรับการโหลดใหม่โดยไม่ต้อง restart application
- **LoadConfigFile()**: Global function สำหรับโหลดไฟล์ config รูปแบบต่างๆ
- **MustLoadConfigFile()**: Global function พร้อม panic on error
- **SetFile()**: Method สำหรับเปลี่ยนไฟล์ config
- **clearEnvironmentVariables()**: Internal function สำหรับล้าง environment variables
- **ไฟล์ตัวอย่าง**: เพิ่มตัวอย่างการใช้งานสำหรับทุก format
- **Comprehensive Tests**: เพิ่ม test cases สำหรับทุก format และฟีเจอร์

### Changed

- **Config struct**: เพิ่ม field `format` และ `loadedConfig` สำหรับ tracking
- **Environment Variable Handling**: เปลี่ยนจาก "set only if not exists" เป็น "always set" เพื่อรองรับ reload
- **Nested Key Conversion**: JSON/YAML nested keys จะถูกแปลงเป็น uppercase พร้อม underscore
- **README.md**: อัปเดตเอกสารให้ครอบคลุมฟีเจอร์ใหม่
- **Error Messages**: ปรับปรุงข้อความ error ให้ชัดเจนขึ้น

### Dependencies

- **เพิ่ม gopkg.in/yaml.v3**: สำหรับ YAML parsing

### Backward Compatibility

- **รองรับ API เดิม**: ทุก API เดิมยังใช้งานได้ปกติ
- **LoadEnvFile()**: ยังคงใช้งานได้เพื่อ backward compatibility
- **MustLoadEnvFile()**: ยังคงใช้งานได้เพื่อ backward compatibility

## [1.0.0] - 2024-12-18

### Added

- **รองรับไฟล์ .env**: โหลดและ parse ไฟล์ .env
- **Type-safe Methods**: `Str()`, `Int()`, `Bool()` พร้อม default values
- **Instance-based และ Global Functions**: รองรับทั้งสองรูปแบบการใช้งาน
- **Comment Support**: รองรับ comments ในไฟล์ .env
- **Quoted Values**: รองรับ single และ double quotes
- **Auto-load**: โหลดไฟล์ .env อัตโนมัติเมื่อสร้าง instance
- **Error Handling**: จัดการ error อย่างเหมาะสม
- **All()**: Method สำหรับดึง environment variables ทั้งหมด

### Initial Features

- Basic .env file parsing
- Environment variable management
- Default value support
- Comment and empty line handling
- Quoted value support
