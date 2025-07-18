package rabbitmq

import (
	"product-service/config"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Cfg     *config.Config
}

func NewRabbitMQ(cfg *config.Config) (*RabbitMQ, error) {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		err := conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
		Cfg:     cfg,
	}, nil
}

func (r *RabbitMQ) Setup() error {
	// 声明交换机
	err := r.Channel.ExchangeDeclare(
		r.Cfg.ProductExchange,
		"direct", // 普通直连交换机
		true,     // durable
		false,    // auto-delete
		false,    // internal
		false,    // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	// 声明队列
	_, err = r.Channel.QueueDeclare(
		r.Cfg.ProductQueue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	// 绑定队列到交换机
	err = r.Channel.QueueBind(
		r.Cfg.ProductQueue,
		"", // routing key
		r.Cfg.ProductExchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQ) PublishEvent(body []byte) error {
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent, // 持久化消息
		ContentType:  "application/json",
		Body:         body,
		Timestamp:    time.Now(),
	}

	return r.Channel.Publish(
		r.Cfg.ProductExchange,
		"",    // routing key
		false, // mandatory
		false, // immediate
		msg,
	)
}

func (r *RabbitMQ) Close() {
	if r.Channel != nil {
		err := r.Channel.Close()
		if err != nil {
			return
		}
	}
	if r.Conn != nil {
		err := r.Conn.Close()
		if err != nil {
			return
		}
	}
}
