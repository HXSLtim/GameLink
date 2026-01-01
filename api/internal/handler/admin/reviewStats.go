package admin

import (
	"bytes"
	"encoding/csv"
	"fmt"
	_ "gamelink/internal/model" // Imported for Swagger annotations
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/repository"
	reviewservice "gamelink/internal/service/review"
	"gamelink/pkg/apierr"
)

// ReviewStatsHandler 评价统计接口
type ReviewStatsHandler struct {
	svc *reviewservice.ReviewStatsService
}

// NewReviewStatsHandler 创建评价统计Handler
func NewReviewStatsHandler(svc *reviewservice.ReviewStatsService) *ReviewStatsHandler {
	return &ReviewStatsHandler{svc: svc}
}

// GetReviewStats 获取评价统计概览
// @Summary      获取评价统计概览
// @Description  获取总评价数、平均评分、各评分段分布
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  model.APIResponse[reviewservice.GetReviewStatsResponse]
// @Failure      500  {object}  model.ErrorResponse
// @Router       /admin/reviews/stats [get]
func (h *ReviewStatsHandler) GetReviewStats(c *gin.Context) {
	stats, err := h.svc.GetReviewStats(c.Request.Context())
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}
	respondSuccess(c, stats)
}

// GetReviewTrend 获取评价趋势
// @Summary      获取评价趋势
// @Description  获取最近N天的评价数量趋势
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Produce      json
// @Param        days  query  int  false  "统计天数（默认30天）"
// @Success      200  {object}  model.APIResponse[reviewservice.GetReviewTrendResponse]
// @Failure      500  {object}  model.ErrorResponse
// @Router       /admin/reviews/trend [get]
func (h *ReviewStatsHandler) GetReviewTrend(c *gin.Context) {
	days := 30
	if daysStr := strings.TrimSpace(c.Query("days")); daysStr != "" {
		if val, err := strconv.Atoi(daysStr); err == nil && val > 0 {
			days = val
		}
	}

	trend, err := h.svc.GetReviewTrend(c.Request.Context(), days)
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}
	respondSuccess(c, trend)
}

// GetTopPlayers 获取陪玩师排行榜
// @Summary      获取陪玩师排行榜
// @Description  获取评价最多或评分最高的陪玩师排行榜
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Produce      json
// @Param        limit    query  int     false  "数量限制（默认10）"
// @Param        sort_by  query  string  false  "排序方式：count（评价数量）或 rating（评分）"
// @Success      200  {object}  model.APIResponse[reviewservice.GetTopPlayersResponse]
// @Failure      500  {object}  model.ErrorResponse
// @Router       /admin/reviews/top-players [get]
func (h *ReviewStatsHandler) GetTopPlayers(c *gin.Context) {
	limit := 10
	if limitStr := strings.TrimSpace(c.Query("limit")); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	sortBy := strings.TrimSpace(c.Query("sort_by"))
	if sortBy == "" {
		sortBy = "count"
	}

	players, err := h.svc.GetTopPlayers(c.Request.Context(), limit, sortBy)
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}
	respondSuccess(c, players)
}

// GetGameStats 获取游戏统计
// @Summary      获取游戏统计
// @Description  获取各游戏的评价数量和平均评分
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  model.APIResponse[reviewservice.GetGameStatsResponse]
// @Failure      500  {object}  model.ErrorResponse
// @Router       /admin/reviews/game-stats [get]
func (h *ReviewStatsHandler) GetGameStats(c *gin.Context) {
	games, err := h.svc.GetGameStats(c.Request.Context())
	if err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}
	respondSuccess(c, games)
}

// ExportReviewStats 导出评价统计数据
// @Summary      导出评价统计数据
// @Description  导出评价统计数据为CSV格式
// @Tags         Admin/Reviews
// @Security     BearerAuth
// @Produce      text/csv
// @Param        type  query  string  false  "导出类型：overview（概览）、trend（趋势）、players（陪玩师排行）、games（游戏统计）"
// @Param        days  query  int     false  "趋势统计天数（默认30天）"
// @Param        limit query  int     false  "排行榜数量限制（默认10）"
// @Success      200  {file}  file  "CSV文件"
// @Failure      500  {object}  model.ErrorResponse
// @Router       /admin/reviews/export [get]
func (h *ReviewStatsHandler) ExportReviewStats(c *gin.Context) {
	exportType := strings.TrimSpace(c.Query("type"))
	if exportType == "" {
		exportType = "overview"
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	var filename string

	switch exportType {
	case "overview":
		filename = "review_stats_overview.csv"
		if err := h.exportOverview(c, writer); err != nil {
			respondError(c, apierr.InternalError(err.Error()))
			return
		}
	case "trend":
		filename = "review_stats_trend.csv"
		if err := h.exportTrend(c, writer); err != nil {
			respondError(c, apierr.InternalError(err.Error()))
			return
		}
	case "players":
		filename = "review_stats_players.csv"
		if err := h.exportPlayers(c, writer); err != nil {
			respondError(c, apierr.InternalError(err.Error()))
			return
		}
	case "games":
		filename = "review_stats_games.csv"
		if err := h.exportGames(c, writer); err != nil {
			respondError(c, apierr.InternalError(err.Error()))
			return
		}
	default:
		respondBadRequest(c, "invalid export type, must be one of: overview, trend, players, games")
		return
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		respondError(c, apierr.InternalError(err.Error()))
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

func (h *ReviewStatsHandler) exportOverview(c *gin.Context, writer *csv.Writer) error {
	stats, err := h.svc.GetReviewStats(c.Request.Context())
	if err != nil {
		return err
	}

	if err := writer.Write([]string{"指标", "值"}); err != nil {
		return err
	}

	rows := [][]string{
		{"总评价数", strconv.FormatInt(stats.TotalReviews, 10)},
		{"平均评分", fmt.Sprintf("%.2f", stats.AverageRating)},
	}

	for score := 1; score <= 5; score++ {
		count := stats.RatingDistribution[score]
		rows = append(rows, []string{fmt.Sprintf("%d星评价数", score), strconv.FormatInt(count, 10)})
	}

	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func (h *ReviewStatsHandler) exportTrend(c *gin.Context, writer *csv.Writer) error {
	days := 30
	if daysStr := strings.TrimSpace(c.Query("days")); daysStr != "" {
		if val, err := strconv.Atoi(daysStr); err == nil && val > 0 {
			days = val
		}
	}

	trend, err := h.svc.GetReviewTrend(c.Request.Context(), days)
	if err != nil {
		return err
	}

	if err := writer.Write([]string{"日期", "评价数量"}); err != nil {
		return err
	}

	for _, item := range trend.Trend {
		if err := writer.Write([]string{item.Date, strconv.FormatInt(item.Value, 10)}); err != nil {
			return err
		}
	}

	return nil
}

func (h *ReviewStatsHandler) exportPlayers(c *gin.Context, writer *csv.Writer) error {
	limit := 10
	if limitStr := strings.TrimSpace(c.Query("limit")); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	sortBy := strings.TrimSpace(c.Query("sort_by"))
	if sortBy == "" {
		sortBy = "count"
	}

	players, err := h.svc.GetTopPlayers(c.Request.Context(), limit, sortBy)
	if err != nil {
		return err
	}

	if err := writer.Write([]string{"排名", "陪玩师ID", "陪玩师名称", "评价数量", "平均评分"}); err != nil {
		return err
	}

	for i, player := range players.Players {
		row := []string{
			strconv.Itoa(i + 1),
			strconv.FormatUint(player.PlayerID, 10),
			player.PlayerName,
			strconv.FormatInt(player.ReviewCount, 10),
			fmt.Sprintf("%.2f", player.AverageRating),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func (h *ReviewStatsHandler) exportGames(c *gin.Context, writer *csv.Writer) error {
	games, err := h.svc.GetGameStats(c.Request.Context())
	if err != nil {
		return err
	}

	if err := writer.Write([]string{"游戏ID", "游戏名称", "评价数量", "平均评分"}); err != nil {
		return err
	}

	for _, game := range games.Games {
		row := []string{
			strconv.FormatUint(game.GameID, 10),
			game.GameName,
			strconv.FormatInt(game.ReviewCount, 10),
			fmt.Sprintf("%.2f", game.AverageRating),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// ReviewStatsResponse types for Swagger documentation
type ReviewStatsOverviewResponse = reviewservice.GetReviewStatsResponse
type ReviewStatsTrendResponse = reviewservice.GetReviewTrendResponse
type ReviewStatsTopPlayersResponse = reviewservice.GetTopPlayersResponse
type ReviewStatsGameStatsResponse = reviewservice.GetGameStatsResponse
type ReviewStatsDateValue = repository.DateValue
type ReviewStatsPlayerStats = repository.PlayerReviewStats
type ReviewStatsGameStats = repository.GameReviewStats
