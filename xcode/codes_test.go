package xcode

import (
	"errors"
	"testing"
)

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

func TestRegisterCodesRejectsReservedCode(t *testing.T) {
	t.Parallel()

	err := RegisterCodes(map[int]string{CodeInternal: "override"})
	if !errors.Is(err, ErrReservedCode) {
		t.Fatalf("expected ErrReservedCode, got %v", err)
	}
}
