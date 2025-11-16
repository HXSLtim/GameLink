#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Fix all UTF-8 encoding issues in Swagger comments
Replace any broken text with English placeholders
"""

import re
from pathlib import Path

def fix_broken_param(line):
    """Fix broken @Param lines"""
    # If line contains broken UTF-8 marker
    if '�' not in line:
        return line

    # Extract parameter name
    match = re.search(r'// @Param\s+(\w+)', line)
    if not match:
        return line

    param_name = match.group(1)

    # Common parameter mappings
    param_descriptions = {
        'dateFrom': 'Start date (YYYY-MM-DD)',
        'dateTo': 'End date (YYYY-MM-DD)',
        'months': 'Number of months (default 12)',
        'subCategory': 'Sub category',
        'isActive': 'Is active',
        'status': 'Status filter',
        'role': 'Role filter',
        'fields': 'Export fields (comma separated)',
        'month': 'Month filter (YYYY-MM)',
        'request': 'Request body',
    }

    description = param_descriptions.get(param_name, f'Parameter: {param_name}')

    # Rebuild the line
    # Extract the parameter definition parts
    parts_match = re.search(r'// @Param\s+(\w+)\s+(\w+)\s+([\w\[\]]+)\s+(true|false)', line)
    if parts_match:
        name, location, type_str, required = parts_match.groups()
        return f'// @Param        {name:<14} {location:<8} {type_str:<12} {required:<6} "{description}"'

    return line

def fix_broken_description(line):
    """Fix broken @Description lines"""
    if '�' not in line:
        return line

    # Just use a generic description
    if '@Description' in line:
        return '// @Description  API endpoint'

    return line

def fix_file(filepath):
    """Fix encoding issues in file"""
    try:
        with open(filepath, 'r', encoding='utf-8', errors='ignore') as f:
            lines = f.readlines()

        modified = False
        new_lines = []

        for line in lines:
            original = line

            # Fix @Param lines
            if '// @Param' in line and '�' in line:
                line = fix_broken_param(line)
                if line != original:
                    modified = True

            # Fix @Description lines
            elif '// @Description' in line and '�' in line:
                line = fix_broken_description(line)
                if line != original:
                    modified = True

            new_lines.append(line)

        if modified:
            with open(filepath, 'w', encoding='utf-8', newline='\n') as f:
                f.writelines(new_lines)
            return True

        return False
    except Exception as e:
        print(f"Error processing {filepath}: {e}")
        return False

def main():
    """Main function"""
    print("Fixing all UTF-8 encoding issues...")

    handler_dir = Path('internal/handler')
    go_files = list(handler_dir.rglob('*.go'))

    print(f"Checking {len(go_files)} files")

    fixed_count = 0
    for filepath in go_files:
        if fix_file(filepath):
            print(f"  [FIXED] {filepath}")
            fixed_count += 1

    print(f"\nCompleted! Fixed {fixed_count} files")

if __name__ == '__main__':
    main()
