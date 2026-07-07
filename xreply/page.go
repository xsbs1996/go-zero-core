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
//
// 参数：
//   - w: HTTP 响应写入器。
//   - list: 当前页数据列表。
//   - total: 总数据量。
//   - pageNo: 当前页码。
//   - pageSize: 每页数量。
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
