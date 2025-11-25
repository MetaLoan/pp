@echo off
REM ============================================
REM M30 信号数据批处理脚本 (Windows)
REM 用途：快速生成和导入M30信号测试数据
REM ============================================

setlocal enabledelayedexpansion
chcp 65001 >nul 2>&1

REM 获取脚本目录
set "SCRIPT_DIR=%~dp0"
set "WEB_DIR=%SCRIPT_DIR%.."
set "PROJECT_ROOT=%WEB_DIR%.."

REM ============================================
REM 函数定义
REM ============================================

:print_header
echo.
echo ╔═══════════════════════════════════════════════════════╗
echo ║ PP Pro Desk - M30 信号数据批处理工具                 ║
echo ╚═══════════════════════════════════════════════════════╝
echo.
exit /b 0

:print_help
echo 用法: generate-signals.bat [命令] [选项]
echo.
echo 命令:
echo   generate  生成信号数据文件
echo   display   显示生成的数据（JSON格式）
echo   insert    交互式插入数据
echo.
echo 选项:
echo   -c, --count ^<数量^>      信号数量 (默认: 20^)
echo   -f, --format ^<格式^>     输出格式: json^|csv^|sql^|js^|vue (默认: json^)
echo   -h, --help              显示帮助信息
echo.
echo 示例:
echo   generate-signals.bat generate -c 20                REM 生成20条JSON信号
echo   generate-signals.bat generate -c 50 -f vue        REM 生成50条Vue格式
echo   generate-signals.bat generate -c 30 -f csv        REM 生成30条CSV格式
echo   generate-signals.bat display                       REM 显示最新生成的数据
echo   generate-signals.bat insert                        REM 交互式插入数据
echo.
exit /b 0

:check_node
where node >nul 2>&1
if errorlevel 1 (
  echo ❌ 未找到 Node.js，请先安装
  exit /b 1
)
for /f "tokens=*" %%i in ('node -v') do set "NODE_VERSION=%%i"
echo ✅ Node.js 已就绪 (!NODE_VERSION!)
exit /b 0

:generate_signals
set "count=%1"
set "format=%2"
echo.
echo ℹ️  正在生成 !count! 条 M30 信号数据 (!format! 格式^)...
echo.
cd /d "%WEB_DIR%"
call node scripts\generate-signals.js !count! !format!
exit /b 0

:display_signals
echo.
echo ℹ️  正在查找最新生成的信号数据...
echo.
for /f "tokens=*" %%F in ('dir /b /o-d generated\signals-*.json 2^>nul ^| findstr /r ".*" ^| (for /f %%A in ('findstr /r ".*"') do @echo %%A ^& exit /b)') do (
  set "latest=%%F"
  goto :found_file
)

echo ❌ 未找到生成的信号数据
echo 请先运行: generate-signals.bat generate
exit /b 1

:found_file
set "filepath=%WEB_DIR%\generated\!latest!"
echo 显示最新数据: !latest!
echo.
if exist "!filepath!" (
  type "!filepath!" | more
  echo.
  echo ℹ️  完整文件位置: !filepath!
) else (
  echo ❌ 文件不存在: !filepath!
)
exit /b 0

:interactive_insert
echo.
echo ℹ️  交互式信号数据插入工具
echo.

set "count=20"
set "format=vue"

set /p count="请输入要生成的信号数量 (默认: 20^): "
if "!count!"=="" set "count=20"

echo.
echo 选择输出格式:
echo   1. JSON (推荐用于测试^)
echo   2. CSV (用于电子表格^)
echo   3. SQL (用于数据库^)
echo   4. JavaScript
echo   5. Vue (推荐用于组件^)

set /p format_choice="请选择 [1-5] (默认: 5^): "
if "!format_choice!"=="" set "format_choice=5"

if "!format_choice!"=="1" set "format=json"
if "!format_choice!"=="2" set "format=csv"
if "!format_choice!"=="3" set "format=sql"
if "!format_choice!"=="4" set "format=js"
if "!format_choice!"=="5" set "format=vue"

echo.
call :generate_signals !count! !format!
exit /b 0

REM ============================================
REM 主逻辑
REM ============================================

call :print_header

if "%1"=="" (
  call :print_help
  exit /b 0
)

call :check_node
if errorlevel 1 exit /b 1

set "command=%1"
shift

if /i "!command!"=="generate" (
  set "count=20"
  set "format=json"
  
  :parse_args
  if "!1!"=="" goto :end_parse
  if /i "!1!"=="-c" (
    set "count=!2!"
    shift
    shift
    goto :parse_args
  )
  if /i "!1!"=="--count" (
    set "count=!2!"
    shift
    shift
    goto :parse_args
  )
  if /i "!1!"=="-f" (
    set "format=!2!"
    shift
    shift
    goto :parse_args
  )
  if /i "!1!"=="--format" (
    set "format=!2!"
    shift
    shift
    goto :parse_args
  )
  if /i "!1!"=="-h" (
    call :print_help
    exit /b 0
  )
  if /i "!1!"=="--help" (
    call :print_help
    exit /b 0
  )
  shift
  goto :parse_args
  
  :end_parse
  call :generate_signals !count! !format!
) else if /i "!command!"=="display" (
  call :display_signals
) else if /i "!command!"=="insert" (
  call :interactive_insert
) else if /i "!command!"=="-h" (
  call :print_help
) else if /i "!command!"=="--help" (
  call :print_help
) else if /i "!command!"=="help" (
  call :print_help
) else (
  echo ❌ 未知命令: !command!
  echo.
  call :print_help
  exit /b 1
)

endlocal
