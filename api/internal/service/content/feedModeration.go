package content

import (
	"context"

	"gamelink/pkg/safety"
)

// FeedModerationDecision enumerates automatic moderation outcomes.
type FeedModerationDecision string

const (
	// FeedModerationDecisionApprove indicates content passed automated checks.
	FeedModerationDecisionApprove FeedModerationDecision = "approve"
	// FeedModerationDecisionReject indicates content must be rejected immediately.
	FeedModerationDecisionReject FeedModerationDecision = "reject"
	// FeedModerationDecisionManual indicates manual review is required.
	FeedModerationDecisionManual FeedModerationDecision = "manual_review"
)

// FeedModerationInput contains content data for evaluation.
type FeedModerationInput struct {
	Content   string
	ImageURLs []string
}

// FeedModerationResult represents evaluation output.
type FeedModerationResult struct {
	Decision FeedModerationDecision
	Reason   string
}

// FeedModerationEngine abstracts content moderation pipeline.
type FeedModerationEngine interface {
	Evaluate(ctx context.Context, input FeedModerationInput) (FeedModerationResult, error)
}

// NewDefaultFeedModerationEngine returns a lightweight moderation engine using safety utilities.
func NewDefaultFeedModerationEngine() FeedModerationEngine {
	return &simpleFeedModerationEngine{}
}

type simpleFeedModerationEngine struct{}

func (s *simpleFeedModerationEngine) Evaluate(ctx context.Context, input FeedModerationInput) (FeedModerationResult, error) {
	if safety.ContainsSensitiveWord(input.Content) {
		return FeedModerationResult{Decision: FeedModerationDecisionReject, Reason: "文本触发敏感词"}, nil
	}
	for _, url := range input.ImageURLs {
		if safety.ContainsSensitiveWord(url) {
			return FeedModerationResult{Decision: FeedModerationDecisionReject, Reason: "图片命中敏感词"}, nil
		}
	}
	// 默认通过自动审核，留给人工复审通道处理举报等。
	return FeedModerationResult{Decision: FeedModerationDecisionApprove}, nil
}
