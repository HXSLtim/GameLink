#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Fix UTF-8 encoding issues in Swagger comments
Replace broken Chinese with English descriptions
"""

import re
from pathlib import Path

# Common replacements for broken UTF-8 in Swagger comments
REPLACEMENTS = {
    r'月数[^"]*"': 'Number of months"',
    r'开始时[^"]*"': 'Start date"',
    r'结束时[^"]*"': 'End date"',
    r'游戏ID[^"]*"': 'Game ID"',
    r'用户ID[^"]*"': 'User ID"',
    r'订单ID[^"]*"': 'Order ID"',
    r'页码[^"]*"': 'Page number"',
    r'每页数[^"]*"': 'Page size"',
    r'关键[^"]*"': 'Keyword"',
    r'状态[^"]*"': 'Status"',
    r'类型[^"]*"': 'Type"',
}

def fix_file(filepath):
    """Fix UTF-8 issues in file"""
    try:
        with open(filepath, 'r', encoding='utf-8', errors='ignore') as f:
            content = f.read()

        original_content = content

        # Fix @Param comments with broken UTF-8
        for pattern, replacement in REPLACEMENTS.items():
            content = re.sub(pattern, replacement, content)

        # Fix any remaining broken param comments
        # Match: // @Param ... false/true  "broken_text�?[more_broken]"
        def fix_param_comment(match):
            param_def = match.group(1)
            # If the description looks broken (contains �), use generic description
            if '�' in param_def or len(param_def) > 100:
                parts = param_def.split()
                if len(parts) >= 4:
                    param_name = parts[0]
                    return f'// @Param        {param_name:<14} {" ".join(parts[1:4]):<25} "Parameter: {param_name}"'
            return match.group(0)

        content = re.sub(r'// @Param\s+(.*?�[^"\n]*"[^"\n]*)', fix_param_comment, content)

        if content != original_content:
            with open(filepath, 'w', encoding='utf-8', newline='\n') as f:
                f.write(content)
            return True
        return False
    except Exception as e:
        print(f"Error processing {filepath}: {e}")
        return False

def main():
    """Main function"""
    print("Fixing UTF-8 encoding issues in Swagger comments...")

    handler_dir = Path('internal/handler')
    go_files = [
        f for f in handler_dir.rglob('*.go')
        if not f.name.endswith('_test.go')
    ]

    print(f"Checking {len(go_files)} files")

    fixed_count = 0
    for filepath in go_files:
        if fix_file(filepath):
            print(f"  [FIXED] {filepath}")
            fixed_count += 1

    print(f"\nCompleted! Fixed {fixed_count} files")
    print("\nRun: swag init -g cmd/main.go -o docs/swagger")

if __name__ == '__main__':
    main()
