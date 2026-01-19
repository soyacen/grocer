#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

# 检查并安装必要的protoc插件
if ! command -v protoc-gen-go &> /dev/null; then
    echo "安装 protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    protoc-gen-go --version
fi

if ! command -v protoc-gen-goose &> /dev/null; then
    echo "安装 protoc-gen-goose..."
	go install github.com/soyacen/goose/cmd/protoc-gen-goose@latest
	protoc-gen-goose --version
fi

if ! command -v protoc-gen-validate-go &> /dev/null; then
    echo "安装 protoc-gen-validate-go..."
	go install github.com/envoyproxy/protoc-gen-validate/cmd/protoc-gen-validate-go@latest
	protoc-gen-validate-go --version
fi

# 获取脚本所在目录（当前模块的根目录）
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 切换到当前模块根目录
cd "${SCRIPT_DIR}"

echo "开始编译当前目录及子目录下的 proto 文件..."

# 使用find命令查找当前目录及子目录下的所有 .proto 文件
PROTO_FILES=()
while IFS= read -r -d '' file; do
    if [[ "$file" != *"third_party"* ]]; then
        PROTO_FILES+=("$file")
    fi
done < <(find . -name "*.proto" -type f -print0)

if [ ${#PROTO_FILES[@]} -eq 0 ]; then
    echo "没有找到 .proto 文件"
    exit 0
fi

echo "发现 ${#PROTO_FILES[@]} 个 proto 文件..."

# 编译所有proto文件，使用当前目录和third_party作为proto_path，这样导入路径可以正确解析
echo "正在编译 proto 文件..."
protoc \
  --proto_path=. \
  --proto_path=../../third_party \
  --go_out=. \
  --go_opt=paths=source_relative \
  --goose_out=. \
  --goose_opt=paths=source_relative \
  "${PROTO_FILES[@]}"

echo "编译完成！生成的文件位于对应的 proto 文件所在目录。"