package kafkax

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/segmentio/kafka-go"
)

var (
	ErrNoAvailableBroker = errors.New("kafkax: no available kafka broker") // ErrNoAvailableBroker 表示没有可用的 Kafka broker。
)

// Ping 检查 Kafka broker 连通性，至少一个 broker 可连接即成功。
func Ping(ctx context.Context, conf Config) error {
	if err := conf.ValidateProducer(); err != nil {
		return err
	}

	conf = conf.WithDefault()
	dialer := &kafka.Dialer{
		ClientID: conf.ClientID,
		Timeout:  conf.DialTimeout,
	}

	var errs []error
	for _, broker := range conf.Brokers {
		broker = strings.TrimSpace(broker)
		if broker == "" {
			continue
		}

		conn, err := dialer.DialContext(ctx, "tcp", broker)
		if err == nil {
			return conn.Close()
		}
		errs = append(errs, fmt.Errorf("dial kafka broker %q failed: %w", broker, err))
	}

	if len(errs) == 0 {
		return ErrMissingBrokers
	}
	return fmt.Errorf("%w: %w", ErrNoAvailableBroker, errors.Join(errs...))
}
