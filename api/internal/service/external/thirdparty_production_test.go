package external

import (
	"context"
	"testing"
)

func TestNewIdentityVerifier_ProductionMustNotSilentMockSuccess(t *testing.T) {
	t.Setenv("APP_ENV", "production")

	verifier := NewIdentityVerifier(nil)
	result, err := verifier.VerifyIdentity(context.Background(), "AB", "110101199003076019")
	if err == nil {
		t.Fatalf("expected production identity verification to fail when provider is not integrated")
	}
	if result != nil && result.Verified {
		t.Fatalf("production identity verification must not silently mock success")
	}
}

func TestNewIdentityVerifier_ProductionRealProviderStillFailsIfNotImplemented(t *testing.T) {
	t.Setenv("APP_ENV", "production")

	verifier := NewIdentityVerifier(&ThirdPartyConfig{
		Mode: ModeReal,
		Identity: IdentityConfig{
			Provider: "aliyun",
		},
	})
	result, err := verifier.VerifyIdentity(context.Background(), "AB", "110101199003076019")
	if err == nil {
		t.Fatalf("expected production identity verification to fail until real provider integration is implemented")
	}
	if result != nil && result.Verified {
		t.Fatalf("production identity verification must not silently return success")
	}
}

func TestNewIdentityVerifier_NonProductionKeepsCurrentDeveloperExperience(t *testing.T) {
	t.Setenv("APP_ENV", "development")

	verifier := NewIdentityVerifier(nil)
	result, err := verifier.VerifyIdentity(context.Background(), "AB", "110101199003076019")
	if err != nil {
		t.Fatalf("expected non-production identity verifier to keep current mock behavior, got error: %v", err)
	}
	if result == nil || !result.Verified {
		t.Fatalf("expected non-production identity verifier to return mock success for valid input")
	}
}

func TestNewPaymentCallbackVerifier_ProductionMustNotSilentMockSuccess(t *testing.T) {
	t.Setenv("APP_ENV", "production")

	verifier := NewPaymentCallbackVerifier(nil)
	wechatOK, wechatErr := verifier.VerifyWeChatCallback([]byte(`{}`), "sig")
	if wechatErr == nil || wechatOK {
		t.Fatalf("expected production wechat callback verification to fail when integration is missing")
	}

	alipayOK, alipayErr := verifier.VerifyAlipayCallback(map[string]string{"sign": "sig"})
	if alipayErr == nil || alipayOK {
		t.Fatalf("expected production alipay callback verification to fail when integration is missing")
	}
}

func TestNewPaymentCallbackVerifier_NonProductionKeepsCurrentDeveloperExperience(t *testing.T) {
	t.Setenv("APP_ENV", "development")

	verifier := NewPaymentCallbackVerifier(nil)
	wechatOK, wechatErr := verifier.VerifyWeChatCallback([]byte(`{}`), "sig")
	if wechatErr != nil || !wechatOK {
		t.Fatalf("expected non-production wechat callback verification to keep mock success behavior")
	}

	alipayOK, alipayErr := verifier.VerifyAlipayCallback(map[string]string{"sign": "sig"})
	if alipayErr != nil || !alipayOK {
		t.Fatalf("expected non-production alipay callback verification to keep mock success behavior")
	}
}
