package xkafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	ErrNilReader         = errors.New("xkafka: reader is nil")          // ErrNilReader 表示 Kafka reader 为空。
	ErrNilConsumeHandler = errors.New("xkafka: consume handler is nil") // ErrNilConsumeHandler 表示消费处理函数为空。
	ErrNilBatchHandler   = errors.New("xkafka: batch handler is nil")   // ErrNilBatchHandler 表示批量消费处理函数为空。
	ErrInvalidBatchSize  = errors.New("xkafka: invalid batch size")     // ErrInvalidBatchSize 表示批量大小非法。
	ErrInvalidBatchTime  = errors.New("xkafka: invalid batch timeout")  // ErrInvalidBatchTime 表示批量等待时间非法。
)

// readerConfig 构建 Kafka 消费者配置。
func readerConfig(conf Config, topic, group string) *kafka.ReaderConfig {
	conf = conf.WithDefault()
	return &kafka.ReaderConfig{
		Brokers:         conf.Brokers,
		Topic:           topic,
		GroupID:         group,
		Dialer:          &kafka.Dialer{ClientID: conf.ClientID, Timeout: conf.DialTimeout},
		ReadLagInterval: -1,
		MaxWait:         conf.ReadTimeout,
	}
}

// NewConsumer 根据 topic、group 和配置创建 Kafka 消费者。
//
// 参数：
//   - topic: Kafka topic。
//   - group: Kafka consumer group。
//   - conf: Kafka 连接配置。
//   - opts: 可选消费者配置，例如 WithReaderConfig、WithoutConsumerPing。
//
// 返回值：
//   - *kafka.Reader: Kafka 消费者。
//   - error: 配置非法、topic/group 为空或 ping 失败时返回错误。
func NewConsumer(topic, group string, conf Config, opts ...ConsumerOption) (*kafka.Reader, error) {
	if err := conf.ValidateConsumer(); err != nil {
		return nil, err
	}

	topic, err := normalizeTopic(topic)
	if err != nil {
		return nil, err
	}
	group, err = normalizeGroup(group)
	if err != nil {
		return nil, err
	}

	conf = conf.WithDefault()
	options := newConsumerOptions(conf, topic, group, opts...)
	return kafka.NewReader(*options.readerConfig), nil
}

// MustNewConsumer 根据 topic、group 和配置创建 Kafka 消费者，失败时直接 panic。
//
// 参数：
//   - topic: Kafka topic。
//   - group: Kafka consumer group。
//   - conf: Kafka 连接配置。
//   - opts: 可选消费者配置。
//
// 返回值：
//   - *kafka.Reader: Kafka 消费者。
func MustNewConsumer(topic, group string, conf Config, opts ...ConsumerOption) *kafka.Reader {
	reader, err := NewConsumer(topic, group, conf, opts...)
	if err != nil {
		panic(err)
	}
	return reader
}

// closeConsumer 关闭 Kafka 消费者。
func closeConsumer(reader *kafka.Reader) error {
	if reader == nil {
		return ErrNilConsumer
	}
	if err := reader.Close(); err != nil {
		return fmt.Errorf("close kafka consumer failed: %w", err)
	}
	return nil
}

// readMessageLoop 使用 Kafka reader 执行单条消息读取循环。
//
// handler 在每条消息处理完成后执行，只有 handler 成功后才会提交 offset。
func readMessageLoop(ctx context.Context, reader *kafka.Reader, handler func(context.Context, kafka.Message) error) error {
	if reader == nil {
		return ErrNilReader
	}
	if handler == nil {
		return ErrNilConsumeHandler
	}

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			return err
		}

		if err := handler(ctx, msg); err != nil {
			return err
		}

		if err := reader.CommitMessages(ctx, msg); err != nil {
			return fmt.Errorf("commit kafka message failed: %w", err)
		}
	}
}

// readBatchLoop 使用 Kafka reader 执行批量消息读取循环。
//
// batchSize 表示单批次最大消息数，达到后立即触发 handler。
// batchTimeout 表示批次最长等待时间，从收到第一条消息开始计时。
// handler 在每批消息处理完成后执行，只有 handler 成功后才会提交 offset。
func readBatchLoop(ctx context.Context, reader *kafka.Reader, batchSize int, batchTimeout time.Duration, handler func(context.Context, []kafka.Message) error) error {
	if reader == nil {
		return ErrNilReader
	}
	if handler == nil {
		return ErrNilBatchHandler
	}
	if batchSize <= 0 {
		return ErrInvalidBatchSize
	}
	if batchTimeout <= 0 {
		return ErrInvalidBatchTime
	}

	batch := make([]kafka.Message, 0, batchSize)
	var deadline time.Time

	flush := func() error {
		if len(batch) == 0 {
			return nil
		}
		if err := handler(ctx, batch); err != nil {
			return err
		}
		if err := reader.CommitMessages(ctx, batch...); err != nil {
			return fmt.Errorf("commit Kafka messages failed: %w", err)
		}
		batch = batch[:0]
		deadline = time.Time{}
		return nil
	}

	for {
		if err := ctx.Err(); err != nil {
			if len(batch) > 0 {
				if flushErr := flush(); flushErr != nil {
					return flushErr
				}
			}
			return err
		}

		fetchCtx := ctx
		var cancel context.CancelFunc

		if len(batch) > 0 {
			wait := time.Until(deadline)
			if wait <= 0 {
				if err := flush(); err != nil {
					return err
				}
				continue
			}

			fetchCtx, cancel = context.WithTimeout(ctx, wait)
		}

		msg, err := reader.FetchMessage(fetchCtx)
		if cancel != nil {
			cancel()
		}
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				if len(batch) > 0 {
					if flushErr := flush(); flushErr != nil {
						return flushErr
					}
				}
				continue
			}
			if errors.Is(err, context.Canceled) {
				if len(batch) > 0 {
					if flushErr := flush(); flushErr != nil {
						return flushErr
					}
				}
				return err
			}
			return err
		}

		if len(batch) == 0 {
			deadline = time.Now().Add(batchTimeout)
		}

		batch = append(batch, msg)
		if len(batch) >= batchSize {
			if err := flush(); err != nil {
				return err
			}
		}
	}
}
