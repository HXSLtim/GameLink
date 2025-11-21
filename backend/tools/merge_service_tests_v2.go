package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ImportInfo struct {
	alias string
	path  string
}

func main() {
	fmt.Println("GameLink Service层测试文件合并工具 v2")
	fmt.Println("=====================================")
	
	// 定义要合并的模块
	modules := []struct {
		path   string
		files  []string
		output string
	}{
		{
			path:   "internal/service/commission",
			files:  []string{"commission_test.go", "commission_extended_test.go", "commission_additional_test.go"},
			output: "commission_test.go",
		},
		{
			path:   "internal/service/item",
			files:  []string{"item_test.go", "item_extended_test.go"},
			output: "item_test.go",
		},
		{
			path:   "internal/service/payment",
			files:  []string{"payment_test.go", "payment_extended_test.go", "payment_additional_test.go", "payment_full_coverage_test.go"},
			output: "payment_test.go",
		},
		{
			path:   "internal/service/order",
			files:  []string{"order_test.go", "order_extended_test.go", "order_autodestroy_test.go", "order_availability_test.go"},
			output: "order_test.go",
		},
	}
	
	successCount := 0
	for _, module := range modules {
		fmt.Printf("\n处理模块: %s\n", module.path)
		if err := mergeModule(module.path, module.files, module.output); err != nil {
			fmt.Printf("失败: %v\n", err)
		} else {
			fmt.Printf("成功\n")
			successCount++
		}
	}
	
	fmt.Printf("\n=====================================\n")
	fmt.Printf("完成: %d/%d 个模块\n", successCount, len(modules))
	
	if successCount == len(modules) {
		fmt.Println("\n所有模块合并成功！")
		fmt.Println("\n下一步:")
		fmt.Println("1. 运行测试: go test ./internal/service/...")
		fmt.Println("2. 验证通过后删除旧文件")
		fmt.Println("3. 提交代码")
	}
}

func mergeModule(dir string, files []string, output string) error {
	// 检查目录
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("目录不存在: %s", dir)
	}
	
	// 收集所有内容
	var packageName string
	imports := make(map[string]ImportInfo) // path -> ImportInfo
	var allCode []string
	
	for _, file := range files {
		filePath := filepath.Join(dir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("  文件不存在: %s\n", file)
			continue
		}
		
		fmt.Printf("  读取: %s\n", file)
		
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("读取失败 %s: %w", file, err)
		}
		
		// 解析文件
		name, fileImports, codeBlocks, err := parseGoFile(string(content))
		if err != nil {
			return fmt.Errorf("解析失败 %s: %w", file, err)
		}
		
		if packageName == "" {
			packageName = name
		}
		
		for path, info := range fileImports {
			imports[path] = info
		}
		
		allCode = append(allCode, codeBlocks...)
	}
	
	if packageName == "" {
		return fmt.Errorf("未找到package声明")
	}
	
	if len(allCode) == 0 {
		return fmt.Errorf("未找到代码块")
	}
	
	// 生成合并内容
	var builder strings.Builder
	
	// Package声明
	builder.WriteString(fmt.Sprintf("package %s\n\n", packageName))
	
	// Import语句
	if len(imports) > 0 {
		builder.WriteString("import (\n")
		
		// 标准库
		var stdLibs, extLibs []ImportInfo
		for path, info := range imports {
			if strings.Contains(path, ".") {
				extLibs = append(extLibs, info)
			} else {
				stdLibs = append(stdLibs, info)
			}
		}
		
		// 按字母排序
		for _, info := range stdLibs {
			if info.alias != "" {
				builder.WriteString(fmt.Sprintf("\t%s \"%s\"\n", info.alias, info.path))
			} else {
				builder.WriteString(fmt.Sprintf("\t\"%s\"\n", info.path))
			}
		}
		
		if len(extLibs) > 0 && len(stdLibs) > 0 {
			builder.WriteString("\n")
		}
		
		for _, info := range extLibs {
			if info.alias != "" {
				builder.WriteString(fmt.Sprintf("\t%s \"%s\"\n", info.alias, info.path))
			} else {
				builder.WriteString(fmt.Sprintf("\t\"%s\"\n", info.path))
			}
		}
		
		builder.WriteString(")\n\n")
	}
	
	// 所有代码块（包括mock和test函数）
	for i, code := range allCode {
		if i > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(code)
		builder.WriteString("\n")
	}
	
	// 写入输出文件
	outputPath := filepath.Join(dir, output)
	if err := os.WriteFile(outputPath, []byte(builder.String()), 0644); err != nil {
		return fmt.Errorf("写入失败: %w", err)
	}
	
	fmt.Printf("  生成: %s (%d 个代码块)\n", output, len(allCode))
	
	return nil
}

func parseGoFile(content string) (packageName string, imports map[string]ImportInfo, allCode []string, err error) {
	imports = make(map[string]ImportInfo)
	
	reader := strings.NewReader(content)
	scanner := bufio.NewScanner(reader)
	
	var inImport bool
	var inCodeBlock bool
	var codeLines []string
	var braceCount int
	
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		
		// Package声明
		if strings.HasPrefix(line, "package ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				packageName = parts[1]
			}
			continue
		}
		
		// Import处理
		if strings.HasPrefix(line, "import (") {
			inImport = true
			continue
		}
		
		if inImport {
			if line == ")" {
				inImport = false
				continue
			}
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "//") {
				// 处理带别名的import，如：commissionrepo "gamelink/internal/repository/commission"
				if strings.Contains(line, " ") {
					parts := strings.Fields(line)
					if len(parts) >= 2 {
						alias := parts[0]
						path := strings.Trim(parts[len(parts)-1], `"`)
						imports[path] = ImportInfo{alias: alias, path: path}
					}
				} else {
					// 普通import
					path := strings.Trim(line, `"`)
					imports[path] = ImportInfo{alias: "", path: path}
				}
			}
			continue
		}
		
		// 跳过package和import后的空行和注释
		if !inCodeBlock && (trimmed == "" || strings.HasPrefix(trimmed, "//")) {
			continue
		}
		
		// 检测代码块开始
		if !inCodeBlock {
			// 任何结构体定义
			if strings.HasPrefix(trimmed, "type ") && strings.HasSuffix(trimmed, "struct {") {
				inCodeBlock = true
				codeLines = []string{line}
				braceCount = strings.Count(line, "{") - strings.Count(line, "}")
				continue
			}
			
			// 任何函数（包括mock方法和Test函数）
			if strings.HasPrefix(trimmed, "func ") {
				inCodeBlock = true
				codeLines = []string{line}
				braceCount = strings.Count(line, "{") - strings.Count(line, "}")
				continue
			}
		}
		
		// 处理代码块
		if inCodeBlock {
			codeLines = append(codeLines, line)
			braceCount += strings.Count(line, "{") - strings.Count(line, "}")
			
			if braceCount == 0 {
				allCode = append(allCode, strings.Join(codeLines, "\n"))
				inCodeBlock = false
			}
		}
	}
	
	if err := scanner.Err(); err != nil && err != io.EOF {
		return "", nil, nil, err
	}
	
	return packageName, imports, allCode, nil
}
