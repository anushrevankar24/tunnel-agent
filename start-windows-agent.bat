@echo off
REM Windows Agent Startup Script
REM This script starts the tunnel agent with proper configuration

echo Starting Tunnel Agent...

REM Load environment variables from config file
if exist "config.env" (
    echo Loading configuration from config.env...
    for /f "usebackq tokens=1,2 delims==" %%a in ("config.env") do (
        set %%a=%%b
    )
) else (
    echo Warning: config.env not found. Using default values.
)

REM Start the agent
echo Connecting to server: %AGENT_SERVER_URL%
echo Local service: %AGENT_LOCAL_URL%
echo Local API port: %AGENT_LOCAL_PORT%

agent.exe -server %AGENT_SERVER_URL% -local %AGENT_LOCAL_URL% -port %AGENT_LOCAL_PORT% -timeout %AGENT_TIMEOUT%

pause

