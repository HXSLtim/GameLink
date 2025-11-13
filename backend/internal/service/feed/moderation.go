package feed

import (
	"context"

	"gamelink/internal/pkg/safety"
)

// ModerationDecision enumerates automatic moderation outcomes.
type ModerationDecision string

const (
	// ModerationDecisionApprove indicates content passed automated checks.
	ModerationDecisionApprove ModerationDecision = "approve"
	// ModerationDecisionReject indicates content must be rejected immediately.
	ModerationDecisionReject ModerationDecision = "reject"
	// ModerationDecisionManual indicates manual review is required.
	ModerationDecisionManual ModerationDecision = "manual_review"
)

// ModerationInput contains content data for evaluation.
type ModerationInput struct {
	Content   string
	ImageURLs []string
}

// ModerationResult represents evaluation output.
type ModerationResult struct {
	Decision ModerationDecision
	Reason   string
}

// ModerationEngine abstracts content moderation pipeline.
type ModerationEngine interface {
	Evaluate(ctx context.Context, input ModerationInput) (ModerationResult, error)
}

// NewDefaultModerationEngine returns a lightweight moderation engine using safety utilities.
func NewDefaultModerationEngine() ModerationEngine {
	return &simpleModerationEngine{}
}

type simpleModerationEngine struct{}

func (s *simpleModerationEngine) Evaluate(ctx context.Context, input ModerationInput) (ModerationResult, error) {
	if safety.ContainsSensitiveWord(input.Content) {
		return ModerationResult{Decision: ModerationDecisionReject, Reason: "文本触发敏感词"}, nil
	}
	for _, url := range input.ImageURLs {
		if safety.ContainsSensitiveWord(url) {
			return ModerationResult{Decision: ModerationDecisionReject, Reason: "图片命中敏感词"}, nil
		}
	}
	// 默认通过自动审核，留给人工复审通道处理举报等。
	return ModerationResult{Decision: ModerationDecisionApprove}, nil
}
