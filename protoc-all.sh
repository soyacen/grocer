#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

# 获取脚本所在目录（项目根目录）
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 切换到项目根目录
cd "${SCRIPT_DIR}"

echo "开始编译项目中所有的 proto 文件..."

# 检查并安装必要的protoc插件
if [ ! $(command -v protoc-gen-go) ]; then
    echo "安装 protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    protoc-gen-go --version
fi

# 查找项目中所有的 .proto 文件
PROTO_FILES=( $(find . -name "*.proto" -type f | grep -v "./protoc-all.sh") )

if [ ${#PROTO_FILES[@]} -eq 0 ]; then
    echo "没有找到 .proto 文件"
    exit 0
fi

echo "发现 ${#PROTO_FILES[@]} 个 proto 文件..."

# 逐个编译每个proto文件，保持相对路径结构
for proto_file in "${PROTO_FILES[@]}"; do
    dir=$(dirname "$proto_file")
    echo "正在编译 $proto_file ..."
    
    protoc \
      --proto_path=. \
      --go_out=. \
      --go_opt=paths=source_relative \
      "$proto_file"
done

echo "编译完成！生成的文件位于对应的 proto 文件所在目录。"