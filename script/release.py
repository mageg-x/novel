#!/usr/bin/env python3
import os
import sys
import re
import json
import hashlib
import urllib.parse
import requests
from pathlib import Path
from pypinyin import lazy_pinyin, Style

# 强制要求 ARCHIVE_DIR 通过命令行参数传入
if len(sys.argv) < 2:
    print("请提供归档目录作为命令行参数（例如：./script.py ./my_archive）", file=sys.stderr)
    sys.exit(1)

ARCHIVE_DIR = Path(sys.argv[1])

GITHUB_REPO = os.getenv("GITHUB_REPO")
GITHUB_TOKEN = os.getenv("GITHUB_TOKEN") or os.popen("gh auth token").read().strip()

if not GITHUB_REPO:
    print("请设置环境变量 GITHUB_REPO", file=sys.stderr)
    sys.exit(1)

API_BASE = f"https://api.github.com/repos/{GITHUB_REPO}"
HEADERS = {
    "Authorization": f"token {GITHUB_TOKEN}",
    "Accept": "application/vnd.github.v3+json",
    "X-GitHub-Api-Version": "2022-11-28"
}

def get_release_tag(dirname: str) -> str:
    """直接使用目录名作为 tag（支持中文）"""
    return dirname  # 不再取首字母！

def to_pinyin_name(filename: str) -> str:
    stem = Path(filename).stem
    ext = Path(filename).suffix
    parts = []
    for char in stem:
        if char.isascii() and (char.isalnum() or char in "._-"):
            parts.append(char)
        else:
            py = lazy_pinyin(char, style=Style.NORMAL)
            parts.extend([p.capitalize() for p in py])
    safe_stem = "".join(parts)
    safe_stem = re.sub(r"[^a-zA-Z0-9]+", "_", safe_stem)
    safe_stem = re.sub(r"^_+|_+$", "", safe_stem)
    return safe_stem + ext

def calc_sha256(path: Path) -> str:
    hash_sha256 = hashlib.sha256()
    with open(path, "rb") as f:
        for chunk in iter(lambda: f.read(4096), b""):
            hash_sha256.update(chunk)
    return hash_sha256.hexdigest()

def ensure_release(tag: str, title: str):
    # 必须对 tag 进行 URL 编码
    encoded_tag = urllib.parse.quote(tag, safe='')
    url = f"{API_BASE}/releases/tags/{encoded_tag}"
    r = requests.get(url, headers=HEADERS)
    if r.status_code == 200:
        return r.json()
    elif r.status_code == 404:
        data = {"tag_name": tag, "name": title, "body": ""}
        r = requests.post(f"{API_BASE}/releases", json=data, headers=HEADERS)
        r.raise_for_status()
        return r.json()
    else:
        r.raise_for_status()

def get_online_asset_hashes(release_id: int):
    """从 assets 的 digest 字段提取 SHA256 哈希"""
    assets_url = f"{API_BASE}/releases/{release_id}/assets"
    r = requests.get(assets_url, headers=HEADERS)
    if r.status_code != 200:
        return {}
    hashes = {}
    for asset in r.json():
        digest = asset.get("digest", "")
        if digest.startswith("sha256:"):
            sha = digest[7:]  # 去掉 "sha256:"
            hashes[asset["name"]] = sha
    return hashes

def upload_asset(release_info: dict, local_path: Path, remote_name: str):
    release_id = release_info["id"]
    upload_url = release_info["upload_url"].split("{")[0]

    # 获取线上所有 asset 的 SHA256（来自 digest）
    online_hashes = get_online_asset_hashes(release_id)

    local_hash = calc_sha256(local_path)

    # 如果哈希一致，跳过
    if online_hashes.get(remote_name) == local_hash:
        print(f"  ➖ {remote_name} (内容未变，跳过)")
        return

    # 删除同名旧 asset
    assets_url = f"{API_BASE}/releases/{release_id}/assets"
    r = requests.get(assets_url, headers=HEADERS)
    if r.status_code == 200:
        for asset in r.json():
            if asset["name"] == remote_name:
                del_resp = requests.delete(asset["url"], headers=HEADERS)
                if del_resp.status_code == 204:
                    print(f"  ♻️ 替换: {remote_name}")
                else:
                    print(f"  ⚠️ 删除失败 ({del_resp.status_code})")
                break

    # 上传新文件
    with open(local_path, "rb") as f:
        resp = requests.post(
            f"{upload_url}?name={requests.utils.quote(remote_name)}",
            headers={**HEADERS, "Content-Type": "application/octet-stream"},
            data=f,
        )
    if resp.status_code == 201:
        print(f"  ✅ {remote_name}")
    else:
        print(f"  ❌ {remote_name} → {resp.status_code}")

# === 主流程 ===
for subdir in sorted(ARCHIVE_DIR.iterdir()):
    if not subdir.is_dir():
        continue

    tag = get_release_tag(subdir.name)  # 直接用目录名，如 "一"
    print(f"\n {subdir.name} → Release: {tag}")
    release_info = ensure_release(tag, subdir.name)

    for file_path in sorted(subdir.glob("*.7z")):
        safe_name = to_pinyin_name(file_path.name)
        upload_asset(release_info, file_path, safe_name)

print("\n✅ 同步完成！")