package xticket

import "errors"

var (
	ErrMissingSecret = errors.New("xticket: missing secret") // ErrMissingSecret 表示签名密钥为空。
	ErrInvalidTicket = errors.New("xticket: invalid ticket") // ErrInvalidTicket 表示票据格式或签名非法。
	ErrTicketExpired = errors.New("xticket: ticket expired") // ErrTicketExpired 表示票据已经过期。
)
