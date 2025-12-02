#!/bin/bash

# 捕获 Ctrl+C 并优雅退出
trap 'echo -e "\n用户中断，正在退出..."; exit 130' INT TERM

# 显示用法
show_usage() {
    echo "用法: $0 <7z文件目录> <解压目标目录>"
    echo "功能: 将 <7z文件目录> 下的每个 .7z 文件解压为独立子目录"
    echo "依赖: p7zip-full (7z 命令)"
    echo "示例:"
    echo "  $0 ./compressed ./restored"
    echo "  $0 ../novel-library/books ../novel-database/books"
}

# 检查参数
if [ $# -ne 2 ]; then
    echo "错误: 需要指定 7z 文件目录 和 解压目标目录" >&2
    show_usage
    exit 1
fi

SOURCE_DIR=$(readlink -f "${1%/}")
TARGET_DIR=$(readlink -f "${2%/}")

# 创建目标目录
mkdir -p "$TARGET_DIR"

# 检查 7z 是否安装
if ! command -v 7z &> /dev/null; then
    echo "错误: '7z' 未安装，请运行: sudo apt install p7zip-full" >&2
    exit 1
fi

# 获取所有 .7z 文件（不递归，仅一级）
mapfile -d '' sevenz_files < <(find "$SOURCE_DIR" -mindepth 1 -maxdepth 1 -type f -name "*.7z" -print0 2>/dev/null)

total=${#sevenz_files[@]}

if [ "$total" -eq 0 ]; then
    echo "警告: 目录 '$SOURCE_DIR' 中没有找到 .7z 文件"
    exit 0
fi

echo "找到 $total 个 .7z 文件，开始解压..."
echo

processed=0
for file in "${sevenz_files[@]}"; do
    filename=$(basename "$file")
    dirname="${filename%.7z}"  # 移除 .7z 后缀

    echo "[$((processed + 1))/$total] 解压: $filename → $dirname/"

    # 使用 7z 解压：
    # -x       : 解压（保留完整路径）
    # -o<TARGET_DIR> : 指定输出目录（注意：-o 后不能有空格！）
    # -bd      : 禁用进度条
    # -y       : 自动确认覆盖
    if 7z x -o"$TARGET_DIR" -bd -y "$file" >/dev/null 2>&1; then
        echo "✓ 成功: $filename → $dirname/"
        ((processed++))
    else
        echo "✗ 失败: $filename"
    fi
done

echo
echo "解压完成！"
echo "成功解压: $processed/$total 个文件"
echo "输出位置: $TARGET_DIR