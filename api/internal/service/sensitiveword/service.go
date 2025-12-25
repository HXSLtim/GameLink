package sensitiveword

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

var (
	// ErrNotFound 敏感词不存在
	ErrNotFound = repository.ErrNotFound
	// ErrValidation 表示输入校验失败
	ErrValidation = errors.New("validation failed")
	// ErrDuplicate 敏感词已存在
	ErrDuplicate = errors.New("sensitive word already exists")
)

// SensitiveWordService 敏感词服务
type SensitiveWordService struct {
	repo  repository.SensitiveWordRepository
	cache *sensitiveWordCache
}

// sensitiveWordCache 敏感词缓存
type sensitiveWordCache struct {
	words      []model.SensitiveWord
	wordsMap   map[string]model.SensitiveWord
	lastUpdate time.Time
	mu         sync.RWMutex
	ttl        time.Duration
}

// NewSensitiveWordService 创建敏感词服务
func NewSensitiveWordService(repo repository.SensitiveWordRepository) *SensitiveWordService {
	return &SensitiveWordService{
		repo: repo,
		cache: &sensitiveWordCache{
			wordsMap: make(map[string]model.SensitiveWord),
			ttl:      5 * time.Minute, // 缓存5分钟
		},
	}
}

// AddSensitiveWordRequest 添加敏感词请求
type AddSensitiveWordRequest struct {
	Word     string                      `json:"word" binding:"required,max=100"`
	Category model.SensitiveWordCategory `json:"category" binding:"required"`
	Severity model.SensitiveWordSeverity `json:"severity" binding:"required"`
}

// UpdateSensitiveWordRequest 更新敏感词请求
type UpdateSensitiveWordRequest struct {
	Word     string                      `json:"word" binding:"required,max=100"`
	Category model.SensitiveWordCategory `json:"category" binding:"required"`
	Severity model.SensitiveWordSeverity `json:"severity" binding:"required"`
}

// SensitiveWordDTO 敏感词DTO
type SensitiveWordDTO struct {
	ID        uint64                      `json:"id"`
	Word      string                      `json:"word"`
	Category  model.SensitiveWordCategory `json:"category"`
	Severity  model.SensitiveWordSeverity `json:"severity"`
	CreatedAt string                      `json:"createdAt"`
	UpdatedAt string                      `json:"updatedAt"`
}

// ListSensitiveWordsRequest 列出敏感词请求
type ListSensitiveWordsRequest struct {
	Page     int                          `form:"page"`
	PageSize int                          `form:"pageSize"`
	Keyword  string                       `form:"keyword"`
	Category *model.SensitiveWordCategory `form:"category"`
	Severity *model.SensitiveWordSeverity `form:"severity"`
}

// ListSensitiveWordsResponse 列出敏感词响应
type ListSensitiveWordsResponse struct {
	Words []SensitiveWordDTO `json:"words"`
	Total int64              `json:"total"`
}

// DetectSensitiveWordsRequest 检测敏感词请求
type DetectSensitiveWordsRequest struct {
	Content string `json:"content" binding:"required"`
}

// DetectedWord 检测到的敏感词
type DetectedWord struct {
	Word      string                      `json:"word"`
	Category  model.SensitiveWordCategory `json:"category"`
	Severity  model.SensitiveWordSeverity `json:"severity"`
	Positions []int                       `json:"positions"` // 在文本中的所有位置
}

// DetectSensitiveWordsResponse 检测敏感词响应
type DetectSensitiveWordsResponse struct {
	HasSensitiveWords  bool           `json:"hasSensitiveWords"`
	DetectedWords      []DetectedWord `json:"detectedWords"`
	HighlightedContent string         `json:"highlightedContent"` // 高亮显示敏感词的文本
}

// AddSensitiveWord 添加敏感词
func (s *SensitiveWordService) AddSensitiveWord(ctx context.Context, req AddSensitiveWordRequest) (*SensitiveWordDTO, error) {
	// 验证输入
	if strings.TrimSpace(req.Word) == "" {
		return nil, ErrValidation
	}

	if !req.Category.Valid() {
		return nil, ErrValidation
	}

	if !req.Severity.Valid() {
		return nil, ErrValidation
	}

	// 创建敏感词
	word := &model.SensitiveWord{
		Word:     strings.TrimSpace(req.Word),
		Category: req.Category,
		Severity: req.Severity,
	}

	if err := s.repo.Create(ctx, word); err != nil {
		if strings.Contains(err.Error(), "已存在") {
			return nil, ErrDuplicate
		}
		return nil, err
	}

	// 清除缓存
	s.invalidateCache()

	return &SensitiveWordDTO{
		ID:        word.ID,
		Word:      word.Word,
		Category:  word.Category,
		Severity:  word.Severity,
		CreatedAt: word.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: word.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// UpdateSensitiveWord 更新敏感词
func (s *SensitiveWordService) UpdateSensitiveWord(ctx context.Context, id uint64, req UpdateSensitiveWordRequest) error {
	// 验证输入
	if strings.TrimSpace(req.Word) == "" {
		return ErrValidation
	}

	if !req.Category.Valid() {
		return ErrValidation
	}

	if !req.Severity.Valid() {
		return ErrValidation
	}

	// 获取现有敏感词
	word, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	// 更新字段
	word.Word = strings.TrimSpace(req.Word)
	word.Category = req.Category
	word.Severity = req.Severity

	if err := s.repo.Update(ctx, word); err != nil {
		if strings.Contains(err.Error(), "已存在") {
			return ErrDuplicate
		}
		return err
	}

	// 清除缓存
	s.invalidateCache()

	return nil
}

// DeleteSensitiveWord 删除敏感词
func (s *SensitiveWordService) DeleteSensitiveWord(ctx context.Context, id uint64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// 清除缓存
	s.invalidateCache()

	return nil
}

// ListSensitiveWords 列出敏感词
func (s *SensitiveWordService) ListSensitiveWords(ctx context.Context, req ListSensitiveWordsRequest) (*ListSensitiveWordsResponse, error) {
	// 默认分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 查询敏感词
	words, total, err := s.repo.List(ctx, repository.SensitiveWordListOptions{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		Category: req.Category,
		Severity: req.Severity,
	})
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	dtos := make([]SensitiveWordDTO, 0, len(words))
	for _, w := range words {
		dtos = append(dtos, SensitiveWordDTO{
			ID:        w.ID,
			Word:      w.Word,
			Category:  w.Category,
			Severity:  w.Severity,
			CreatedAt: w.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: w.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &ListSensitiveWordsResponse{
		Words: dtos,
		Total: total,
	}, nil
}

// DetectSensitiveWords 检测敏感词（检测并高亮）
func (s *SensitiveWordService) DetectSensitiveWords(ctx context.Context, req DetectSensitiveWordsRequest) (*DetectSensitiveWordsResponse, error) {
	if strings.TrimSpace(req.Content) == "" {
		return &DetectSensitiveWordsResponse{
			HasSensitiveWords:  false,
			DetectedWords:      []DetectedWord{},
			HighlightedContent: req.Content,
		}, nil
	}

	// 获取所有敏感词（使用缓存）
	words, err := s.getCachedSensitiveWords(ctx)
	if err != nil {
		return nil, err
	}

	// 检测敏感词，按词分组收集位置
	wordPositions := make(map[string]*DetectedWord)
	content := strings.ToLower(req.Content)

	for _, word := range words {
		searchWord := strings.ToLower(word.Word)
		pos := 0
		var positions []int

		for {
			index := strings.Index(content[pos:], searchWord)
			if index == -1 {
				break
			}
			actualPos := pos + index
			positions = append(positions, actualPos)
			pos = actualPos + len(searchWord)
		}

		if len(positions) > 0 {
			wordPositions[word.Word] = &DetectedWord{
				Word:      word.Word,
				Category:  word.Category,
				Severity:  word.Severity,
				Positions: positions,
			}
		}
	}

	// 转换为切片
	detectedWords := make([]DetectedWord, 0, len(wordPositions))
	for _, dw := range wordPositions {
		detectedWords = append(detectedWords, *dw)
	}

	// 高亮显示
	highlightedContent := req.Content
	if len(detectedWords) > 0 {
		// 收集所有位置并从后往前替换
		type posInfo struct {
			pos      int
			word     string
			severity model.SensitiveWordSeverity
		}
		var allPositions []posInfo
		for _, dw := range detectedWords {
			for _, p := range dw.Positions {
				allPositions = append(allPositions, posInfo{pos: p, word: dw.Word, severity: dw.Severity})
			}
		}
		// 按位置降序排序
		for i := 0; i < len(allPositions)-1; i++ {
			for j := i + 1; j < len(allPositions); j++ {
				if allPositions[i].pos < allPositions[j].pos {
					allPositions[i], allPositions[j] = allPositions[j], allPositions[i]
				}
			}
		}
		// 从后往前替换
		for _, pi := range allPositions {
			start := pi.pos
			end := start + len(pi.word)
			if end > len(highlightedContent) {
				continue
			}
			var marker string
			switch pi.severity {
			case model.SensitiveWordSeverityHigh:
				marker = "***"
			case model.SensitiveWordSeverityMedium:
				marker = "**"
			default:
				marker = "*"
			}
			highlightedContent = highlightedContent[:start] + marker + req.Content[start:end] + marker + highlightedContent[end:]
		}
	}

	return &DetectSensitiveWordsResponse{
		HasSensitiveWords:  len(detectedWords) > 0,
		DetectedWords:      detectedWords,
		HighlightedContent: highlightedContent,
	}, nil
}

// getCachedSensitiveWords 获取缓存的敏感词列表
func (s *SensitiveWordService) getCachedSensitiveWords(ctx context.Context) ([]model.SensitiveWord, error) {
	s.cache.mu.RLock()
	if time.Since(s.cache.lastUpdate) < s.cache.ttl && len(s.cache.words) > 0 {
		words := s.cache.words
		s.cache.mu.RUnlock()
		return words, nil
	}
	s.cache.mu.RUnlock()

	// 缓存过期或为空，重新加载
	s.cache.mu.Lock()
	defer s.cache.mu.Unlock()

	// 双重检查
	if time.Since(s.cache.lastUpdate) < s.cache.ttl && len(s.cache.words) > 0 {
		return s.cache.words, nil
	}

	// 从数据库加载
	words, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// 更新缓存
	s.cache.words = words
	s.cache.wordsMap = make(map[string]model.SensitiveWord)
	for _, w := range words {
		s.cache.wordsMap[strings.ToLower(w.Word)] = w
	}
	s.cache.lastUpdate = time.Now()

	return words, nil
}

// invalidateCache 清除缓存
func (s *SensitiveWordService) invalidateCache() {
	s.cache.mu.Lock()
	defer s.cache.mu.Unlock()

	s.cache.words = nil
	s.cache.wordsMap = make(map[string]model.SensitiveWord)
	s.cache.lastUpdate = time.Time{}
}
