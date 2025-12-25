// Package logging provides generic logging utilities.
package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
)

var logger *slog.Logger

// ANSI颜色代码
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorGray    = "\033[37m"
	colorWhite   = "\033[97m"
)

// levelColors 定义不同日志级别的颜色
var levelColors = map[slog.Level]string{
	slog.LevelDebug: colorMagenta,
	slog.LevelInfo:  colorGreen,
	slog.LevelWarn:  colorYellow,
	slog.LevelError: colorRed,
}

// PrettyHandler 自定义的彩色表格格式Handler
type PrettyHandler struct {
	attrs       []slog.Attr
	groupStack  []string
	level       slog.Leveler
	output      io.Writer
	enableColor bool
	timeFormat  string
}

func init() {
	// 默认初始化
	Init("info")
}

// Init sets the default slog logger with colored table format handler and given level.
func Init(level string) *slog.Logger {
	lvl := parseLevel(level)
	enableColor := runtime.GOOS == "windows" || isTerminal(os.Stdout)

	if enableColor {
		// 使用彩色表格格式输出
		handler := NewPrettyHandler(os.Stdout, lvl, enableColor)
		logger = slog.New(handler)
	} else {
		// 非终端环境使用JSON格式（便于日志收集）
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl, AddSource: false})
		logger = slog.New(handler)
	}

	slog.SetDefault(logger)
	return logger
}

// NewPrettyHandler creates a new colored table format handler.
func NewPrettyHandler(output io.Writer, level slog.Leveler, enableColor bool) *PrettyHandler {
	return &PrettyHandler{
		output:      output,
		level:       level,
		enableColor: enableColor,
		timeFormat:  "2006-01-02 15:04:05.000",
	}
}

// Enabled implements slog.Handler interface.
func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

// Handle implements slog.Handler interface.
func (h *PrettyHandler) Handle(ctx context.Context, record slog.Record) error {
	buf := &bytes.Buffer{}

	// 时间
	timeStr := record.Time.Format(h.timeFormat)
	if h.enableColor {
		buf.WriteString(colorGray + "[" + timeStr + "]" + colorReset + " ")
	} else {
		buf.WriteString("[" + timeStr + "] ")
	}

	// 日志级别
	levelStr := strings.ToUpper(record.Level.String())
	if h.enableColor {
		color := levelColors[record.Level]
		buf.WriteString(color)
		buf.WriteString(fmt.Sprintf("%-5s", levelStr))
		buf.WriteString(colorReset)
	} else {
		buf.WriteString(fmt.Sprintf("%-5s", levelStr))
	}

	// 消息
	buf.WriteString(" | ")
	if h.enableColor {
		buf.WriteString(colorWhite)
		buf.WriteString(record.Message)
		buf.WriteString(colorReset)
	} else {
		buf.WriteString(record.Message)
	}

	// 处理属性
	attrsMap := make(map[string]interface{})

	// 添加handler级别的属性
	for _, attr := range h.attrs {
		attrsMap[attr.Key] = attr.Value.Any()
	}

	// 添加record级别的属性
	record.Attrs(func(attr slog.Attr) bool {
		attrsMap[attr.Key] = attr.Value.Any()
		return true
	})

	// 如果属性存在，输出它们
	if len(attrsMap) > 0 {
		buf.WriteString("\n")

		// 获取所有key并排序
		var keys []string
		for k := range attrsMap {
			keys = append(keys, k)
		}

		// 特殊处理http_request，尝试对齐输出
		if record.Message == "http_request" {
			httpAttrs := h.formatHTTPRequest(attrsMap)
			buf.WriteString(httpAttrs)
		} else {
			// 普通属性按JSON格式显示
			attrsJSON, _ := json.MarshalIndent(attrsMap, "", "  ")
			if h.enableColor && len(attrsJSON) > 0 {
				// 高亮JSON中的键
				jsonStr := string(attrsJSON)
				jsonStr = strings.ReplaceAll(jsonStr, `\"`, `"`)
				lines := strings.Split(jsonStr, "\n")
				for _, line := range lines {
					if strings.Contains(line, `":`) {
						parts := strings.SplitN(line, `":`, 2)
						buf.WriteString(colorCyan + parts[0] + `":` + colorReset)
						if len(parts) > 1 {
							buf.WriteString(parts[1])
						}
					} else {
						buf.WriteString(line)
					}
					buf.WriteString("\n")
				}
			} else {
				buf.WriteString("  " + string(attrsJSON))
			}
		}
	}
	buf.WriteString("\n")

	_, err := h.output.Write(buf.Bytes())
	return err
}

// httpFieldDef 定义HTTP日志字段
type httpFieldDef struct {
	key         string
	label       string
	colorizeVal func(h *PrettyHandler, v string) string
}

// getHTTPFields 返回HTTP日志字段定义
func getHTTPFields() []httpFieldDef {
	return []httpFieldDef{
		{"method", "HTTP方法", (*PrettyHandler).colorizeMethod},
		{"status", "状态码", (*PrettyHandler).colorizeStatus},
		{"path", "请求路径", func(h *PrettyHandler, v string) string { return v }},
		{"ip", "客户端IP", func(h *PrettyHandler, v string) string { return v }},
		{"duration", "响应时间", func(h *PrettyHandler, v string) string { return v }},
		{"request_id", "请求ID", func(h *PrettyHandler, v string) string { return v }},
		{"user_id", "用户ID", func(h *PrettyHandler, v string) string { return v }},
	}
}

func (h *PrettyHandler) colorizeMethod(v string) string {
	if !h.enableColor {
		return v
	}
	switch v {
	case "GET":
		return colorBlue + v + colorReset
	case "POST":
		return colorGreen + v + colorReset
	case "PUT":
		return colorYellow + v + colorReset
	case "DELETE":
		return colorRed + v + colorReset
	default:
		return v
	}
}

func (h *PrettyHandler) colorizeStatus(v string) string {
	if !h.enableColor {
		return v
	}
	status := 0
	fmt.Sscanf(v, "%d", &status)
	switch {
	case status >= 500:
		return colorRed + v + colorReset
	case status >= 400:
		return colorYellow + v + colorReset
	case status >= 300:
		return colorCyan + v + colorReset
	case status >= 200:
		return colorGreen + v + colorReset
	default:
		return v
	}
}

// formatHTTPRequest 格式化HTTP请求日志为表格样式
func (h *PrettyHandler) formatHTTPRequest(attrsMap map[string]interface{}) string {
	var buf bytes.Buffer
	col2Width := 60

	h.writeTableHeader(&buf, col2Width)
	plainAttrs := h.writeKnownFields(&buf, attrsMap, col2Width)
	h.writeExtraFields(&buf, plainAttrs, col2Width)
	h.writeTableFooter(&buf, col2Width)

	return buf.String()
}

func (h *PrettyHandler) writeTableHeader(buf *bytes.Buffer, col2Width int) {
	if h.enableColor {
		buf.WriteString(colorGray + "┌──────────────┬" + strings.Repeat("─", col2Width+2) + "┐\n" + colorReset)
	} else {
		buf.WriteString("┌──────────────┬" + strings.Repeat("─", col2Width+2) + "┐\n")
	}
}

func (h *PrettyHandler) writeTableFooter(buf *bytes.Buffer, col2Width int) {
	if h.enableColor {
		buf.WriteString(colorGray + "└──────────────┴" + strings.Repeat("─", col2Width+2) + "┘" + colorReset)
	} else {
		buf.WriteString("└──────────────┴" + strings.Repeat("─", col2Width+2) + "┘")
	}
}

func (h *PrettyHandler) writeKnownFields(buf *bytes.Buffer, attrsMap map[string]interface{}, col2Width int) map[string]string {
	fields := getHTTPFields()
	knownKeys := make(map[string]bool)
	for _, f := range fields {
		knownKeys[f.key] = true
		if val, ok := attrsMap[f.key]; ok {
			vStr := fmt.Sprintf("%v", val)
			coloredValue := f.colorizeVal(h, vStr)
			h.writeRow(buf, f.label, coloredValue, colorCyan, col2Width)
		}
	}

	plainAttrs := make(map[string]string)
	for key, val := range attrsMap {
		if !knownKeys[key] {
			plainAttrs[key] = fmt.Sprintf("%v", val)
		}
	}
	return plainAttrs
}

func (h *PrettyHandler) writeExtraFields(buf *bytes.Buffer, plainAttrs map[string]string, col2Width int) {
	if len(plainAttrs) == 0 {
		return
	}
	if h.enableColor {
		buf.WriteString(colorGray + "├──────────────┼" + strings.Repeat("─", col2Width+2) + "┤\n" + colorReset)
	} else {
		buf.WriteString("├──────────────┼" + strings.Repeat("─", col2Width+2) + "┤\n")
	}
	for key, val := range plainAttrs {
		h.writeRow(buf, key, val, colorWhite, col2Width)
	}
}

func (h *PrettyHandler) writeRow(buf *bytes.Buffer, key, value, colorCode string, col2Width int) {
	col1Width := 12
	if h.enableColor {
		buf.WriteString(colorGray + "│ " + colorReset)
		buf.WriteString(colorCode)
		buf.WriteString(fmt.Sprintf("%-*s", col1Width, key))
		buf.WriteString(colorGray + " │ " + colorReset)
		buf.WriteString(value)
		buf.WriteString(colorGray)
		if remaining := col2Width - len(value); remaining > 0 {
			buf.WriteString(strings.Repeat(" ", remaining))
		}
		buf.WriteString(" │" + colorReset + "\n")
	} else {
		buf.WriteString(fmt.Sprintf("│ %-12s │ %s", key, value))
		buf.WriteString(strings.Repeat(" ", col2Width-len(value)))
		buf.WriteString(" │\n")
	}
}

// WithAttrs implements slog.Handler interface.
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		attrs:       append(h.attrs, attrs...),
		output:      h.output,
		level:       h.level,
		enableColor: h.enableColor,
		timeFormat:  h.timeFormat,
	}
}

// WithGroup implements slog.Handler interface.
func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	// 简单处理group（这里不做特殊处理，因为表格格式下group意义不大）
	return h
}

// parseLevel 解析日志级别
func parseLevel(v string) slog.Leveler {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// isTerminal 检查文件描述符是否为终端（Windows下总是返回true）
func isTerminal(f *os.File) bool {
	// 在Windows环境下，powerShell和CMD都支持ANSI颜色码
	return runtime.GOOS == "windows"
}

// Debug logs a debug message with the given key-value pairs.
func Debug(msg string, args ...interface{}) {
	logger.Debug(msg, args...)
}

// Info logs an info message with the given key-value pairs.
func Info(msg string, args ...interface{}) {
	logger.Info(msg, args...)
}

// Warn logs a warning message with the given key-value pairs.
func Warn(msg string, args ...interface{}) {
	logger.Warn(msg, args...)
}

// Error logs an error message with the given key-value pairs.
func Error(msg string, args ...interface{}) {
	logger.Error(msg, args...)
}

// DebugContext logs a debug message with context and the given key-value pairs.
func DebugContext(ctx context.Context, msg string, args ...interface{}) {
	logger.DebugContext(ctx, msg, args...)
}

// InfoContext logs an info message with context and the given key-value pairs.
func InfoContext(ctx context.Context, msg string, args ...interface{}) {
	logger.InfoContext(ctx, msg, args...)
}

// WarnContext logs a warning message with context and the given key-value pairs.
func WarnContext(ctx context.Context, msg string, args ...interface{}) {
	logger.WarnContext(ctx, msg, args...)
}

// ErrorContext logs an error message with context and the given key-value pairs.
func ErrorContext(ctx context.Context, msg string, args ...interface{}) {
	logger.ErrorContext(ctx, msg, args...)
}
