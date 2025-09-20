@echo off
echo 📋 Checking System Requirements...
echo ================================

echo.
echo 🐳 Checking Docker...
docker --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker is not installed or not in PATH
    echo.
    echo 📥 Please install Docker Desktop for Windows:
    echo    1. Visit: https://www.docker.com/products/docker-desktop/
    echo    2. Download Docker Desktop for Windows
    echo    3. Run the installer and follow the instructions
    echo    4. Restart your computer if prompted
    echo    5. Start Docker Desktop
    echo    6. Run this script again
    echo.
    goto :check_go
) else (
    docker --version
    echo ✅ Docker is installed

    echo.
    echo 🔄 Checking Docker service...
    docker info >nul 2>&1
    if errorlevel 1 (
        echo ❌ Docker service is not running
        echo 💡 Please start Docker Desktop and try again
        goto :check_go
    ) else (
        echo ✅ Docker service is running
    )
)

:check_go
echo.
echo 🏃 Checking Go...
go version >nul 2>&1
if errorlevel 1 (
    echo ❌ Go is not installed or not in PATH
    echo.
    echo 📥 Please install Go:
    echo    1. Visit: https://golang.org/dl/
    echo    2. Download Go 1.21+ for Windows
    echo    3. Run the installer
    echo    4. Restart your command prompt
    echo    5. Run this script again
    echo.
    goto :check_make
) else (
    go version
    echo ✅ Go is installed
)

:check_make
echo.
echo 🔨 Checking Make...
make --version >nul 2>&1
if errorlevel 1 (
    echo ⚠️  Make is not installed (optional for development)
    echo.
    echo 📥 To install Make (optional):
    echo    Option 1 - Install via Chocolatey:
    echo      1. Install Chocolatey: https://chocolatey.org/install
    echo      2. Run: choco install make
    echo.
    echo    Option 2 - Use batch files instead:
    echo      - Use start.bat instead of 'make start'
    echo      - Use dev.bat instead of 'make dev'
    echo      - Use stop.bat instead of 'make stop'
    echo.
) else (
    make --version | findstr "GNU Make"
    echo ✅ Make is installed
)

echo.
echo 📊 Summary:
echo =========
docker --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker: Not installed
    set "docker_ok=false"
) else (
    docker info >nul 2>&1
    if errorlevel 1 (
        echo ⚠️  Docker: Installed but not running
        set "docker_ok=false"
    ) else (
        echo ✅ Docker: Ready
        set "docker_ok=true"
    )
)

go version >nul 2>&1
if errorlevel 1 (
    echo ❌ Go: Not installed
) else (
    echo ✅ Go: Ready
)

if "%docker_ok%"=="true" (
    echo.
    echo 🎉 You can now start the services using:
    echo    start.bat    ^(recommended^)
    echo    OR
    echo    make start   ^(if Make is installed^)
) else (
    echo.
    echo 🔧 Please install/start Docker first, then run the services
)

echo.
pause