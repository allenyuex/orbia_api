#!/bin/bash

# 数据库初始化脚本
# 用于执行 sql/init.sql 文件来初始化数据库结构

# 设置脚本错误时退出
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 获取脚本所在目录的父目录（项目根目录）
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
SQL_FILE="$PROJECT_ROOT/sql/init.sql"

# 默认数据库配置（从 config.yaml 读取的默认值）
DB_HOST=${DB_HOST:-"127.0.0.1"}
DB_PORT=${DB_PORT:-"3306"}
DB_USER=${DB_USER:-"root"}
DB_PASSWORD=${DB_PASSWORD:-"root123"}
DB_NAME=${DB_NAME:-"orbia"}

print_info "开始数据库初始化..."
print_info "项目根目录: $PROJECT_ROOT"
print_info "SQL文件路径: $SQL_FILE"

# 检查 SQL 文件是否存在
if [ ! -f "$SQL_FILE" ]; then
    print_error "SQL文件不存在: $SQL_FILE"
    exit 1
fi

# 检查 mysql 命令是否可用
if ! command -v mysql &> /dev/null; then
    print_error "mysql 命令未找到，请确保已安装 MySQL 客户端"
    print_info "在 macOS 上可以通过以下方式安装："
    print_info "  brew install mysql-client"
    print_info "  或者安装完整的 MySQL: brew install mysql"
    exit 1
fi

# 显示连接信息
print_info "数据库连接信息:"
print_info "  主机: $DB_HOST"
print_info "  端口: $DB_PORT"
print_info "  用户: $DB_USER"
print_info "  数据库: $DB_NAME"

# 测试数据库连接
print_info "测试数据库连接..."
if ! mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1;" &> /dev/null; then
    print_error "无法连接到数据库，请检查："
    print_error "  1. MySQL 服务是否正在运行"
    print_error "  2. 数据库连接参数是否正确"
    print_error "  3. 用户权限是否足够"
    print_info ""
    print_info "你可以通过环境变量覆盖默认配置："
    print_info "  export DB_HOST=your_host"
    print_info "  export DB_PORT=your_port"
    print_info "  export DB_USER=your_user"
    print_info "  export DB_PASSWORD=your_password"
    print_info "  export DB_NAME=your_database"
    exit 1
fi

print_info "数据库连接成功！"

# 执行 SQL 文件
print_info "执行数据库初始化脚本..."
if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" < "$SQL_FILE"; then
    print_info "数据库初始化完成！"
    
    # 显示创建的表
    print_info "已创建的表："
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "USE $DB_NAME; SHOW TABLES;" 2>/dev/null | grep -v "Tables_in_" | while read table; do
        if [ -n "$table" ]; then
            print_info "  - $table"
        fi
    done
else
    print_error "数据库初始化失败！"
    exit 1
fi

print_info "数据库初始化脚本执行完成！"

# 生成数据库模型代码
print_info "开始生成数据库模型代码..."

# 检查 Go 是否已安装
if ! command -v go &> /dev/null; then
    print_error "Go 命令未找到，请确保已安装 Go"
    exit 1
fi

# 检查 db_gen.go 文件是否存在
DB_GEN_FILE="$SCRIPT_DIR/db_gen.go"
if [ ! -f "$DB_GEN_FILE" ]; then
    print_error "代码生成脚本不存在: $DB_GEN_FILE"
    exit 1
fi

# 设置环境变量（如果需要的话）
export DB_HOST="$DB_HOST"
export DB_PORT="$DB_PORT"
export DB_USER="$DB_USER"
export DB_PASSWORD="$DB_PASSWORD"
export DB_NAME="$DB_NAME"

# 执行代码生成脚本
print_info "执行代码生成脚本..."
cd "$PROJECT_ROOT"
if go run "$DB_GEN_FILE"; then
    print_info "数据库模型代码生成完成！"
    
    # 显示生成的文件
    MODEL_DIR="$PROJECT_ROOT/biz/dal/mysql"
    if [ -d "$MODEL_DIR" ]; then
        print_info "生成的模型文件："
        find "$MODEL_DIR" -name "*.go" -type f | while read file; do
            print_info "  - $(basename "$file")"
        done
    fi
else
    print_error "数据库模型代码生成失败！"
    exit 1
fi

print_info "完整的数据库初始化和代码生成流程执行完成！"