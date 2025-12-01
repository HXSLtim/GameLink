package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

/**
 * 测试骨架生成工具
 *
 * 使用方法:
 * go run scripts/generate_test_skeleton.go <package_path>
 *
 * 示例:
 * go run scripts/generate_test_skeleton.go internal/cache
 */

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run generate_test_skeleton.go <package_path>")
		fmt.Println("Example: go run generate_test_skeleton.go internal/cache")
		os.Exit(1)
	}

	packagePath := os.Args[1]

	// 遍历包中的所有Go文件
	files, err := filepath.Glob(filepath.Join(packagePath, "*.go"))
	if err != nil {
		fmt.Printf("Error reading package: %v\n", err)
		os.Exit(1)
	}

	for _, file := range files {
		// 跳过测试文件和生成的文件
		if strings.HasSuffix(file, "_test.go") ||
			strings.HasSuffix(file, "_gen.go") ||
			strings.Contains(file, "doc.go") {
			continue
		}

		// 检查是否已经有测试文件
		testFile := strings.TrimSuffix(file, ".go") + "_test.go"
		if _, err := os.Stat(testFile); err == nil {
			fmt.Printf("Test file already exists: %s\n", testFile)
			continue
		}

		// 解析源文件
		functions, packageName, err := parseSourceFile(file)
		if err != nil {
			fmt.Printf("Error parsing %s: %v\n", file, err)
			continue
		}

		if len(functions) == 0 {
			continue
		}

		// 生成测试文件
		testCode := generateTestFile(packageName, functions, filepath.Base(file))

		// 写入测试文件
		if err := os.WriteFile(testFile, []byte(testCode), 0644); err != nil {
			fmt.Printf("Error writing test file %s: %v\n", testFile, err)
			continue
		}

		fmt.Printf("Generated test file: %s (%d functions)\n", testFile, len(functions))
	}
}

func parseSourceFile(filePath string) ([]string, string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, "", err
	}

	var functions []string
	packageName := node.Name.Name

	// 遍历所有声明
	for _, decl := range node.Decls {
		// 只处理函数声明
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// 跳过私有函数和方法(可选)
		funcName := funcDecl.Name.Name
		if !ast.IsExported(funcName) {
			continue // 可以注释这行以包含私有函数
		}

		functions = append(functions, funcName)
	}

	return functions, packageName, nil
}

func generateTestFile(packageName string, functions []string, sourceFile string) string {
	var sb strings.Builder

	// 包声明和导入
	sb.WriteString(fmt.Sprintf("package %s\n\n", packageName))
	sb.WriteString("import (\n")
	sb.WriteString("\t\"testing\"\n\n")
	sb.WriteString("\t\"github.com/stretchr/testify/assert\"\n")
	sb.WriteString(")\n\n")

	sb.WriteString(fmt.Sprintf("// Generated test file for %s\n", sourceFile))
	sb.WriteString("// TODO: 实现具体的测试逻辑\n\n")

	// 为每个函数生成测试骨架
	for _, funcName := range functions {
		sb.WriteString(fmt.Sprintf("/**\n * Test_%s 测试 %s 函数\n */\n", funcName, funcName))
		sb.WriteString(fmt.Sprintf("func Test_%s(t *testing.T) {\n", funcName))
		sb.WriteString("\t// Arrange\n")
		sb.WriteString("\t// TODO: 设置测试数据和环境\n\n")
		sb.WriteString("\t// Act\n")
		sb.WriteString(fmt.Sprintf("\t// TODO: 调用 %s 函数\n\n", funcName))
		sb.WriteString("\t// Assert\n")
		sb.WriteString("\t// TODO: 验证结果\n")
		sb.WriteString("\tt.Skip(\"TODO: 实现测试逻辑\")\n")
		sb.WriteString("\tassert.NotNil(t, t)\n")
		sb.WriteString("}\n\n")

		// 额外生成边界条件测试
		sb.WriteString(fmt.Sprintf("/**\n * Test_%s_EdgeCases 测试 %s 的边界条件\n */\n", funcName, funcName))
		sb.WriteString(fmt.Sprintf("func Test_%s_EdgeCases(t *testing.T) {\n", funcName))
		sb.WriteString("\t// TODO: 测试边界条件\n")
		sb.WriteString("\tt.Skip(\"TODO: 实现边界条件测试\")\n")
		sb.WriteString("}\n\n")

		// 生成错误处理测试
		sb.WriteString(fmt.Sprintf("/**\n * Test_%s_ErrorHandling 测试 %s 的错误处理\n */\n", funcName, funcName))
		sb.WriteString(fmt.Sprintf("func Test_%s_ErrorHandling(t *testing.T) {\n", funcName))
		sb.WriteString("\t// TODO: 测试错误处理\n")
		sb.WriteString("\tt.Skip(\"TODO: 实现错误处理测试\")\n")
		sb.WriteString("}\n\n")
	}

	return sb.String()
}
