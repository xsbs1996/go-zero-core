package xpacket

import "errors"

var (
	ErrPacketTooShort       = errors.New("xpacket: packet too short")       // ErrPacketTooShort 表示二进制包长度不足 6 字节头部。
	ErrPacketLengthMismatch = errors.New("xpacket: packet length mismatch") // ErrPacketLengthMismatch 表示头部 bodyLen 与实际 body 长度不一致。
	ErrNilMessage           = errors.New("xpacket: nil decode target")      // ErrNilMessage 表示解码目标为空。
	ErrBodyTooLarge         = errors.New("xpacket: body too large")         // ErrBodyTooLarge 表示 body 长度超过 uint32 可表达的最大值。
)
