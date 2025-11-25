#!/bin/bash

# ============================================
# M30 信号数据批处理脚本
# 用途：快速生成和导入M30信号测试数据
# ============================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WEB_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$WEB_DIR")"

# ============================================
# 函数定义
# ============================================

print_header() {
  echo -e "${BLUE}╔═══════════════════════════════════════════════════════╗${NC}"
  echo -e "${BLUE}║ PP Pro Desk - M30 信号数据批处理工具                 ║${NC}"
  echo -e "${BLUE}╚═══════════════════════════════════════════════════════╝${NC}"
  echo
}

print_help() {
  cat << EOF
用法：bash generate-signals.sh [命令] [选项]

命令：
  generate  生成信号数据文件
  display   显示生成的数据（JSON格式）
  copy      复制Vue代码到剪贴板（macOS）
  insert    交互式插入数据

选项：
  -c, --count <数量>      信号数量 (默认: 20)
  -f, --format <格式>     输出格式: json|csv|sql|js|vue (默认: json)
  -h, --help              显示帮助信息

示例：
  bash generate-signals.sh generate -c 20                # 生成20条JSON信号
  bash generate-signals.sh generate -c 50 -f vue        # 生成50条Vue格式
  bash generate-signals.sh generate -c 30 -f csv        # 生成30条CSV格式
  bash generate-signals.sh display                       # 显示最新生成的数据
  bash generate-signals.sh copy                          # 复制Vue代码到剪贴板
  bash generate-signals.sh insert                        # 交互式插入数据

EOF
}

print_success() {
  echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
  echo -e "${RED}❌ $1${NC}"
}

print_warning() {
  echo -e "${YELLOW}⚠️  $1${NC}"
}

print_info() {
  echo -e "${BLUE}ℹ️  $1${NC}"
}

# 检查Node.js
check_node() {
  if ! command -v node &> /dev/null; then
    print_error "未找到 Node.js，请先安装"
    exit 1
  fi
  print_success "Node.js 已就绪 ($(node -v))"
}

# 生成信号数据
generate_signals() {
  local count=$1
  local format=$2
  
  print_info "正在生成 $count 条 M30 信号数据 ($format 格式)..."
  cd "$WEB_DIR"
  node scripts/generate-signals.js "$count" "$format"
}

# 显示生成的数据
display_signals() {
  local latest=$(ls -t generated/signals-*.json 2>/dev/null | head -1)
  
  if [ -z "$latest" ]; then
    print_error "未找到生成的信号数据，请先运行: bash generate-signals.sh generate"
    return 1
  fi
  
  print_info "显示最新数据: $(basename $latest)"
  echo
  cat "$latest" | head -50
  echo
  echo "..."
  echo
  print_info "完整文件位置: $latest"
}

# 复制到剪贴板（macOS）
copy_to_clipboard() {
  local latest=$(ls -t generated/signals-*.js 2>/dev/null | head -1)
  
  if [ -z "$latest" ]; then
    print_info "未找到 .js 格式的数据，正在生成..."
    generate_signals 20 vue
    latest=$(ls -t generated/signals-*.js 2>/dev/null | head -1)
  fi
  
  if [[ "$OSTYPE" == "darwin"* ]]; then
    cat "$latest" | pbcopy
    print_success "Vue 代码已复制到剪贴板"
    print_info "可直接粘贴到 TradeView.vue 的 signalFeed 变量"
  else
    print_error "本功能仅支持 macOS"
    print_info "请手动复制文件: $latest"
  fi
}

# 交互式插入
interactive_insert() {
  echo
  print_info "交互式信号数据插入工具"
  echo
  
  read -p "请输入要生成的信号数量 (默认: 20): " count
  count=${count:-20}
  
  echo "选择输出格式:"
  echo "  1. JSON (推荐用于测试)"
  echo "  2. CSV (用于电子表格)"
  echo "  3. SQL (用于数据库)"
  echo "  4. JavaScript"
  echo "  5. Vue (推荐用于组件)"
  read -p "请选择 [1-5] (默认: 5): " format_choice
  format_choice=${format_choice:-5}
  
  case $format_choice in
    1) format="json" ;;
    2) format="csv" ;;
    3) format="sql" ;;
    4) format="js" ;;
    5) format="vue" ;;
    *) 
      print_error "无效选择"
      return 1
      ;;
  esac
  
  echo
  generate_signals "$count" "$format"
}

# ============================================
# 主逻辑
# ============================================

main() {
  print_header
  
  if [ $# -eq 0 ]; then
    print_help
    return 0
  fi
  
  check_node
  echo
  
  local command=$1
  shift
  
  case $command in
    generate)
      local count=20
      local format="json"
      
      while [[ $# -gt 0 ]]; do
        case $1 in
          -c|--count)
            count="$2"
            shift 2
            ;;
          -f|--format)
            format="$2"
            shift 2
            ;;
          -h|--help)
            print_help
            return 0
            ;;
          *)
            print_error "未知选项: $1"
            return 1
            ;;
        esac
      done
      
      generate_signals "$count" "$format"
      ;;
      
    display)
      display_signals
      ;;
      
    copy)
      copy_to_clipboard
      ;;
      
    insert)
      interactive_insert
      ;;
      
    -h|--help|help)
      print_help
      ;;
      
    *)
      print_error "未知命令: $command"
      echo
      print_help
      return 1
      ;;
  esac
}

# 执行主函数
main "$@"
