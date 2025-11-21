#!/usr/bin/env python3
"""
测试文件合并脚本
用于将分散的测试文件合并为统一的主测试文件
"""

import os
import re
import sys
import shutil
from pathlib import Path

# ANSI颜色代码
GREEN = '\033[92m'
YELLOW = '\033[93m'
RED = '\033[91m'
BLUE = '\033[94m'
RESET = '\033[0m'

def print_colored(color, message):
    """打印彩色文本"""
    try:
        print(f"{color}{message}{RESET}")
    except UnicodeEncodeError:
        # Windows CMD不支持Unicode，移除emoji
        message = message.encode('ascii', errors='ignore').decode('ascii')
        print(f"{color}{message}{RESET}")

def extract_imports(content):
    """提取import语句"""
    imports = []
    
    # 匹配单行import
    single_import = re.search(r'^import "([^"]+)"', content, re.MULTILINE)
    if single_import:
        imports.append(single_import.group(1))
    
    # 匹配多行import
    multi_import = re.search(r'import \((.*?)\)', content, re.DOTALL)
    if multi_import:
        import_lines = multi_import.group(1).strip().split('\n')
        for line in import_lines:
            line = line.strip()
            if line and not line.startswith('//'):
                # 移除引号
                line = line.strip('"')
                imports.append(line)
    
    return imports

def extract_test_functions(content):
    """提取所有test函数"""
    # 匹配test函数（包括子测试）
    pattern = r'func (Test\w+\([^)]*\)) \{(?:[^{}]*|\{[^}]*\})*?\n\}'
    matches = re.finditer(pattern, content, re.MULTILINE | re.DOTALL)
    
    functions = []
    for match in matches:
        func_code = match.group(0)
        # 检查是否是子测试（包含t.Run）
        if 't.Run' in func_code:
            # 提取子测试
            subtest_pattern = r'(func Test\w+\(t \*testing\.T\) \{[^}]*t\.Run[^}]*\})'
            sub_matches = re.finditer(subtest_pattern, func_code, re.MULTILINE | re.DOTALL)
            for sub_match in sub_matches:
                functions.append(sub_match.group(0))
        else:
            functions.append(func_code)
    
    return functions

def extract_package_name(content):
    """提取package名称"""
    match = re.search(r'^package (\w+)', content, re.MULTILINE)
    return match.group(1) if match else "main"

def merge_test_files(module_path, files_to_merge, output_file):
    """
    合并测试文件
    
    Args:
        module_path: 模块路径（如 internal/service/commission）
        files_to_merge: 要合并的文件列表
        output_file: 输出文件名
    """
    print_colored(BLUE, f"\n📦 正在处理模块: {module_path}")
    
    full_path = Path(module_path)
    if not full_path.exists():
        print_colored(RED, f"❌ 路径不存在: {module_path}")
        return False
    
    # 收集所有文件内容
    package_name = None
    all_imports = set()
    all_functions = []
    file_count = 0
    
    for file_name in files_to_merge:
        file_path = full_path / file_name
        if not file_path.exists():
            print_colored(YELLOW, f"[WARN] 文件不存在: {file_name}")
            continue
        
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            print_colored(GREEN, f"  [OK] 读取: {file_name}")
            
            # 提取package名称
            if package_name is None:
                package_name = extract_package_name(content)
            
            # 提取imports
            imports = extract_imports(content)
            all_imports.update(imports)
            
            # 提取test函数
            functions = extract_test_functions(content)
            all_functions.extend(functions)
            
            file_count += 1
            
        except Exception as e:
            print_colored(RED, f"[ERROR] 读取失败 {file_name}: {e}")
            return False
    
    if file_count == 0:
        print_colored(RED, "❌ 没有成功读取任何文件")
        return False
    
    # 构建合并后的内容
    merged_content = []
    
    # 1. Package声明
    merged_content.append(f"package {package_name}\n")
    
    # 2. Import语句
    if all_imports:
        merged_content.append("import (")
        # 标准库在前
        std_libs = [imp for imp in all_imports if not imp.startswith('.') and '/' not in imp]
        for lib in sorted(std_libs):
            merged_content.append(f'\t"{lib}"')
        
        # 外部库在后
        ext_libs = [imp for imp in all_imports if '/' in imp]
        if ext_libs:
            merged_content.append("")
            for lib in sorted(ext_libs):
                merged_content.append(f'\t"{lib}"')
        
        merged_content.append(")\n")
    
    # 3. Test函数
    for i, func in enumerate(all_functions):
        if i > 0:
            merged_content.append("")
        merged_content.append(func)
    
    # 写入文件
    output_path = full_path / output_file
    try:
        with open(output_path, 'w', encoding='utf-8') as f:
            f.write('\n'.join(merged_content))
        
        print_colored(GREEN, f"[SUCCESS] 合并成功: {output_file}")
        print_colored(BLUE, f"   📊 统计: {file_count} 个文件, {len(all_functions)} 个测试函数")
        
        return True
        
    except Exception as e:
        print_colored(RED, f"❌ 写入失败: {e}")
        return False

def backup_files(module_path, files_to_backup, backup_dir):
    """备份文件"""
    print_colored(YELLOW, f"\n💾 备份文件到: {backup_dir}")
    
    backup_path = Path(backup_dir)
    backup_path.mkdir(parents=True, exist_ok=True)
    
    for file_name in files_to_backup:
        src = Path(module_path) / file_name
        if src.exists():
            dst = backup_path / file_name
            shutil.copy2(src, dst)
            print_colored(GREEN, f"  [OK] 备份: {file_name}")

def main():
    """主函数"""
    print_colored(BLUE, "GameLink 测试文件合并工具")
    print_colored(BLUE, "=" * 50)
    
    # 检查是否在项目根目录
    if not Path("go.mod").exists():
        print_colored(RED, "❌ 错误: 请在项目根目录运行此脚本")
        sys.exit(1)
    
    # 配置
    modules = [
        {
            "path": "internal/service/commission",
            "files": [
                "commission_test.go",
                "commission_extended_test.go",
                "commission_additional_test.go"
            ],
            "output": "commission_test.go"
        },
        {
            "path": "internal/service/item",
            "files": [
                "item_test.go",
                "item_extended_test.go"
            ],
            "output": "item_test.go"
        },
        {
            "path": "internal/service/payment",
            "files": [
                "payment_test.go",
                "payment_extended_test.go",
                "payment_additional_test.go",
                "payment_full_coverage_test.go"
            ],
            "output": "payment_test.go"
        },
        {
            "path": "internal/service/order",
            "files": [
                "order_test.go",
                "order_extended_test.go",
                "order_autodestroy_test.go",
                "order_availability_test.go"
            ],
            "output": "order_test.go"
        }
    ]
    
    # 创建备份目录
    backup_base = Path("backup/merge_") / Path(str(Path.cwd().name))
    backup_base.mkdir(parents=True, exist_ok=True)
    
    success_count = 0
    
    for module in modules:
        print_colored(BLUE, "\n" + "=" * 50)
        
        # 创建模块备份
        module_backup = backup_base / module["path"].replace("/", "_")
        backup_files(module["path"], module["files"], module_backup)
        
        # 合并文件
        if merge_test_files(module["path"], module["files"], module["output"]):
            success_count += 1
            print_colored(GREEN, f"[SUCCESS] {module['path']} 处理成功")
        else:
            print_colored(RED, f"[ERROR] {module['path']} 处理失败")
    
    # 总结
    print_colored(BLUE, "\n" + "=" * 50)
    print_colored(BLUE, "📊 合并完成总结")
    print_colored(GREEN, f"✅ 成功: {success_count}/{len(modules)} 个模块")
    
    if success_count == len(modules):
        print_colored(GREEN, "\n[SUCCESS] 所有模块合并成功！")
        print_colored(YELLOW, "\n下一步操作:")
        print_colored(YELLOW, "1. 运行测试: go test ./internal/service/...")
        print_colored(YELLOW, "2. 如果测试通过，删除备份目录")
        print_colored(YELLOW, "3. 提交代码: git add . && git commit -m '合并Service层测试文件'")
    else:
        print_colored(RED, "\n[ERROR] 部分模块合并失败")
        print_colored(YELLOW, "请检查错误信息并修复问题")
        sys.exit(1)

if __name__ == "__main__":
    main()
