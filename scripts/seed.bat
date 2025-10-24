@echo off
setlocal enabledelayedexpansion

REM Inventory API Seeder Script for Windows
REM This script provides easy access to run seeders with common configurations

set DEFAULT_COUNT=20
set SCRIPT_DIR=%~dp0
set PROJECT_ROOT=%SCRIPT_DIR%..

REM Function to show help
if "%1"=="--help" goto :show_help
if "%1"=="-h" goto :show_help
if "%1"=="/?" goto :show_help

REM Parse arguments
set TYPE=all
set COUNT=%DEFAULT_COUNT%
set QUICK_SETUP=false
set DEMO_DATA=false
set LOAD_TEST=false

:parse_args
if "%1"=="" goto :validate_args

if "%1"=="-t" (
    set TYPE=%2
    shift
    shift
    goto :parse_args
)
if "%1"=="--type" (
    set TYPE=%2
    shift
    shift
    goto :parse_args
)
if "%1"=="-c" (
    set COUNT=%2
    shift
    shift
    goto :parse_args
)
if "%1"=="--count" (
    set COUNT=%2
    shift
    shift
    goto :parse_args
)
if "%1"=="--quick-setup" (
    set QUICK_SETUP=true
    shift
    goto :parse_args
)
if "%1"=="--demo-data" (
    set DEMO_DATA=true
    shift
    goto :parse_args
)
if "%1"=="--load-test" (
    set LOAD_TEST=true
    shift
    goto :parse_args
)

REM Handle --type=value format
echo %1 | findstr "^--type=" >nul
if !errorlevel! equ 0 (
    for /f "tokens=2 delims==" %%a in ("%1") do set TYPE=%%a
    shift
    goto :parse_args
)

REM Handle --count=value format
echo %1 | findstr "^--count=" >nul
if !errorlevel! equ 0 (
    for /f "tokens=2 delims==" %%a in ("%1") do set COUNT=%%a
    shift
    goto :parse_args
)

echo Unknown option: %1
goto :show_help

:validate_args
REM Check for special modes
if "%QUICK_SETUP%"=="true" goto :quick_setup
if "%DEMO_DATA%"=="true" goto :demo_data
if "%LOAD_TEST%"=="true" goto :load_test

REM Validate type
if not "%TYPE%"=="users" if not "%TYPE%"=="categories" if not "%TYPE%"=="locations" if not "%TYPE%"=="all" (
    echo Invalid type: %TYPE%
    echo Valid types: users, categories, locations, all
    exit /b 1
)

REM Validate count
echo %COUNT% | findstr "^[0-9][0-9]*$" >nul
if !errorlevel! neq 0 (
    echo Invalid count: %COUNT%
    echo Count must be a positive integer
    exit /b 1
)

REM Check if we're in the right directory
if not exist "%PROJECT_ROOT%\go.mod" (
    echo Error: go.mod not found. Are you in the correct directory?
    exit /b 1
)

if not exist "%PROJECT_ROOT%\cmd\seed\main.go" (
    echo Error: Seeder not found at cmd\seed\main.go
    exit /b 1
)

REM Check if .env exists
if not exist "%PROJECT_ROOT%\.env" (
    echo Warning: .env file not found. Make sure environment variables are set.
)

goto :run_seeder

:quick_setup
echo 🚀 Quick Setup: Creating minimal dataset for development...
set TYPE=all
set COUNT=10
goto :run_seeder

:demo_data
echo 🎯 Demo Data: Creating demo dataset...
set TYPE=all
set COUNT=50
goto :run_seeder

:load_test
echo 🏋️ Load Test: Creating large dataset...
echo ⚠️ This will create 500 records of each type. Continue? (y/N)
set /p response=
if /i "!response!"=="y" (
    set TYPE=all
    set COUNT=500
    goto :run_seeder
) else (
    echo Cancelled.
    exit /b 0
)

:run_seeder
echo 🌱 Running seeder: type=%TYPE%, count=%COUNT%

cd /d "%PROJECT_ROOT%"

go run cmd\seed\main.go -type=%TYPE% -count=%COUNT%

if !errorlevel! neq 0 (
    echo ❌ Seeding failed!
    exit /b 1
)

echo ✅ Seeding completed successfully!
goto :end

:show_help
echo Inventory API Seeder Script for Windows
echo ======================================
echo.
echo Usage: %~nx0 [OPTIONS]
echo.
echo Options:
echo   -t, --type TYPE     Type of seed: users, categories, locations, all (default: all)
echo   -c, --count COUNT   Number of records to create (default: 20)
echo   -h, --help          Show this help message
echo.
echo Quick Commands:
echo   %~nx0 --quick-setup    Seed a small dataset for development (10 records each)
echo   %~nx0 --demo-data      Seed demo dataset (50 records each)
echo   %~nx0 --load-test      Seed large dataset for load testing (500 records each)
echo.
echo Examples:
echo   %~nx0                          # Seed all with default count (20)
echo   %~nx0 -t users -c 50          # Seed 50 users
echo   %~nx0 --type=categories --count=30  # Seed 30 categories
echo   %~nx0 --quick-setup            # Quick development setup
echo.

:end
endlocal
