#!/usr/bin/env python3
"""
UTF-8 编码检测脚本
用于扫描整个项目的编码问题
"""

import os
import sys

def check_file_encoding(filepath):
    """检查文件编码"""
    try:
        with open(filepath, 'rb') as f:
            content = f.read()
        
        # 尝试用 UTF-8 解码
        try:
            content.decode('utf-8')
            return "UTF-8", True
        except UnicodeDecodeError as e:
            # 检测失败的字符位置
            error_pos = e.start if hasattr(e, 'start') else 0
            error_bytes = content[error_pos:error_pos+10] if error_pos < len(content) else b""
            return f"Non-UTF-8 (error at pos {error_pos}, bytes: {error_bytes.hex()})", False
            
    except Exception as e:
        return f"Error: {e}", False

def scan_directory(directory):
    """扫描目录中的 Go 文件"""
    problem_files = []
    total_files = 0
    
    for root, dirs, files in os.walk(directory):
        # 跳过隐藏目录和缓存目录
        dirs[:] = [d for d in dirs if not d.startswith('.') and d not in ['vendor', 'node_modules']]
        
        for file in files:
            if file.endswith('.go'):
                filepath = os.path.join(root, file)
                total_files += 1
                
                encoding, is_utf8 = check_file_encoding(filepath)
                if not is_utf8:
                    problem_files.append((filepath, encoding))
                    print(f"❌ {filepath}: {encoding}")
                else:
                    # 可选：显示正常文件
                    # print(f"✅ {filepath}: {encoding}")
                    pass
    
    return problem_files, total_files

def main():
    """主函数"""
    # 扫描的目录
    directories_to_scan = [
        "/mnt/c/Users/a2778/Desktop/code/GameLink/backend/internal",
        "/mnt/c/Users/a2778/Desktop/code/GameLink/backend/cmd",
        "/mnt/c/Users/a2778/Desktop/code/GameLink/backend/configs",
        "/mnt/c/Users/a2778/Desktop/code/GameLink/backend/docs"
    ]
    
    all_problem_files = []
    total_files_scanned = 0
    
    print("扫描 UTF-8 编码问题...")
    print("=" * 80)
    
    for directory in directories_to_scan:
        if os.path.exists(directory):
            print(f"\n扫描目录: {directory}")
            print("-" * 40)
            problem_files, total_files = scan_directory(directory)
            all_problem_files.extend(problem_files)
            total_files_scanned += total_files
            
            if not problem_files:
                print(f"✅ 目录 {directory} 中的所有文件都是 UTF-8 编码")
        else:
            print(f"⚠️  目录不存在: {directory}")
    
    print("\n" + "=" * 80)
    print("扫描完成!")
    print(f"总共扫描了 {total_files_scanned} 个 Go 文件")
    print(f"发现 {len(all_problem_files)} 个文件存在编码问题")
    
    if all_problem_files:
        print("\n存在编码问题的文件:")
        for filepath, encoding in all_problem_files:
            print(f"  ❌ {filepath}: {encoding}")
        return 1
    else:
        print("✅ 所有文件都是 UTF-8 编码!")
        return 0

if __name__ == "__main__":
    sys.exit(main())