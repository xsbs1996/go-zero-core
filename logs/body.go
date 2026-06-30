package logs

import (
	"encoding/json"
	"fmt"
)

// Content 日志正文结构。
type Content struct {
	Msg     string `json:"msg"`     // Msg 日志消息。
	Error   any    `json:"error"`   // Error 错误信息。
	Content any    `json:"content"` // Content 日志正文内容。
}

// Body 将日志正文包装为统一 JSON 结构。
func Body(msg string, content any, err any) string {
	data, marshalErr := json.Marshal(Content{
		Msg:     msg,
		Error:   normalizeBody(err),
		Content: normalizeBody(content),
	})
	if marshalErr != nil {
		return `{"msg":"","error":null,"content":null}`
	}

	return string(data)
}

// BodyContent 将日志正文结构包装为 JSON 字符串。
func BodyContent(content Content) string {
	return Body(content.Msg, content.Content, content.Error)
}

func normalizeBody(content any) any {
	switch value := content.(type) {
	case nil:
		return nil
	case error:
		return value.Error()
	case fmt.Stringer:
		return value.String()
	default:
		return value
	}
}
