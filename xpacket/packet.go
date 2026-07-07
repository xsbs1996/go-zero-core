// Package xpacket 提供与传输层无关的二进制封包能力。
//
// 默认封包格式固定为 6 字节头部加 body：
//
//	[0:2] action  uint16，大端序，表示业务操作码。
//	[2:6] bodyLen uint32，大端序，表示 body 字节长度。
//	[6:]  body    任意二进制载荷，例如 protobuf、JSON、MsgPack。
//
// 当前包把通用 envelope 与具体 body 编解码拆开：
//   - EncodeBody、DecodeBody 只处理 action + bodyLen + body。
//   - EncodeProto、DecodeProto 在通用 envelope 之上处理 protobuf。
//
// 后续如果增加新的 body 编码方式，可以新增 EncodeJson、DecodeJson 等函数；
// 如果增加新的封包格式，可以新增独立文件和函数，避免影响当前默认格式。
package xpacket

import (
	"encoding/binary"
	"math"
)

// HeaderLen 表示默认二进制封包头部长度。
const HeaderLen = 6

// EncodeBody 将 action 和任意二进制 body 编码为默认封包格式。
//
// 参数：
//   - action: 业务操作码，会写入包头 [0:2]。
//   - body: 二进制载荷，会写入包体 [6:]；可以为 nil 或空切片。
//
// 返回值：
//   - []byte: 编码后的完整二进制包。
//   - error: body 超过 uint32 长度上限或封包长度校验失败时返回错误。
//
// body 可以为 nil 或空切片，此时 bodyLen 写入 0。
// body 长度超过 uint32 最大值时返回 ErrBodyTooLarge。
func EncodeBody(action uint16, body []byte) ([]byte, error) {
	if len(body) > math.MaxUint32 {
		return nil, ErrBodyTooLarge
	}

	packet := make([]byte, HeaderLen+len(body))
	binary.BigEndian.PutUint16(packet[:2], action)
	binary.BigEndian.PutUint32(packet[2:HeaderLen], uint32(len(body)))
	copy(packet[HeaderLen:], body)

	if err := ValidateLength(packet); err != nil {
		return nil, err
	}
	return packet, nil
}

// DecodeBody 校验并解出默认封包格式中的 action 和 body。
//
// 参数：
//   - data: 完整二进制包，格式必须为 HeaderLen + body。
//
// 返回值：
//   - action: 包头中的业务操作码。
//   - body: 包体切片。
//   - err: 包过短或长度不匹配时返回错误。
//
// 返回的 body 是 data 的切片视图，不会额外拷贝；调用方如果需要长期持有或修改，
// 应自行 copy。
func DecodeBody(data []byte) (action uint16, body []byte, err error) {
	if err := ValidateLength(data); err != nil {
		return 0, nil, err
	}

	action = binary.BigEndian.Uint16(data[:2])
	bodyLen := int(binary.BigEndian.Uint32(data[2:HeaderLen]))
	return action, data[HeaderLen : HeaderLen+bodyLen], nil
}

// ValidateLength 校验默认封包格式的长度是否合法。
//
// 参数：
//   - data: 完整二进制包。
//
// 返回值：
//   - error: 长度合法返回 nil；包头不足或 bodyLen 与实际长度不一致时返回错误。
//
// 合法包必须至少包含 6 字节头部，并且实际长度必须等于 HeaderLen + bodyLen。
func ValidateLength(data []byte) error {
	if len(data) < HeaderLen {
		return ErrPacketTooShort
	}

	bodyLen := binary.BigEndian.Uint32(data[2:HeaderLen])
	if len(data) != HeaderLen+int(bodyLen) {
		return ErrPacketLengthMismatch
	}
	return nil
}

// GetAction 只读取默认封包格式头部中的 action。
//
// 参数：
//   - data: 完整二进制包。
//
// 返回值：
//   - action: 包头中的业务操作码。
//   - err: 包长度非法时返回错误。
//
// 该函数仍会先校验完整包长度，适合在路由分发前做轻量预解析。
func GetAction(data []byte) (action uint16, err error) {
	if err := ValidateLength(data); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(data[:2]), nil
}
