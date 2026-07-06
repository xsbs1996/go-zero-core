package xreply

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xsbs1996/go-zero-core/xcode"
)

// TestNewResult 验证统一响应结构创建。
func TestNewResult(t *testing.T) {
	got := NewResult(100, "ok", "data")
	if got.Code != 100 || got.Msg != "ok" || got.Data != "data" {
		t.Fatalf("NewResult() = %#v", got)
	}
}

// TestSuccessAndFailResponses 验证成功和失败响应的 JSON 输出。
func TestSuccessAndFailResponses(t *testing.T) {
	w := httptest.NewRecorder()
	Success(w, map[string]any{"id": 1})
	assertResult(t, w, http.StatusOK, xcode.CodeSuccess, xcode.MsgSuccess)

	w = httptest.NewRecorder()
	FailStatus(w, http.StatusBadRequest, xcode.CodeError)
	assertResult(t, w, http.StatusBadRequest, xcode.CodeError, xcode.MsgError)
}

// TestSuccessPageResponse 验证分页成功响应的 JSON 输出。
func TestSuccessPageResponse(t *testing.T) {
	w := httptest.NewRecorder()
	SuccessPage(w, []int{1, 2}, 2, 1, 10)
	assertResult(t, w, http.StatusOK, xcode.CodeSuccess, xcode.MsgSuccess)
}

// assertResult 解析响应并校验状态码、业务码和文案。
func assertResult(t *testing.T, w *httptest.ResponseRecorder, status int, code int, msg string) {
	t.Helper()
	if w.Code != status {
		t.Fatalf("status = %d, want %d", w.Code, status)
	}

	var result Result[any]
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if result.Code != code || result.Msg != msg {
		t.Fatalf("response = %#v, want code=%d msg=%q", result, code, msg)
	}
}
