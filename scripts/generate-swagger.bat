@echo off
REM Swagger generation script for Go project (Windows)
REM This script generates Swagger documentation from Go annotations

echo Generating Swagger documentation...

REM Check if swag is installed
swag version >nul 2>&1
if %errorlevel% neq 0 (
    echo swag command not found. Installing swag...
    go install github.com/swaggo/swag/cmd/swag@latest
)

REM Generate docs
echo Running swag init...
swag init -g main.go -o ./docs

REM Check if generation was successful
if %errorlevel% equ 0 (
    echo ✅ Swagger documentation generated successfully!
    echo 📁 Documentation files created in ./docs/
    echo 🌐 Access Swagger UI at: http://localhost:8080/swagger/index.html
) else (
    echo ❌ Failed to generate Swagger documentation
    exit /b 1
)

echo 🎉 Swagger setup complete!
