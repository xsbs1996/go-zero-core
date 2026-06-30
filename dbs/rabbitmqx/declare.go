package rabbitmqx

import amqp "github.com/rabbitmq/amqp091-go"

// declareProducer 声明生产者相关 exchange、queue 和 binding。
func declareProducer(channel *amqp.Channel, conf ProducerConfig) error {
	if err := declareExchange(channel, conf.Exchange); err != nil {
		return err
	}
	if err := declareQueue(channel, conf.Queue); err != nil {
		return err
	}
	return bindQueue(channel, conf.Binding)
}

// declareConsumer 声明消费者相关 exchange、queue、binding 和 qos。
func declareConsumer(channel *amqp.Channel, conf ConsumerConfig) error {
	if err := declareExchange(channel, conf.Exchange); err != nil {
		return err
	}
	if err := declareQueue(channel, conf.Queue); err != nil {
		return err
	}
	if err := bindQueue(channel, conf.Binding); err != nil {
		return err
	}
	if conf.PrefetchCount > 0 || conf.PrefetchSize > 0 {
		return channel.Qos(conf.PrefetchCount, conf.PrefetchSize, conf.GlobalQos)
	}
	return nil
}

func declareExchange(channel *amqp.Channel, conf ExchangeConfig) error {
	if conf.Name == "" {
		return nil
	}
	kind := conf.Kind
	if kind == "" {
		kind = "direct"
	}
	return channel.ExchangeDeclare(conf.Name, kind, conf.Durable, conf.AutoDelete, conf.Internal, conf.NoWait, conf.Args)
}

func declareQueue(channel *amqp.Channel, conf QueueConfig) error {
	if conf.Name == "" {
		return nil
	}
	_, err := channel.QueueDeclare(conf.Name, conf.Durable, conf.AutoDelete, conf.Exclusive, conf.NoWait, conf.Args)
	return err
}

func bindQueue(channel *amqp.Channel, conf BindingConfig) error {
	if conf.Queue == "" || conf.Exchange == "" {
		return nil
	}
	return channel.QueueBind(conf.Queue, conf.RoutingKey, conf.Exchange, conf.NoWait, conf.Args)
}
