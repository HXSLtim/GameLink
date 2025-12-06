package content

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeedModerationDecision_Constants(t *testing.T) {
	// Verify constant values
	assert.Equal(t, FeedModerationDecision("approve"), FeedModerationDecisionApprove)
	assert.Equal(t, FeedModerationDecision("reject"), FeedModerationDecisionReject)
	assert.Equal(t, FeedModerationDecision("manual_review"), FeedModerationDecisionManual)
}

func TestNewDefaultFeedModerationEngine(t *testing.T) {
	engine := NewDefaultFeedModerationEngine()
	require.NotNil(t, engine)

	// Verify it implements the interface
	_, ok := engine.(FeedModerationEngine)
	assert.True(t, ok, "should implement FeedModerationEngine interface")
}

func TestSimpleFeedModerationEngine_Evaluate(t *testing.T) {
	engine := NewDefaultFeedModerationEngine()
	ctx := context.Background()

	tests := []struct {
		name          string
		input         FeedModerationInput
		wantDecision  FeedModerationDecision
		wantReasonNot string // reason should not be empty if this is set
	}{
		{
			name: "clean content should be approved",
			input: FeedModerationInput{
				Content:   "这是一条正常的动态内容",
				ImageURLs: []string{},
			},
			wantDecision: FeedModerationDecisionApprove,
		},
		{
			name: "empty content should be approved",
			input: FeedModerationInput{
				Content:   "",
				ImageURLs: []string{},
			},
			wantDecision: FeedModerationDecisionApprove,
		},
		{
			name: "content with clean images should be approved",
			input: FeedModerationInput{
				Content:   "分享一下今天的游戏截图",
				ImageURLs: []string{"https://example.com/image1.jpg", "https://example.com/image2.png"},
			},
			wantDecision: FeedModerationDecisionApprove,
		},
		{
			name: "multiple clean images should be approved",
			input: FeedModerationInput{
				Content:   "",
				ImageURLs: []string{"https://cdn.example.com/a.jpg", "https://cdn.example.com/b.jpg", "https://cdn.example.com/c.jpg"},
			},
			wantDecision: FeedModerationDecisionApprove,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.Evaluate(ctx, tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.wantDecision, result.Decision)
		})
	}
}

func TestSimpleFeedModerationEngine_Evaluate_ContextCancellation(t *testing.T) {
	engine := NewDefaultFeedModerationEngine()

	// Test with cancelled context - should still work as current impl doesn't check context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	input := FeedModerationInput{
		Content: "正常内容",
	}

	result, err := engine.Evaluate(ctx, input)
	// Current implementation doesn't check context, so it should succeed
	require.NoError(t, err)
	assert.Equal(t, FeedModerationDecisionApprove, result.Decision)
}

func TestFeedModerationInput_Fields(t *testing.T) {
	input := FeedModerationInput{
		Content:   "测试内容",
		ImageURLs: []string{"url1", "url2"},
	}

	assert.Equal(t, "测试内容", input.Content)
	assert.Len(t, input.ImageURLs, 2)
	assert.Equal(t, "url1", input.ImageURLs[0])
	assert.Equal(t, "url2", input.ImageURLs[1])
}

func TestFeedModerationResult_Fields(t *testing.T) {
	result := FeedModerationResult{
		Decision: FeedModerationDecisionReject,
		Reason:   "违规内容",
	}

	assert.Equal(t, FeedModerationDecisionReject, result.Decision)
	assert.Equal(t, "违规内容", result.Reason)
}

// Benchmark tests
func BenchmarkSimpleFeedModerationEngine_Evaluate(b *testing.B) {
	engine := NewDefaultFeedModerationEngine()
	ctx := context.Background()
	input := FeedModerationInput{
		Content:   "这是一条测试动态，分享游戏心得和攻略",
		ImageURLs: []string{"https://example.com/img1.jpg", "https://example.com/img2.jpg"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Evaluate(ctx, input)
	}
}

func BenchmarkSimpleFeedModerationEngine_Evaluate_LongContent(b *testing.B) {
	engine := NewDefaultFeedModerationEngine()
	ctx := context.Background()

	// Generate long content
	longContent := ""
	for i := 0; i < 100; i++ {
		longContent += "这是一段很长的测试内容，用于测试性能。"
	}

	input := FeedModerationInput{
		Content:   longContent,
		ImageURLs: []string{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Evaluate(ctx, input)
	}
}
