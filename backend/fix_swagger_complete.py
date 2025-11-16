#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Complete Swagger generic syntax fixer
Automatically creates Response types and updates annotations
"""

import os
import re
from pathlib import Path
from collections import defaultdict

def extract_generic_types(content):
    """Extract all unique generic types from Success annotations"""
    pattern = r'// @Success\s+\d+\s+\{object\}\s+model\.APIResponse\[([^\]]+)\]'
    matches = re.findall(pattern, content)
    return set(matches)

def create_response_type_name(generic_type):
    """Generate a response type name from generic type"""
    # Remove package prefixes and clean up
    parts = generic_type.split('.')
    if len(parts) > 1:
        type_name = parts[-1]
    else:
        type_name = generic_type

    # Handle special cases
    if type_name == 'any' or type_name == 'interface{}':
        return None
    if '[' in type_name:  # Nested generics
        return None

    # Generate response name
    if 'Response' in type_name:
        return type_name.replace('Response', '') + 'SuccessResponse'
    else:
        return type_name + 'Response'

def generate_response_struct(response_name, data_type):
    """Generate Go response struct definition"""
    return f"""// {response_name} API响应
type {response_name} struct {{
\tSuccess bool        `json:"success"`
\tCode    int         `json:"code"`
\tMessage string      `json:"message"`
\tData    {data_type} `json:"data"`
}}
"""

def fix_file_complete(filepath):
    """Fix file with complete response type generation"""
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()

        original_content = content

        # Extract all generic types used in this file
        generic_types = extract_generic_types(content)

        if not generic_types:
            return False

        # Generate response types
        response_types = {}
        for gen_type in generic_types:
            resp_name = create_response_type_name(gen_type)
            if resp_name:
                response_types[gen_type] = (resp_name, gen_type)

        if not response_types:
            return False

        # Find the package and import block
        package_match = re.search(r'(package\s+\w+\s*\n\s*import\s*\([^)]+\))', content, re.DOTALL)
        if not package_match:
            package_match = re.search(r'(package\s+\w+.*?)\n\n', content, re.DOTALL)

        if package_match:
            insert_pos = package_match.end()

            # Generate all response type definitions
            type_defs = '\n'
            for gen_type, (resp_name, data_type) in sorted(response_types.items()):
                type_defs += generate_response_struct(resp_name, data_type)

            # Insert type definitions after imports
            content = content[:insert_pos] + type_defs + content[insert_pos:]

            # Replace @Success annotations
            for gen_type, (resp_name, _) in response_types.items():
                escaped_type = re.escape(gen_type)
                pattern = r'(// @Success\s+)(\d+)(\s+\{object\}\s+)model\.APIResponse\[' + escaped_type + r'\]'
                replacement = r'\1\2\3' + resp_name
                content = re.sub(pattern, replacement, content)

        # Write back if changed
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
    print("Fixing all Swagger generic syntax...")

    # Find all .go files in internal/handler (exclude tests)
    handler_dir = Path('internal/handler')
    go_files = [
        f for f in handler_dir.rglob('*.go')
        if not f.name.endswith('_test.go')
    ]

    print(f"Found {len(go_files)} files to process")

    fixed_count = 0
    for filepath in go_files:
        if fix_file_complete(filepath):
            print(f"  [FIXED] {filepath}")
            fixed_count += 1

    print(f"\nCompleted! Fixed {fixed_count} files")
    print("\nRun: swag init -g cmd/main.go -o docs/swagger")

if __name__ == '__main__':
    main()
