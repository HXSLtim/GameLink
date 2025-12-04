#!/usr/bin/env python3
"""
文件名重命名脚本
用于将backend目录下的文件按照Go规范重命名
"""

import os
import re
import sys
from pathlib import Path
from typing import Tuple, List


class FileRenamer:
    def __init__(self, root_dir: str, dry_run: bool = True):
        """
        初始化文件重命名器

        Args:
            root_dir: 根目录路径
            dry_run: 是否只预览不执行
        """
        self.root_dir = Path(root_dir)
        self.dry_run = dry_run
        self.renamed_files: List[Tuple[str, str]] = []
        self.skipped_files: List[str] = []

    def is_test_file(self, filename: str) -> bool:
        """
        判断是否为测试文件（以Test.go结尾）

        Args:
            filename: 文件名

        Returns:
            bool: 是否为测试文件
        """
        return filename.endswith('Test.go')

    def is_snake_case(self, filename: str) -> bool:
        """
        判断是否为snake_case命名

        Args:
            filename: 文件名

        Returns:
            bool: 是否为snake_case
        """
        # 去除后缀
        name_without_ext = filename[:-3] if filename.endswith('.go') else filename
        return '_' in name_without_ext and not self.is_camel_case(name_without_ext)

    def is_camel_case(self, name: str) -> bool:
        """
        判断是否为驼峰命名

        Args:
            name: 名称（不含后缀）

        Returns:
            bool: 是否为驼峰命名
        """
        # 如果包含下划线，不是驼峰命名
        if '_' in name:
            return False

        # 如果包含大小写混合，可能是驼峰命名
        has_lower = any(c.islower() for c in name)
        has_upper = any(c.isupper() for c in name)

        return has_lower and has_upper

    def convert_camel_to_snake(self, name: str) -> str:
        """
        将驼峰命名转换为snake_case

        Args:
            name: 驼峰命名字符串

        Returns:
            str: snake_case字符串
        """
        # 使用正则表达式进行转换
        # 在单词边界插入下划线
        s1 = re.sub('(.)([A-Z][a-z]+)', r'\1_\2', name)
        s2 = re.sub('([a-z0-9])([A-Z])', r'\1_\2', s1)
        return s2.lower()

    def rename_file(self, old_path: Path) -> str:
        """
        根据规则生成新的文件名

        Args:
            old_path: 原文件路径

        Returns:
            str: 新文件名
        """
        filename = old_path.name

        # 情况1: 小驼峰测试文件（如 userTest.go）
        if self.is_test_file(filename):
            # 移除Test后缀，添加_test前缀
            base_name = filename[:-6]  # 移除 'Test.go'
            new_name = f"{self.convert_camel_to_snake(base_name)}_test.go"
            return new_name

        # 情况2: snake_case文件（如 user_service.go）
        if self.is_snake_case(filename):
            # 保持不变
            return filename

        # 情况3: 标准小驼峰文件（如 userService.go）
        # 判断是否为驼峰命名
        name_without_ext = filename[:-3]
        if self.is_camel_case(name_without_ext):
            # 转换为snake_case
            new_name = f"{self.convert_camel_to_snake(name_without_ext)}.go"
            return new_name

        # 其他情况保持不变
        return filename

    def scan_directory(self) -> None:
        """
        扫描目录并收集需要重命名的文件
        """
        print(f"扫描目录: {self.root_dir}\n")

        for file_path in self.root_dir.rglob("*.go"):
            # 跳过 vendor 目录和隐藏目录
            if 'vendor' in file_path.parts or any(part.startswith('.') for part in file_path.parts):
                continue

            old_name = file_path.name
            new_name = self.rename_file(file_path)

            if old_name != new_name:
                new_path = file_path.parent / new_name
                self.renamed_files.append((str(file_path), str(new_path)))
            else:
                self.skipped_files.append(str(file_path))

    def print_summary(self) -> None:
        """
        打印重命名摘要
        """
        print("=" * 80)
        print("重命名摘要")
        print("=" * 80)

        if self.renamed_files:
            print(f"\n需要重命名的文件 ({len(self.renamed_files)} 个):")
            for old_path, new_path in self.renamed_files:
                old_name = Path(old_path).name
                new_name = Path(new_path).name
                print(f"  {old_name} → {new_name}")
                print(f"    路径: {old_path}")
        else:
            print("\n没有需要重命名的文件。")

        if self.skipped_files:
            print(f"\n\n跳过的文件 ({len(self.skipped_files)} 个):")
            for i, file_path in enumerate(self.skipped_files[:10], 1):
                print(f"  {i}. {Path(file_path).name}")

            if len(self.skipped_files) > 10:
                print(f"  ... 还有 {len(self.skipped_files) - 10} 个文件")

    def execute_renaming(self) -> None:
        """
        执行重命名操作
        """
        if self.dry_run:
            print("\n\n⚠️  当前为预览模式，未执行实际重命名操作")
            print("   使用 --apply 参数执行实际重命名")
            return

        print("\n\n" + "=" * 80)
        print("开始执行重命名操作...")
        print("=" * 80)

        success_count = 0
        fail_count = 0
        failed_files: List[Tuple[str, str, str]] = []

        for old_path, new_path in self.renamed_files:
            old_file = Path(old_path)
            new_file = Path(new_path)

            try:
                if new_file.exists():
                    raise FileExistsError(f"目标文件已存在: {new_path}")

                old_file.rename(new_file)
                print(f"✓ {old_file.name} → {new_file.name}")
                success_count += 1
            except Exception as e:
                print(f"✗ {old_file.name} → 失败: {str(e)}")
                fail_count += 1
                failed_files.append((old_path, new_path, str(e)))

        print("\n" + "=" * 80)
        print(f"执行完成: 成功 {success_count} 个, 失败 {fail_count} 个")
        print("=" * 80)

        if fail_count > 0:
            print("\n失败的文件:")
            for old_path, new_path, error in failed_files:
                print(f"  - {Path(old_path).name}")
                print(f"    错误: {error}")

def main():
    import argparse

    parser = argparse.ArgumentParser(
        description='批量重命名Go文件，按照规范转换为snake_case命名',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  # 预览模式（默认）
  python rename_files.py

  # 执行实际重命名
  python rename_files.py --apply

  # 指定目录
  python rename_files.py --dir /path/to/backend
        """
    )

    parser.add_argument(
        '--dir', '-d',
        default='backend',
        help='目标目录路径 (默认: backend)'
    )

    parser.add_argument(
        '--apply', '-a',
        action='store_true',
        help='执行实际重命名操作（默认只预览）'
    )

    args = parser.parse_args()

    if not os.path.isdir(args.dir):
        print(f"错误: 目录不存在: {args.dir}", file=sys.stderr)
        sys.exit(1)

    renamer = FileRenamer(args.dir, dry_run=not args.apply)
    renamer.scan_directory()
    renamer.print_summary()
    renamer.execute_renaming()

    if not args.apply:
        print("\n💡 提示: 使用 --apply 参数执行重命名")


if __name__ == "__main__":
    main()
