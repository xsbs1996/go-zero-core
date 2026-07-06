package xpacket

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

// EncodeProto 将 action 和 protobuf 消息编码为默认二进制封包格式。
//
// msg 为 nil 时会编码一个只有头部、bodyLen 为 0 的空 body 包。
// msg 不为空时会先执行 proto.Marshal，序列化失败会返回错误。
func EncodeProto(action uint16, msg proto.Message) ([]byte, error) {
	var body []byte
	if msg != nil {
		var err error
		body, err = proto.Marshal(msg)
		if err != nil {
			return nil, fmt.Errorf("xpacket: marshal protobuf: %w", err)
		}
	}
	return EncodeBody(action, body)
}

// DecodeProto 将默认二进制封包格式解码到 protobuf 消息，并返回包头中的 action。
//
// 解码流程固定为：
//  1. 校验包长度，避免读取越界或半包数据被误处理。
//  2. 按 bodyLen 切出 protobuf body，并执行 proto.Unmarshal。
//  3. 使用 protovalidate 做 protobuf 规则校验。
//  4. 如果 msg 实现 ProtoBusinessValidator，则继续执行 protobuf 载荷业务校验。
func DecodeProto(data []byte, msg proto.Message) (action uint16, err error) {
	if msg == nil {
		return 0, ErrNilMessage
	}

	action, body, err := DecodeBody(data)
	if err != nil {
		return 0, err
	}

	if err := proto.Unmarshal(body, msg); err != nil {
		return action, fmt.Errorf("xpacket: unmarshal protobuf: %w", err)
	}

	if err := validateProto(msg); err != nil {
		return action, err
	}

	if err := validateBusinessProto(msg); err != nil {
		return action, err
	}

	return action, nil
}
