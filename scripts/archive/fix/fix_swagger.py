#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Fix Swagger generic syntax in Go files
Replaces model.APIResponse[T] with concrete types
"""

import os
import re
from pathlib import Path

def fix_file(filepath):
    """Fix Swagger annotations in a single file"""
    try:
        # Read file with UTF-8 encoding
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()

        original_content = content

        # Replace Failure annotations with model.ErrorResponse
        content = re.sub(
            r'(// @Failure\s+)(\d+)(\s+\{object\}\s+)model\.APIResponse\[any\]',
            r'\1\2            {object}  model.ErrorResponse',
            content
        )
        content = re.sub(
            r'(// @Failure\s+)(\d+)(\s+\{object\}\s+)model\.APIResponse\[interface\{\}\]',
            r'\1\2            {object}  model.ErrorResponse',
            content
        )

        # Replace Success annotations with model.SuccessResponse (for generic any type)
        content = re.sub(
            r'(// @Success\s+)(\d+)(\s+\{object\}\s+)model\.APIResponse\[any\]',
            r'\1\2            {object}  model.SuccessResponse',
            content
        )

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
    print("Fixing Swagger generic syntax...")

    # Find all .go files in internal/handler (exclude tests)
    handler_dir = Path('internal/handler')
    go_files = [
        f for f in handler_dir.rglob('*.go')
        if not f.name.endswith('_test.go')
    ]

    print(f"Found {len(go_files)} files to process")

    fixed_count = 0
    for filepath in go_files:
        if fix_file(filepath):
            print(f"  [FIXED] {filepath}")
            fixed_count += 1

    print(f"\nCompleted! Fixed {fixed_count} files")
    print("\nNext steps:")
    print("  1. Run: swag init -g cmd/main.go -o docs/swagger")
    print("  2. Check remaining errors and create specific Response types as needed")
    print("\nSee SWAGGER_FIX_GUIDE.md for detailed instructions")

if __name__ == '__main__':
    main()
