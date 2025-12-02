#!/bin/bash

# 捕获 Ctrl+C 并优雅退出
trap 'echo -e "\n用户中断，正在退出..."; exit 130' INT TERM

# 显示用法
show_usage() {
    echo "用法: $0 <源目录> <目标目录>"
    echo "功能: 将 <源目录> 下所有包含 .txt 的「叶子目录」压缩为 .7z"
    echo "     输出路径保持原相对结构，如 books/A/Book → target/A/Book.7z"
    echo "依赖: p7zip-full (7z 命令)"
    echo "示例:"
    echo "  $0 ./data/books ./novel-archive"
}

# 检查参数
if [ $# -ne 2 ]; then
    echo "错误: 需要指定源目录和目标目录" >&2
    show_usage
    exit 1
fi

SOURCE_DIR=$(readlink -f "${1%/}")
TARGET_BASE=$(readlink -f "${2%/}")

# 创建目标根目录
mkdir -p "$TARGET_BASE"

# 检查 7z 是否安装
if ! command -v 7z &> /dev/null; then
    echo "错误: '7z' 未安装，请运行: sudo apt install p7zip-full" >&2
    exit 1
fi

# 找出所有「叶子目录」：包含至少一个 .txt，且无子目录
echo " 正在扫描源目录中的小说（叶子目录）..."
mapfile -d '' leaf_dirs < <(
    find "$SOURCE_DIR" -type d -exec sh -c '
        for d; do
            # 检查是否有 .txt 文件
            if [ -n "$(find "$d" -maxdepth 1 -name "*.txt" -print -quit)" ]; then
                # 检查是否有子目录
                if [ -z "$(find "$d" -mindepth 1 -maxdepth 1 -type d -print -quit)" ]; then
                    printf "%s\0" "$d"
                fi
            fi
        done
    ' _ {} +
)

total=${#leaf_dirs[@]}

if [ "$total" -eq 0 ]; then
    echo "警告: 未找到符合条件的小说目录（需包含 .txt 且无子目录）"
    exit 0
fi

echo "✅ 找到 $total 个小说目录，开始压缩..."
echo

processed=0
for full_path in "${leaf_dirs[@]}"; do
    # 计算相对于 SOURCE_DIR 的路径（如 D/DOTA之最强血脉）
    rel_path="${full_path#$SOURCE_DIR/}"

    # 构造目标 .7z 路径（保留层级）
    target_7z="$TARGET_BASE/${rel_path}.7z"
    target_dir=$(dirname "$target_7z")

    # 创建目标父目录
    mkdir -p "$target_dir"

    novel_name=$(basename "$rel_path")
    echo "[$((processed + 1))/$total] 压缩: $rel_path → ${rel_path}.7z"

    # 进入 SOURCE_DIR 执行 7z（避免绝对路径写入 archive）
    if (cd "$SOURCE_DIR" && 7z a -t7z -m0=lzma2 -mx=5 -mmt=on -ms=on -bd -y "$target_7z" "$rel_path" >/dev/null 2>&1); then
        size=$(du -h "$target_7z" 2>/dev/null | cut -f1)
        echo "✓ 成功: ${rel_path}.7z ($size)"
        ((processed++))
    else
        echo "✗ 失败: $rel_path"
    fi
done

echo
echo "处理完成！"
echo "成功压缩: $processed/$total 个小说"
echo "输出位置: $TARGET_BASE"