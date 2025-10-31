#!/usr/bin/env python3
"""
为所有 repository 子包添加 repository 包导入
"""
import re
import os
from pathlib import Path

REPOS = [
    'player',
    'player_tag',
    'payment',
    'permission',
    'role',
    'stats',
    'operation_log',
    'review',
]

def add_import(repo_name):
    """为单个 repository 添加导入"""
    file_path = Path(f'internal/repository/{repo_name}/{repo_name}_gorm_repository.go')
    
    if not file_path.exists():
        # operation_log 有不同的文件名
        file_path = Path(f'internal/repository/{repo_name}/operation_log_gorm_repository.go')
        if not file_path.exists():
            print(f"[ERROR] File not found: {repo_name}")
            return False
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 检查是否已经有导入
    if '"gamelink/internal/repository"' in content:
        print(f"[SKIP] {repo_name} already has repository import")
        return True
    
    # 在最后一个内部导入后添加
    pattern = r'("gamelink/internal/model")'
    replacement = r'\1\n\t"gamelink/internal/repository"'
    
    if '"gamelink/internal/model"' not in content:
        print(f"[WARN] {repo_name} doesn't import model package, trying different pattern")
        # 尝试在 import 块结束前添加
        pattern = r'(\n\)\n\n)'
        replacement = r'\t"gamelink/internal/repository"\n\1'
    
    content = re.sub(pattern, replacement, content, count=1)
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"[OK] Added import to {repo_name}")
    return True

def main():
    os.chdir(Path(__file__).parent.parent)
    
    print("Adding repository imports...\n")
    
    success_count = 0
    for repo in REPOS:
        if add_import(repo):
            success_count += 1
    
    print(f"\nProcessed {success_count}/{len(REPOS)} repositories")

if __name__ == '__main__':
    main()


