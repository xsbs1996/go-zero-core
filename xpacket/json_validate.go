package xpacket

// JsonBusinessValidator 表示 JSON 反序列化之后的业务校验器。
//
// DecodeJson 会在 json.Unmarshal 成功后调用该接口，适合校验余额、状态、
// 跨字段业务关系等 JSON 载荷对应的业务规则。
type JsonBusinessValidator interface {
	ValidateBusinessJson() error
}

func validateBusinessJson(v any) error {
	if vdt, ok := v.(JsonBusinessValidator); ok {
		return vdt.ValidateBusinessJson()
	}
	return nil
}
