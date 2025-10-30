#!/usr/bin/env python3
"""
删除 repository 子包中的类型断言行
"""
import re
from pathlib import Path

def remove_assertions(file_path):
    """删除类型断言"""
    with open(file_path, 'r', encoding='utf-8') as f:
        lines = f.readlines()
    
    # 过滤掉类型断言行
    filtered = [line for line in lines if not re.match(r'^\s*var\s+_\s+repository\.', line)]
    
    # 移除文件末尾的空行
    while filtered and filtered[-1].strip() == '':
        filtered.pop()
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.writelines(filtered)
    
    return len(lines) - len(filtered)

def main():
    import os
    os.chdir(Path(__file__).parent.parent)
    
    count = 0
    repo_dir = Path('internal/repository')
    
    for go_file in repo_dir.rglob('*.go'):
        if 'gorm' in go_file.name or go_file.parent == repo_dir:
            removed = remove_assertions(go_file)
            if removed > 0:
                print(f"[OK] Removed {removed} assertion(s) from {go_file.relative_to(repo_dir)}")
                count += 1
    
    print(f"\nProcessed {count} files")

if __name__ == '__main__':
    main()


