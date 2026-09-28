package transport

import (
	"context"
	"errors"

	v1 "github.com/velonyapp/asset/gen/api/v1"
	"github.com/velonyapp/asset/internal/conf"
	"github.com/velonyapp/asset/internal/presentation/api"

	"github.com/Azure/go-amqp"
	"github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
	"google.golang.org/protobuf/proto"
)

var (
	ErrUnsupportedMessage = errors.New("unsupported message")
	ErrInvalidMessage     = errors.New("invalid message")
)

type RabbitMQConsumer struct {
	c   *conf.Transport
	svc *api.Service

	conn     *rabbitmqamqp.AmqpConnection
	consumer *rabbitmqamqp.Consumer
}

func NewRabbitMQConsumer(
	c *conf.Transport,
	svc *api.Service,
) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		c:   c,
		svc: svc,
	}
}

func (s *RabbitMQConsumer) Start(ctx context.Context) error {
	conn, err := rabbitmqamqp.Dial(ctx,
		s.c.Rabbitmq.Address,
		&rabbitmqamqp.AmqpConnOptions{
			SASLType: amqp.SASLTypePlain(
				s.c.Rabbitmq.Username,
				s.c.Rabbitmq.Password,
			),
		},
	)
	if err != nil {
		return err
	}

	s.conn = conn

	consumer, err := conn.NewConsumer(ctx, s.c.Rabbitmq.Image.Queue, nil)
	if err != nil {
		_ = conn.Close(ctx)
		s.conn = nil
		return err
	}

	s.consumer = consumer

	for {
		delivery, err := consumer.Receive(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || ctx.Err() != nil {
				return nil
			}

			return err
		}

		message := delivery.Message()

		routingKey, ok := message.Annotations["x-routing-key"].(string)
		if !ok {
			if err := delivery.Discard(ctx, nil); err != nil {
				return err
			}

			continue
		}

		if err := s.handle(ctx, routingKey, message.GetData()); err != nil {
			switch {
			case errors.Is(err, ErrUnsupportedMessage):
				if err := delivery.Discard(ctx, nil); err != nil {
					return err
				}

			case errors.Is(err, ErrInvalidMessage):
				if err := delivery.Discard(ctx, nil); err != nil {
					return err
				}

			default:
				if err := delivery.Requeue(ctx); err != nil {
					return err
				}
			}

			continue
		}

		if err := delivery.Accept(ctx); err != nil {
			return err
		}
	}
}

func (s *RabbitMQConsumer) handle(ctx context.Context, routingKey string, data []byte) error {
	switch routingKey {
	case s.c.Rabbitmq.Image.RemoveRoutingKey:
		req := new(v1.RemoveImageRequest)

		if err := proto.Unmarshal(data, req); err != nil {
			return ErrInvalidMessage
		}

		_, err := s.svc.RemoveImage(ctx, req)
		return err

	default:
		return ErrUnsupportedMessage
	}
}

func (s *RabbitMQConsumer) Stop(ctx context.Context) error {
	if s.consumer != nil {
		if err := s.consumer.Close(ctx); err != nil {
			return err
		}

		s.consumer = nil
	}

	if s.conn != nil {
		if err := s.conn.Close(ctx); err != nil {
			return err
		}

		s.conn = nil
	}

	return nil
}
