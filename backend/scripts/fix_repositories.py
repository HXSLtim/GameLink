#!/usr/bin/env python3
"""
批量修复 repository 子包的结构体命名冲突
"""
import re
import os
from pathlib import Path

# 需要修复的 repository 列表
REPOS = [
    'user',
    'player',
    'player_tag',
    'payment',
    'permission',
    'role',
    'stats',
    'operation_log',
    'review',
]

def fix_repository(repo_name):
    """修复单个 repository"""
    file_path = Path(f'internal/repository/{repo_name}/{repo_name}_gorm_repository.go')
    
    if not file_path.exists():
        # 尝试其他可能的文件名
        alt_path = Path(f'internal/repository/{repo_name}/operation_log_gorm_repository.go')
        if alt_path.exists():
            file_path = alt_path
        else:
            print(f"[ERROR] File not found: {file_path}")
            return False
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 获取接口名（首字母大写）
    interface_name = ''.join(word.capitalize() for word in repo_name.split('_')) + 'Repository'
    struct_name = 'gorm' + interface_name
    
    print(f"Fixing {repo_name}: {interface_name} -> {struct_name}")
    
    # 1. 重命名结构体定义
    content = re.sub(
        rf'^type {interface_name} struct',
        f'type {struct_name} struct',
        content,
        flags=re.MULTILINE
    )
    
    # 2. 修改构造函数返回类型
    content = re.sub(
        rf'func New{interface_name}\(db \*gorm\.DB\) \*{interface_name}',
        f'func New{interface_name}(db *gorm.DB) {interface_name}',
        content
    )
    
    # 修改构造函数返回语句
    content = re.sub(
        rf'return &{interface_name}{{db: db}}',
        f'return &{struct_name}{{db: db}}',
        content
    )
    
    # 3. 修改所有方法接收者
    content = re.sub(
        rf'\(r \*{interface_name}\)',
        f'(r *{struct_name})',
        content
    )
    
    # 写回文件
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"[OK] Fixed {repo_name}")
    return True

def main():
    os.chdir(Path(__file__).parent.parent)
    
    print("Fixing repository struct name conflicts...\n")
    
    success_count = 0
    for repo in REPOS:
        if fix_repository(repo):
            success_count += 1
        print()
    
    print(f"\nFixed {success_count}/{len(REPOS)} repositories")

if __name__ == '__main__':
    main()

