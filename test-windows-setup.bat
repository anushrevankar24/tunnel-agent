@echo off
REM Windows Tunnel Agent Test Script
REM This script helps test the agent setup and connection

echo ========================================
echo    Windows Tunnel Agent Test Script
echo ========================================
echo.

REM Check if agent.exe exists
if not exist "agent.exe" (
    echo ❌ ERROR: agent.exe not found!
    echo Please ensure agent.exe is in the same folder as this script.
    pause
    exit /b 1
)

REM Check if config file exists
if not exist "config.env" (
    echo ❌ ERROR: config.env not found!
    echo Please ensure config.env is in the same folder as this script.
    pause
    exit /b 1
)

echo ✅ Found agent.exe
echo ✅ Found config.env
echo.

REM Load configuration
echo Loading configuration...
for /f "usebackq tokens=1,2 delims==" %%a in ("config.env") do (
    set %%a=%%b
)

echo Configuration loaded:
echo   Server URL: %AGENT_SERVER_URL%
echo   Local URL: %AGENT_LOCAL_URL%
echo   Local API Port: %AGENT_LOCAL_PORT%
echo.

REM Test server connectivity
echo Testing server connectivity...
curl -s --max-time 10 https://your-app-name.onrender.com/health >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ Server is reachable
) else (
    echo ❌ WARNING: Cannot reach server. Check your internet connection.
)
echo.

REM Check if local service is running
echo Testing local service...
curl -s --max-time 5 %AGENT_LOCAL_URL% >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ Local service is running on %AGENT_LOCAL_URL%
) else (
    echo ❌ WARNING: Local service not accessible on %AGENT_LOCAL_URL%
    echo Please start your local service before running the agent.
)
echo.

echo ========================================
echo           Setup Test Complete
echo ========================================
echo.
echo Next steps:
echo 1. Start your local service (if not already running)
echo 2. Run start-windows-agent.bat to start the agent
echo 3. Test the tunnel: https://your-app-name.onrender.com/tunnel/windows-test/
echo.
echo Press any key to continue...
pause >nul
