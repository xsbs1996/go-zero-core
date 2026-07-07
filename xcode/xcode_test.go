package xcode

import (
	"errors"
	"testing"
)

// TestMsgRendersVarsAndCleansMissingPlaceholders 验证错误文案变量渲染和缺失变量清理。
func TestMsgRendersVarsAndCleansMissingPlaceholders(t *testing.T) {
	t.Parallel()

	got := Msg(CodeInvalidParam, Vars{"field": "name"})
	if got != "invalid param: name" {
		t.Fatalf("Msg() = %q, want %q", got, "invalid param: name")
	}

	got = Msg(CodeInvalidParam)
	if got != "invalid param" {
		t.Fatalf("Msg() without vars = %q, want %q", got, "invalid param")
	}
}

// TestRegisterCodesRejectsReservedCode 验证 0-99 内置错误码区间不允许业务注册。
func TestRegisterCodesRejectsReservedCode(t *testing.T) {
	t.Parallel()

	err := RegisterCodes(map[int]string{CodeInternal: "override"})
	if !errors.Is(err, ErrReservedCode) {
		t.Fatalf("expected ErrReservedCode, got %v", err)
	}
}

// TestRegisterCodesAllowsBusinessCodeFrom100 验证业务错误码从 100 开始允许注册。
func TestRegisterCodesAllowsBusinessCodeFrom100(t *testing.T) {
	code := 100
	msg := "business error"

	if err := RegisterCodes(map[int]string{code: msg}); err != nil {
		t.Fatalf("RegisterCodes() error = %v", err)
	}
	if got := Msg(code); got != msg {
		t.Fatalf("Msg(%d) = %q, want %q", code, got, msg)
	}
}
