@echo off
REM Test Swagger endpoint

echo Testing Swagger endpoint...
echo.

REM Generate docs first
echo Generating Swagger documentation...
swag init -g main.go -o ./docs
if %errorlevel% neq 0 (
    echo Failed to generate Swagger docs
    exit /b 1
)

echo.
echo Starting server in background...
start /B go run main.go

echo Waiting for server to start...
timeout /t 5 /nobreak >nul

echo.
echo Testing Swagger endpoints...
echo.

REM Test the JSON endpoint
echo Testing /swagger/doc.json...
powershell -Command "try { $response = Invoke-WebRequest -Uri 'http://localhost:8080/swagger/doc.json' -UseBasicParsing; Write-Host 'SUCCESS: JSON endpoint working'; Write-Host $response.StatusCode } catch { Write-Host 'ERROR: JSON endpoint failed'; Write-Host $_.Exception.Message }"

echo.
echo Testing /swagger/index.html...
powershell -Command "try { $response = Invoke-WebRequest -Uri 'http://localhost:8080/swagger/index.html' -UseBasicParsing; Write-Host 'SUCCESS: HTML endpoint working'; Write-Host $response.StatusCode } catch { Write-Host 'ERROR: HTML endpoint failed'; Write-Host $_.Exception.Message }"

echo.
echo Opening Swagger UI in browser...
start http://localhost:8080/swagger/index.html

echo.
echo Server is running. Press any key to stop...
pause >nul

REM Kill the background process
taskkill /F /IM go.exe >nul 2>&1
