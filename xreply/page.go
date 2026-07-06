package xreply

import (
	"net/http"

	"github.com/xsbs1996/go-zero-core/xcode"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type page[T any] struct {
	List     []T   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page,omitempty"`
	PageSize int   `json:"pageSize,omitempty"`
}

// SuccessPage 输出分页成功 JSON 响应。
func SuccessPage[T any](w http.ResponseWriter, list []T, total int64, pageNo int, pageSize int) {
	httpx.OkJson(w, successPageResult(list, total, pageNo, pageSize))
}

func successPageResult[T any](list []T, total int64, pageNo int, pageSize int) Result[page[T]] {
	return NewResult(xcode.CodeSuccess, xcode.MsgSuccess, page[T]{
		List:     list,
		Total:    total,
		Page:     pageNo,
		PageSize: pageSize,
	})
}
