package content

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"github.com/xuri/excelize/v2"
)

// ContentStatsService 内容统计服务
type ContentStatsService struct {
	feedRepo    repository.FeedRepository
	messageRepo repository.ChatMessageRepository
}

// NewContentStatsService 创建内容统计服务
func NewContentStatsService(
	feedRepo repository.FeedRepository,
	messageRepo repository.ChatMessageRepository,
) *ContentStatsService {
	return &ContentStatsService{
		feedRepo:    feedRepo,
		messageRepo: messageRepo,
	}
}

// ContentStatsDTO 内容统计DTO
type ContentStatsDTO struct {
	TotalFeeds    int64                                `json:"totalFeeds"`
	PendingFeeds  int64                                `json:"pendingFeeds"`
	ApprovedFeeds int64                                `json:"approvedFeeds"`
	RejectedFeeds int64                                `json:"rejectedFeeds"`
	TotalMessages int64                                `json:"totalMessages"`
	FeedsByStatus map[model.FeedModerationStatus]int64 `json:"feedsByStatus"`
	FeedTrend     []repository.DateValue               `json:"feedTrend"`
}

// GetStats 获取内容统计
func (s *ContentStatsService) GetStats(ctx context.Context, trendDays int) (*ContentStatsDTO, error) {
	if trendDays < 1 {
		trendDays = 30
	}

	// 获取动态状态统计
	feedsByStatus, err := s.feedRepo.CountByStatus(ctx)
	if err != nil {
		return nil, err
	}

	// 计算总数
	var totalFeeds int64
	for _, count := range feedsByStatus {
		totalFeeds += count
	}

	// 获取动态趋势
	feedTrend, err := s.feedRepo.GetTrend(ctx, trendDays)
	if err != nil {
		feedTrend = []repository.DateValue{}
	}

	// 获取消息统计（简化版，可扩展）
	messages, totalMessages, _ := s.messageRepo.ListForModeration(ctx, repository.ChatMessageModerationListOptions{
		Page:     1,
		PageSize: 1,
	})
	_ = messages // 只需要总数

	return &ContentStatsDTO{
		TotalFeeds:    totalFeeds,
		PendingFeeds:  feedsByStatus[model.FeedModerationPending],
		ApprovedFeeds: feedsByStatus[model.FeedModerationApproved],
		RejectedFeeds: feedsByStatus[model.FeedModerationRejected],
		TotalMessages: totalMessages,
		FeedsByStatus: feedsByStatus,
		FeedTrend:     feedTrend,
	}, nil
}

// ExportStats 导出内容统计数据为Excel
func (s *ContentStatsService) ExportStats(ctx context.Context, trendDays int) (*bytes.Buffer, string, error) {
	stats, err := s.GetStats(ctx, trendDays)
	if err != nil {
		return nil, "", err
	}

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	// 创建概览工作表
	sheetName := "内容统计概览"
	index, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(index)

	// 设置标题样式
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	// 写入标题
	_ = f.SetCellValue(sheetName, "A1", "GameLink 内容管理统计报表")
	_ = f.MergeCell(sheetName, "A1", "D1")
	_ = f.SetCellStyle(sheetName, "A1", "D1", titleStyle)

	// 写入生成时间
	_ = f.SetCellValue(sheetName, "A2", fmt.Sprintf("生成时间: %s", time.Now().Format("2006-01-02 15:04:05")))
	_ = f.MergeCell(sheetName, "A2", "D2")

	// 写入统计概览
	_ = f.SetCellValue(sheetName, "A4", "统计项目")
	_ = f.SetCellValue(sheetName, "B4", "数量")
	_ = f.SetCellStyle(sheetName, "A4", "B4", headerStyle)

	_ = f.SetCellValue(sheetName, "A5", "动态总数")
	_ = f.SetCellValue(sheetName, "B5", stats.TotalFeeds)

	_ = f.SetCellValue(sheetName, "A6", "待审核动态")
	_ = f.SetCellValue(sheetName, "B6", stats.PendingFeeds)

	_ = f.SetCellValue(sheetName, "A7", "已通过动态")
	_ = f.SetCellValue(sheetName, "B7", stats.ApprovedFeeds)

	_ = f.SetCellValue(sheetName, "A8", "已拒绝动态")
	_ = f.SetCellValue(sheetName, "B8", stats.RejectedFeeds)

	_ = f.SetCellValue(sheetName, "A9", "聊天消息总数")
	_ = f.SetCellValue(sheetName, "B9", stats.TotalMessages)

	// 设置列宽
	_ = f.SetColWidth(sheetName, "A", "A", 20)
	_ = f.SetColWidth(sheetName, "B", "B", 15)

	// 创建趋势工作表
	trendSheet := "发布趋势"
	_, _ = f.NewSheet(trendSheet)

	_ = f.SetCellValue(trendSheet, "A1", fmt.Sprintf("最近%d天动态发布趋势", trendDays))
	_ = f.MergeCell(trendSheet, "A1", "B1")
	_ = f.SetCellStyle(trendSheet, "A1", "B1", titleStyle)

	_ = f.SetCellValue(trendSheet, "A3", "日期")
	_ = f.SetCellValue(trendSheet, "B3", "发布数量")
	_ = f.SetCellStyle(trendSheet, "A3", "B3", headerStyle)

	for i, trend := range stats.FeedTrend {
		row := i + 4
		_ = f.SetCellValue(trendSheet, fmt.Sprintf("A%d", row), trend.Date)
		_ = f.SetCellValue(trendSheet, fmt.Sprintf("B%d", row), trend.Value)
	}

	_ = f.SetColWidth(trendSheet, "A", "A", 15)
	_ = f.SetColWidth(trendSheet, "B", "B", 15)

	// 删除默认的Sheet1
	_ = f.DeleteSheet("Sheet1")

	// 写入buffer
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("content_stats_%s.xlsx", time.Now().Format("20060102_150405"))
	return buf, filename, nil
}
