package xjwt

import (
	"errors"
	"testing"
	"time"
)

// TestGenerateParseAndRefresh 验证 JWT 生成、解析和刷新流程。
func TestGenerateParseAndRefresh(t *testing.T) {
	conf := Config{
		Secret:  "secret",
		Issuer:  "issuer",
		Subject: "subject",
		Expire:  time.Minute,
	}

	token, err := Generate(conf, map[string]any{"user_id": float64(1001)})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	claims, err := Parse(conf, token)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if claims.Issuer != conf.Issuer || claims.Subject != conf.Subject {
		t.Fatalf("claims issuer/subject mismatch: %#v", claims.RegisteredClaims)
	}
	if claims.Data["user_id"] != float64(1001) {
		t.Fatalf("claims data mismatch: %#v", claims.Data)
	}

	refreshed, err := Refresh(conf, token)
	if err != nil {
		t.Fatalf("Refresh() error = %v", err)
	}
	if refreshed == "" {
		t.Fatal("Refresh() returned empty token")
	}
}

// TestJWTMissingSecret 验证缺少签名密钥时返回统一错误。
func TestJWTMissingSecret(t *testing.T) {
	if _, err := Generate(Config{}, nil); !errors.Is(err, ErrMissingSecret) {
		t.Fatalf("Generate() error = %v, want %v", err, ErrMissingSecret)
	}
	if _, err := Parse(Config{}, "token"); !errors.Is(err, ErrMissingSecret) {
		t.Fatalf("Parse() error = %v, want %v", err, ErrMissingSecret)
	}
}
