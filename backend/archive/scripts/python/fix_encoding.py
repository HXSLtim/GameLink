#!/usr/bin/env python3
"""
修复 repository 文件的编码问题和接口返回类型
"""
import re
from pathlib import Path

def fix_file(file_path):
    """修复单个文件"""
    try:
        # 尝试多种编码读取
        content = None
        for encoding in ['utf-8', 'gb2312', 'gbk', 'latin-1']:
            try:
                with open(file_path, 'r', encoding=encoding) as f:
                    content = f.read()
                break
            except:
                continue
        
        if content is None:
            print(f"[ERROR] Cannot read {file_path}")
            return False
        
        # 修复常见的乱码模式
        content = content.replace('�?', '。\n')
        
        # 修复返回类型（添加 repository. 前缀）
        # 匹配模式：func NewXxxRepository(...) XxxRepository {
        pattern = r'func (New\w+Repository\([^)]+\)) (\w+Repository) \{'
        
        def replace_return_type(match):
            func_sig = match.group(1)
            repo_type = match.group(2)
            return f'func {func_sig} repository.{repo_type} {{'
        
        content = re.sub(pattern, replace_return_type, content)
        
        # 保存为UTF-8
        with open(file_path, 'w', encoding='utf-8', newline='\n') as f:
            f.write(content)
        
        return True
    except Exception as e:
        print(f"[ERROR] {file_path}: {e}")
        return False

def main():
    import os
    os.chdir(Path(__file__).parent.parent)
    
    repo_dir = Path('internal/repository')
    count = 0
    success = 0
    
    for go_file in repo_dir.rglob('*_gorm_repository.go'):
        count += 1
        if fix_file(go_file):
            print(f"[OK] Fixed {go_file.relative_to(repo_dir)}")
            success += 1
    
    # 也修复非gorm的repository文件
    for go_file in repo_dir.rglob('*_repository.go'):
        if 'gorm' not in go_file.name:
            count += 1
            if fix_file(go_file):
                print(f"[OK] Fixed {go_file.relative_to(repo_dir)}")
                success += 1
    
    print(f"\nProcessed {count} files, {success} successful")

if __name__ == '__main__':
    main()


