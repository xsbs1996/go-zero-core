package xpacket

import (
	"fmt"
	"sync"

	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

var (
	validatorOnce sync.Once
	validator     protovalidate.Validator
	validatorErr  error
)

// SetProtoValidator 替换包级 protobuf 结构校验器。
//
// 通常在服务启动阶段调用；传入 nil 表示关闭 protovalidate 结构校验，只保留
// protobuf 反序列化和 ProtoBusinessValidator 业务校验。
func SetProtoValidator(v protovalidate.Validator) {
	validatorOnce.Do(func() {})
	validator = v
	validatorErr = nil
}

// ProtoBusinessValidator 表示 protobuf 反序列化和 protovalidate 之后的业务校验器。
//
// DecodeProto 会在 proto.Unmarshal 和 protovalidate 成功后调用该接口，
// 适合校验余额、状态、跨字段业务关系等 protobuf 载荷对应的业务规则。
type ProtoBusinessValidator interface {
	ValidateBusinessProto() error
}

func validateBusinessProto(msg proto.Message) error {
	if vdt, ok := msg.(ProtoBusinessValidator); ok {
		return vdt.ValidateBusinessProto()
	}
	return nil
}

func validateProto(msg proto.Message) error {
	validatorOnce.Do(func() {
		validator, validatorErr = protovalidate.New()
	})
	if validatorErr != nil {
		return fmt.Errorf("xpacket: init protobuf validator: %w", validatorErr)
	}
	if validator == nil {
		return nil
	}
	if err := validator.Validate(msg); err != nil {
		return fmt.Errorf("xpacket: validate protobuf: %w", err)
	}
	return nil
}
