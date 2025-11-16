#!/usr/bin/env python3
"""
Fix Swagger generic annotations in Go files.
Replaces model.APIResponse[any] and map[string]any with proper types.
"""

import os
import re
import sys

def fix_swagger_annotations(content):
    """Fix Swagger annotations in Go content."""

    # Replace failure responses
    content = re.sub(
        r'(@Failure\s+\d+\s+{object}\s+)model\.APIResponse\[any\]',
        r'\1model.ErrorResponse',
        content
    )

    content = re.sub(
        r'(@Failure\s+\d+\s+{object}\s+)map\[string\]any',
        r'\1model.ErrorResponse',
        content
    )

    # Replace success responses that use any
    content = re.sub(
        r'(@Success\s+\d+\s+{object}\s+)map\[string\]any',
        r'\1model.SuccessResponse',
        content
    )

    content = re.sub(
        r'(@Success\s+\d+\s+{object}\s+)model\.APIResponse\[any\]',
        r'\1model.SuccessResponse',
        content
    )

    return content

def process_file(file_path):
    """Process a single Go file."""
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()

        original_content = content
        content = fix_swagger_annotations(content)

        if content != original_content:
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(content)
            print(f"Fixed: {file_path}")
            return True
        else:
            print(f"No changes needed: {file_path}")
            return False

    except Exception as e:
        print(f"Error processing {file_path}: {e}")
        return False

def main():
    """Main function."""
    if len(sys.argv) != 2:
        print("Usage: python fix_swagger_annotations.py <directory>")
        sys.exit(1)

    root_dir = sys.argv[1]
    fixed_count = 0

    # Walk through all Go files in the directory
    for root, dirs, files in os.walk(root_dir):
        for file in files:
            if file.endswith('.go'):
                file_path = os.path.join(root, file)
                if process_file(file_path):
                    fixed_count += 1

    print(f"\nTotal files fixed: {fixed_count}")

if __name__ == "__main__":
    main()