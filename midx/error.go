package midx

import (
	"errors"
	"net/http"
)

var (
	ErrMissingToken = errors.New("midx: missing token") // ErrMissingToken 表示 token 缺失。
)

// defaultUnauthorized 输出默认鉴权失败响应。
func defaultUnauthorized(w http.ResponseWriter, _ *http.Request, _ error) {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}
