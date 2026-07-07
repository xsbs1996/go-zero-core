package xticket

import (
	"errors"
	"strings"
	"testing"
	"time"
)

type payload struct {
	UserID   uint64 `json:"user_id"`
	NickName string `json:"nickname"`
}

// TestGenerateVerify 验证票据生成、验签和业务载荷解析。
func TestGenerateVerify(t *testing.T) {
	t.Parallel()

	conf := Config{
		Secret: []byte("secret"),
		Issuer: "test",
		TTL:    time.Minute,
	}

	value, err := Generate(conf, payload{UserID: 1001, NickName: "alice"})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	claims, err := Verify[payload](conf, value)
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if claims.Payload.UserID != 1001 || claims.Payload.NickName != "alice" {
		t.Fatalf("payload = %+v, want user 1001 alice", claims.Payload)
	}
	if claims.Issuer != "test" || claims.Nonce == "" {
		t.Fatalf("claims = %+v, want issuer and nonce", claims)
	}
}

// TestVerifyRejectsTamperedTicket 验证票据被篡改后验签失败。
func TestVerifyRejectsTamperedTicket(t *testing.T) {
	t.Parallel()

	conf := Config{Secret: []byte("secret"), TTL: time.Minute}
	value, err := Generate(conf, payload{UserID: 1001})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	tampered := strings.Replace(value, ".", ".x", 1)
	_, err = Verify[payload](conf, tampered)
	if !errors.Is(err, ErrInvalidTicket) {
		t.Fatalf("expected ErrInvalidTicket, got %v", err)
	}
}

// TestVerifyRejectsExpiredTicket 验证过期票据会被拒绝。
func TestVerifyRejectsExpiredTicket(t *testing.T) {
	t.Parallel()

	conf := Config{Secret: []byte("secret"), TTL: time.Nanosecond}
	value, err := Generate(conf, payload{UserID: 1001})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	time.Sleep(time.Millisecond)
	_, err = Verify[payload](conf, value)
	if !errors.Is(err, ErrTicketExpired) {
		t.Fatalf("expected ErrTicketExpired, got %v", err)
	}
}
