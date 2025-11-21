package safety

import "testing"

func TestContainsSensitiveWord_Positive(t *testing.T) {
	cases := []string{
		"这是违法内容",
		"这条消息涉嫌违规",
		"包含spam关键词",
		"FAKE news here",
	}

	for _, tc := range cases {
		if !ContainsSensitiveWord(tc) {
			t.Fatalf("expected text %q to contain sensitive word", tc)
		}
	}
}

func TestContainsSensitiveWord_Negative(t *testing.T) {
	cases := []string{
		"正常交流内容",
		"欢迎使用 GameLink 服务",
	}

	for _, tc := range cases {
		if ContainsSensitiveWord(tc) {
			t.Fatalf("expected text %q to be clean", tc)
		}
	}
}

func TestValidateText_Empty(t *testing.T) {
	if err := ValidateText("   \t\n  ", 10); err == nil {
		t.Fatalf("expected error for empty content")
	}
}

func TestValidateText_TooLong(t *testing.T) {
	// 这里使用简单的 ASCII 文本，超过 maxLen 即视为超长
	content := "abcdef"
	if err := ValidateText(content, 3); err == nil {
		t.Fatalf("expected error for too long content")
	}
}

func TestValidateText_Sensitive(t *testing.T) {
	if err := ValidateText("这是一条违法内容", 50); err == nil {
		t.Fatalf("expected error when sensitive word is present")
	}
}

func TestValidateText_Valid(t *testing.T) {
	if err := ValidateText("正常的聊天内容", 50); err != nil {
		t.Fatalf("expected no error for valid content, got: %v", err)
	}
}
