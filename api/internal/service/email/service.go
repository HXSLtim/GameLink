package email

import (
	"context"
	"fmt"
	"os"
)

type Service struct {
	enabled bool
}

func NewService() *Service {
	enabled := os.Getenv("EMAIL_ENABLED")
	return &Service{
		enabled: enabled == "true" || enabled == "1",
	}
}

var (
	ErrEmailDisabled = fmt.Errorf("email service is disabled")
)

func (s *Service) SendPasswordReset(ctx context.Context, to, resetURL string) error {
	if !s.enabled {
		fmt.Printf("[Email] Email not enabled, password reset URL: %s\n", resetURL)
		return ErrEmailDisabled
	}

	subject := "重置您的 GameLink 密码"
	body := fmt.Sprintf(`
您好，

您收到此邮件是因为您请求重置 GameLink 账号的密码。

请点击以下链接重置密码（链接有效期为1小时）：
%s

如果您没有请求重置密码，请忽略此邮件。

GameLink 团队
`, resetURL)

	fmt.Printf("[Email] Sending password reset to %s\n", to)
	fmt.Printf("[Email] Subject: %s\n", subject)
	fmt.Printf("[Email] Body length: %d characters\n", len(body))

	return nil
}
