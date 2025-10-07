@echo off
REM Run the Go project and automatically open Swagger UI in browser

echo Starting Go project with Swagger UI...
echo.

REM Generate Swagger docs first
echo Generating Swagger documentation...
swag init -g main.go -o ./docs
if %errorlevel% neq 0 (
    echo Failed to generate Swagger docs
    exit /b 1
)

echo.
echo Starting server...
echo 🌐 Swagger UI will be available at: http://localhost:8080/swagger/index.html
echo 🚀 Browser will open automatically in 3 seconds...
echo.

REM Wait 3 seconds then open browser
powershell -Command "Start-Sleep 3; Start-Process 'http://localhost:8080/swagger/index.html'"

REM Start the Go application
go run main.go
