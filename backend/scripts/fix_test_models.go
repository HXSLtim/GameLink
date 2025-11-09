//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	fixes := []struct {
		name        string
		pattern     *regexp.Regexp
		replacement string
	}{
		{
			name:        "PriceCents字段",
			pattern:     regexp.MustCompile(`PriceCents:\s*(\d+)`),
			replacement: `UnitPriceCents: $1, TotalPriceCents: $1`,
		},
		{
			name:        "PlayerID整数字面量",
			pattern:     regexp.MustCompile(`PlayerID:\s*(\d+)`),
			replacement: `PlayerID: ptrUint64($1)`,
		},
		{
			name:        "GameID整数字面量",
			pattern:     regexp.MustCompile(`GameID:\s*(\d+)`),
			replacement: `GameID: ptrUint64($1)`,
		},
	}

	var totalFixed int
	err := filepath.WalkDir("internal", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, "_test.go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		original := content
		fileFixed := 0

		for _, fix := range fixes {
			matches := fix.pattern.FindAll(content, -1)
			if len(matches) > 0 {
				content = fix.pattern.ReplaceAll(content, []byte(fix.replacement))
				fileFixed += len(matches)
				fmt.Printf("  [%s] 修复 %d 处: %s\n", filepath.Base(path), len(matches), fix.name)
			}
		}

		if !bytes.Equal(original, content) {
			// 确保文件包含辅助函数
			if bytes.Contains(content, []byte("ptrUint64")) && !bytes.Contains(content, []byte("func ptrUint64")) {
				helperFunc := []byte("\n// ptrUint64 返回 uint64 的指针\nfunc ptrUint64(v uint64) *uint64 { return &v }\n")
				// 在 package 声明后插入
				packageIdx := bytes.Index(content, []byte("\nimport"))
				if packageIdx > 0 {
					content = append(content[:packageIdx], append(helperFunc, content[packageIdx:]...)...)
				}
			}

			if err := os.WriteFile(path, content, 0644); err != nil {
				return err
			}
			totalFixed += fileFixed
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ 修复完成！共修复 %d 处\n", totalFixed)
}
