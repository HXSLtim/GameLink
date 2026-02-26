package sms

import (
	"context"
	"testing"
)

func TestAliyunProvider_ProductionMustFailWhenIntegrationMissing(t *testing.T) {
	t.Setenv("APP_ENV", "production")

	provider := &AliyunProvider{SignName: "GameLink"}

	if err := provider.SendCode(context.Background(), "13800138000", "123456"); err == nil {
		t.Fatalf("expected production SendCode to fail when aliyun integration is missing")
	}

	if err := provider.SendNotification(context.Background(), "13800138000", "tpl-1", map[string]string{"code": "123456"}); err == nil {
		t.Fatalf("expected production SendNotification to fail when aliyun integration is missing")
	}
}

func TestAliyunProvider_NonProductionKeepsCurrentDeveloperExperience(t *testing.T) {
	t.Setenv("APP_ENV", "development")

	provider := &AliyunProvider{SignName: "GameLink"}

	if err := provider.SendCode(context.Background(), "13800138000", "123456"); err != nil {
		t.Fatalf("expected non-production SendCode to keep current behavior, got error: %v", err)
	}

	if err := provider.SendNotification(context.Background(), "13800138000", "tpl-1", map[string]string{"code": "123456"}); err != nil {
		t.Fatalf("expected non-production SendNotification to keep current behavior, got error: %v", err)
	}
}

func TestTencentProvider_ProductionMustFailWhenIntegrationMissing(t *testing.T) {
	t.Setenv("APP_ENV", "production")

	provider := &TencentProvider{SignName: "GameLink"}

	if err := provider.SendCode(context.Background(), "13800138000", "123456"); err == nil {
		t.Fatalf("expected production SendCode to fail when tencent integration is missing")
	}

	if err := provider.SendNotification(context.Background(), "13800138000", "tpl-1", map[string]string{"1": "123456"}); err == nil {
		t.Fatalf("expected production SendNotification to fail when tencent integration is missing")
	}
}

func TestTencentProvider_NonProductionKeepsCurrentDeveloperExperience(t *testing.T) {
	t.Setenv("APP_ENV", "development")

	provider := &TencentProvider{SignName: "GameLink"}

	if err := provider.SendCode(context.Background(), "13800138000", "123456"); err != nil {
		t.Fatalf("expected non-production SendCode to keep current behavior, got error: %v", err)
	}

	if err := provider.SendNotification(context.Background(), "13800138000", "tpl-1", map[string]string{"1": "123456"}); err != nil {
		t.Fatalf("expected non-production SendNotification to keep current behavior, got error: %v", err)
	}
}
