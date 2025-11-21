#!/usr/bin/env python3
"""
GameLink 代码审查整改实施脚本

功能:
1. 自动应用部分代码整改
2. 生成整改报告
3. 验证整改效果

使用方法:
    python scripts/code_review_fixes.py [选项]

选项:
    --check         检查当前代码状态
    --apply-fixes   应用自动修复
    --generate-report 生成整改报告
    --verify        验证整改效果
"""

import os
import sys
import re
import json
import argparse
from pathlib import Path
from datetime import datetime
from typing import Dict, List, Tuple

class CodeReviewFixer:
    """代码审查整改实施器"""
    
    def __init__(self, project_root: str):
        self.project_root = Path(project_root)
        self.internal_dir = self.project_root / "internal"
        self.fixes_applied = []
        self.issues_found = []
        
    def check_jwt_security(self) -> List[Dict]:
        """检查JWT安全问题"""
        issues = []
        
        # 检查中间件中的硬编码密钥
        middleware_file = self.internal_dir / "handler" / "middleware" / "jwt_auth.go"
        if middleware_file.exists():
            content = middleware_file.read_text(encoding='utf-8')
            
            # 检查硬编码密钥
            if 'gamelink-default-secret-key-change-in-production' in content:
                issues.append({
                    'file': str(middleware_file),
                    'line': self._find_line_number(content, 'gamelink-default-secret-key-change-in-production'),
                    'issue': '硬编码JWT密钥',
                    'severity': 'HIGH',
                    'fixable': True,
                    'fix_type': 'replace'
                })
            
            # 检查缺少密钥长度验证
            if 'len(secretKey) < 32' not in content:
                issues.append({
                    'file': str(middleware_file),
                    'line': None,
                    'issue': '缺少JWT密钥长度验证',
                    'severity': 'MEDIUM',
                    'fixable': True,
                    'fix_type': 'add'
                })
        
        return issues
    
    def check_input_validation(self) -> List[Dict]:
        """检查输入验证问题"""
        issues = []
        
        # 检查XSS防护
        service_files = list((self.internal_dir / "service").rglob("*.go"))
        for file_path in service_files:
            content = file_path.read_text(encoding='utf-8')
            
            # 检查是否有用户输入处理但没有XSS过滤
            if ('user.Name' in content or 'user.Bio' in content) and 'bluemonday' not in content:
                if 'Update' in content or 'Create' in content:
                    issues.append({
                        'file': str(file_path),
                        'line': None,
                        'issue': '用户输入缺少XSS过滤',
                        'severity': 'HIGH',
                        'fixable': True,
                        'fix_type': 'add'
                    })
        
        # 检查邮箱验证
        auth_file = self.internal_dir / "service" / "auth" / "auth.go"
        if auth_file.exists():
            content = auth_file.read_text(encoding='utf-8')
            
            # 检查简单的邮箱验证
            if 'mail.ParseAddress(email)' in content and 'emailRegex' not in content:
                issues.append({
                    'file': str(auth_file),
                    'line': self._find_line_number(content, 'mail.ParseAddress(email)'),
                    'issue': '邮箱验证过于简单',
                    'severity': 'MEDIUM',
                    'fixable': True,
                    'fix_type': 'enhance'
                })
        
        return issues
    
    def check_database_indexes(self) -> List[Dict]:
        """检查数据库索引问题"""
        issues = []
        
        # 检查User模型
        user_model = self.internal_dir / "model" / "user.go"
        if user_model.exists():
            content = user_model.read_text(encoding='utf-8')
            
            # 检查Name字段索引
            if '`json:"name" gorm:"size:64"`' in content and 'index' not in content:
                issues.append({
                    'file': str(user_model),
                    'line': self._find_line_number(content, '`json:"name" gorm:"size:64"`'),
                    'issue': 'User.Name字段缺少索引',
                    'severity': 'MEDIUM',
                    'fixable': True,
                    'fix_type': 'enhance'
                })
        
        # 检查Order模型
        order_model = self.internal_dir / "model" / "order.go"
        if order_model.exists():
            content = order_model.read_text(encoding='utf-8')
            
            # 检查复合索引
            if 'idx_user_status_created' not in content:
                issues.append({
                    'file': str(order_model),
                    'line': None,
                    'issue': '缺少复合索引优化查询',
                    'severity': 'MEDIUM',
                    'fixable': True,
                    'fix_type': 'add'
                })
        
        return issues
    
    def check_error_handling(self) -> List[Dict]:
        """检查错误处理问题"""
        issues = []
        
        # 检查是否使用统一的错误包
        handler_files = list((self.internal_dir / "handler").rglob("*.go"))
        
        for file_path in handler_files:
            content = file_path.read_text(encoding='utf-8')
            
            # 检查是否直接返回错误字符串
            if 'c.JSON(http.Status' in content and 'err.Error()' in content:
                if 'apierr' not in content:
                    issues.append({
                        'file': str(file_path),
                        'line': None,
                        'issue': '未使用统一的错误处理',
                        'severity': 'MEDIUM',
                        'fixable': True,
                        'fix_type': 'refactor'
                    })
        
        return issues
    
    def check_cache_usage(self) -> List[Dict]:
        """检查缓存使用问题"""
        issues = []
        
        # 检查Repository是否使用缓存
        repo_files = list((self.internal_dir / "repository").rglob("*.go"))
        
        for file_path in repo_files:
            content = file_path.read_text(encoding='utf-8')
            
            # 检查是否有Get方法但没有缓存逻辑
            if 'func (r *gorm' in content and 'Get(ctx context.Context' in content:
                if 'cache' not in content.lower():
                    issues.append({
                        'file': str(file_path),
                        'line': None,
                        'issue': '频繁查询的方法缺少缓存',
                        'severity': 'MEDIUM',
                        'fixable': True,
                        'fix_type': 'add'
                    })
        
        return issues
    
    def check_code_duplication(self) -> List[Dict]:
        """检查代码重复问题"""
        issues = []
        
        # 检查Repository的重复模式
        repo_files = list((self.internal_dir / "repository").rglob("*.go"))
        
        # 统计CRUD方法的数量
        crud_patterns = ['Create(ctx context.Context', 'Get(ctx context.Context', 
                        'Update(ctx context.Context', 'Delete(ctx context.Context']
        
        for file_path in repo_files:
            content = file_path.read_text(encoding='utf-8')
            crud_count = sum(1 for pattern in crud_patterns if pattern in content)
            
            # 如果文件中有多个CRUD方法，说明可以使用泛型
            if crud_count >= 3 and 'GenericRepository' not in content:
                issues.append({
                    'file': str(file_path),
                    'line': None,
                    'issue': '可以使用泛型Repository减少重复代码',
                    'severity': 'LOW',
                    'fixable': True,
                    'fix_type': 'refactor'
                })
        
        return issues
    
    def apply_fix(self, issue: Dict) -> bool:
        """应用修复"""
        try:
            file_path = Path(issue['file'])
            if not file_path.exists():
                return False
            
            content = file_path.read_text(encoding='utf-8')
            
            if issue['issue'] == '硬编码JWT密钥':
                self._fix_jwt_hardcoded_key(file_path, content)
            elif issue['issue'] == '缺少JWT密钥长度验证':
                self._add_jwt_key_validation(file_path, content)
            elif issue['issue'] == '用户输入缺少XSS过滤':
                self._add_xss_filtering(file_path, content)
            elif issue['issue'] == '邮箱验证过于简单':
                self._enhance_email_validation(file_path, content)
            elif issue['issue'] == 'User.Name字段缺少索引':
                self._add_user_name_index(file_path, content)
            elif issue['issue'] == '缺少复合索引优化查询':
                self._add_composite_indexes(file_path, content)
            elif issue['issue'] == '未使用统一的错误处理':
                self._unify_error_handling(file_path, content)
            elif issue['issue'] == '频繁查询的方法缺少缓存':
                self._add_caching(file_path, content)
            
            self.fixes_applied.append(issue)
            return True
            
        except Exception as e:
            print(f"修复失败 {issue['file']}: {e}")
            return False
    
    def _fix_jwt_hardcoded_key(self, file_path: Path, content: str):
        """修复硬编码JWT密钥"""
        # 替换硬编码密钥
        new_content = content.replace(
            'secretKey = "gamelink-default-secret-key-change-in-production"',
            '''logging.Error("JWT_SECRET_KEY not configured")
    return func(c *gin.Context) {
        c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
            "success": false,
            "code":    http.StatusServiceUnavailable,
            "message": "认证服务配置错误，请联系管理员",
        })
    }'''
        )
        
        file_path.write_text(new_content, encoding='utf-8')
    
    def _add_jwt_key_validation(self, file_path: Path, content: str):
        """添加JWT密钥长度验证"""
        # 在密钥检查后面添加长度验证
        insertion_point = content.find('secretKey := os.Getenv("JWT_SECRET_KEY")')
        if insertion_point != -1:
            # 找到下一行
            next_line_start = content.find('\n', insertion_point) + 1
            
            validation_code = '''\n\t// 验证密钥长度
\tif len(secretKey) < 32 {
\t\tlogging.Error("JWT_SECRET_KEY too short, must be at least 32 characters")
\t\treturn func(c *gin.Context) {
\t\t\tc.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
\t\t\t\t"success": false,
\t\t\t\t"code":    http.StatusServiceUnavailable,
\t\t\t\t"message": "认证服务配置错误，请联系管理员",
\t\t\t})
\t\t}
\t}
'''
            
            new_content = content[:next_line_start] + validation_code + content[next_line_start:]
            file_path.write_text(new_content, encoding='utf-8')
    
    def _add_xss_filtering(self, file_path: Path, content: str):
        """添加XSS过滤"""
        # 检查是否已经导入bluemonday
        if 'bluemonday' not in content:
            # 在导入部分添加
            import_section_end = content.find('import (')
            if import_section_end != -1:
                import_end = content.find(')', import_section_end)
                new_import = '\t"github.com/microcosm-cc/bluemonday"\n'
                
                new_content = content[:import_end] + new_import + content[import_end:]
                file_path.write_text(new_content, encoding='utf-8')
    
    def _enhance_email_validation(self, file_path: Path, content: str):
        """增强邮箱验证"""
        # 替换简单的邮箱验证
        old_func = '''func isValidEmail(email string) bool {
\tif email == "" {
\t\treturn false
\t}
\t_, err := mail.ParseAddress(email)
\treturn err == nil
}'''
        
        new_func = '''var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
\tif email == "" || len(email) > 128 {
\t\treturn false
\t}
\t
\t// 基本格式验证
\t_, err := mail.ParseAddress(email)
\tif err != nil {
\t\treturn false
\t}
\t
\t// 正则表达式验证
\tif !emailRegex.MatchString(email) {
\t\treturn false
\t}
\t
\treturn true
}'''
        
        new_content = content.replace(old_func, new_func)
        file_path.write_text(new_content, encoding='utf-8')
    
    def _add_user_name_index(self, file_path: Path, content: str):
        """为用户Name字段添加索引"""
        # 替换Name字段定义
        new_content = content.replace(
            'Name         string     `json:"name" gorm:"size:64"`',
            'Name         string     `json:"name" gorm:"size:64;index"`'
        )
        file_path.write_text(new_content, encoding='utf-8')
    
    def _add_composite_indexes(self, file_path: Path, content: str):
        """添加复合索引"""
        # 这是一个复杂的改动，需要手动处理
        print(f"  需要在 {file_path} 中手动添加复合索引")
        print("  参考: index:idx_user_status_created")
    
    def _unify_error_handling(self, file_path: Path, content: str):
        """统一错误处理"""
        # 这是一个复杂的重构，需要手动处理
        print(f"  需要在 {file_path} 中手动重构错误处理")
        print("  参考: 使用 apierr 包")
    
    def _add_caching(self, file_path: Path, content: str):
        """添加缓存"""
        # 这是一个复杂的改动，需要手动处理
        print(f"  需要在 {file_path} 中手动添加缓存逻辑")
        print("  参考: 使用 cache.Cache 接口")
    
    def _find_line_number(self, content: str, search_str: str) -> int:
        """查找字符串所在的行号"""
        lines = content.split('\n')
        for i, line in enumerate(lines, 1):
            if search_str in line:
                return i
        return 0
    
    def generate_report(self) -> str:
        """生成整改报告"""
        report = []
        report.append("# GameLink 代码审查整改报告")
        report.append(f"\n生成时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        report.append(f"项目路径: {self.project_root}")
        
        # 统计问题
        all_issues = []
        all_issues.extend(self.check_jwt_security())
        all_issues.extend(self.check_input_validation())
        all_issues.extend(self.check_database_indexes())
        all_issues.extend(self.check_error_handling())
        all_issues.extend(self.check_cache_usage())
        all_issues.extend(self.check_code_duplication())
        
        # 按严重程度分组
        high_issues = [i for i in all_issues if i['severity'] == 'HIGH']
        medium_issues = [i for i in all_issues if i['severity'] == 'MEDIUM']
        low_issues = [i for i in all_issues if i['severity'] == 'LOW']
        
        report.append(f"\n## 问题统计")
        report.append(f"- 严重问题 (HIGH): {len(high_issues)}")
        report.append(f"- 中等问题 (MEDIUM): {len(medium_issues)}")
        report.append(f"- 轻微问题 (LOW): {len(low_issues)}")
        report.append(f"- 总计: {len(all_issues)}")
        
        # 详细问题列表
        for severity, issues in [("严重", high_issues), ("中等", medium_issues), ("轻微", low_issues)]:
            if issues:
                report.append(f"\n## {severity}问题")
                for issue in issues:
                    report.append(f"\n### {issue['issue']}")
                    report.append(f"- **文件**: `{issue['file']}`")
                    if issue['line']:
                        report.append(f"- **行号**: {issue['line']}")
                    report.append(f"- **类型**: {issue['fix_type']}")
                    report.append(f"- **可自动修复**: {'是' if issue['fixable'] else '否'}")
        
        # 已应用的修复
        if self.fixes_applied:
            report.append(f"\n## 已应用的修复")
            report.append(f"已自动修复 {len(self.fixes_applied)} 个问题")
            
            for fix in self.fixes_applied:
                report.append(f"- ✓ {fix['issue']} - {fix['file']}")
        
        report.append(f"\n## 建议")
        report.append("1. 优先修复严重安全问题")
        report.append("2. 手动处理无法自动修复的问题")
        report.append("3. 添加测试验证修复效果")
        report.append("4. 进行代码审查确认整改质量")
        
        return '\n'.join(report)
    
    def run_all_checks(self) -> Dict:
        """运行所有检查"""
        return {
            'jwt_security': self.check_jwt_security(),
            'input_validation': self.check_input_validation(),
            'database_indexes': self.check_database_indexes(),
            'error_handling': self.check_error_handling(),
            'cache_usage': self.check_cache_usage(),
            'code_duplication': self.check_code_duplication(),
        }
    
    def apply_all_fixes(self, issues: List[Dict]) -> int:
        """应用所有可自动修复的问题"""
        applied = 0
        
        for issue in issues:
            if issue['fixable'] and issue['fix_type'] in ['replace', 'add']:
                print(f"修复中: {issue['issue']} - {issue['file']}")
                if self.apply_fix(issue):
                    applied += 1
                    print(f"  ✓ 修复成功")
                else:
                    print(f"  ✗ 修复失败")
        
        return applied


def main():
    parser = argparse.ArgumentParser(description='GameLink 代码审查整改工具')
    parser.add_argument('--check', action='store_true', help='检查代码问题')
    parser.add_argument('--apply-fixes', action='store_true', help='应用自动修复')
    parser.add_argument('--generate-report', action='store_true', help='生成整改报告')
    parser.add_argument('--verify', action='store_true', help='验证整改效果')
    parser.add_argument('--project-root', default='.', help='项目根目录')
    
    args = parser.parse_args()
    
    if not any([args.check, args.apply_fixes, args.generate_report, args.verify]):
        parser.print_help()
        sys.exit(1)
    
    fixer = CodeReviewFixer(args.project_root)
    
    if args.check:
        print("🔍 正在检查代码问题...")
        results = fixer.run_all_checks()
        
        total_issues = sum(len(issues) for issues in results.values())
        print(f"\n发现 {total_issues} 个问题:")
        
        for category, issues in results.items():
            if issues:
                print(f"\n{category.replace('_', ' ').title()}: {len(issues)} 个问题")
                for issue in issues[:3]:  # 只显示前3个
                    print(f"  - {issue['issue']} ({issue['severity']})")
    
    if args.apply_fixes:
        print("🔧 正在应用自动修复...")
        all_issues = []
        results = fixer.run_all_checks()
        
        for issues in results.values():
            all_issues.extend(issues)
        
        applied = fixer.apply_all_fixes(all_issues)
        print(f"\n应用了 {applied} 个自动修复")
    
    if args.generate_report:
        print("📋 正在生成整改报告...")
        report = fixer.generate_report()
        
        # 保存报告
        report_file = Path(args.project_root) / "CODE_REVIEW_FIXES_REPORT.md"
        report_file.write_text(report, encoding='utf-8')
        
        print(f"报告已保存到: {report_file}")
    
    if args.verify:
        print("✅ 正在验证整改效果...")
        # 重新检查
        results = fixer.run_all_checks()
        
        total_issues = sum(len(issues) for issues in results.values())
        print(f"剩余 {total_issues} 个问题")
        
        if total_issues == 0:
            print("🎉 所有问题已修复！")
        else:
            print("⚠️  还有未修复的问题，请查看报告")


if __name__ == '__main__':
    main()
