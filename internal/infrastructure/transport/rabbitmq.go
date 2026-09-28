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

type rabbitMQConsumer struct {
	c   *conf.Transport
	svc *api.Service

	conn      *rabbitmqamqp.AmqpConnection
	consumers map[string]*rabbitmqamqp.Consumer
}

func NewRabbitMQConsumer(
	c *conf.Transport,
	svc *api.Service,
) *rabbitMQConsumer {
	return &rabbitMQConsumer{
		c:   c,
		svc: svc,
	}
}

func (s *rabbitMQConsumer) registerConsumers(ctx context.Context) error {
	var errs []error

	if err := s.registerConsumer(ctx, s.c.Rabbitmq.Queue.RemoveImage); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (s *rabbitMQConsumer) handleMessage(ctx context.Context, queue string, data []byte) error {
	switch queue {
	case s.c.Rabbitmq.Queue.RemoveImage:
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

func (s *rabbitMQConsumer) registerConsumer(ctx context.Context, queue string) error {
	consumer, err := s.conn.NewConsumer(ctx, queue, nil)
	if err != nil {
		return err
	}

	if s.consumers == nil {
		s.consumers = make(map[string]*rabbitmqamqp.Consumer)
	}

	s.consumers[queue] = consumer

	return nil
}

func (s *rabbitMQConsumer) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := rabbitmqamqp.Dial(
		ctx,
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

	if err := s.registerConsumers(ctx); err != nil {
		return err
	}

	errCh := make(chan error, len(s.consumers))

	for queue, consumer := range s.consumers {
		go func() {
			for {
				delivery, err := consumer.Receive(ctx)
				if err != nil {
					if errors.Is(err, context.Canceled) || ctx.Err() != nil {
						errCh <- nil
						return
					}

					errCh <- err
					return
				}

				message := delivery.Message()

				if err := s.handleMessage(ctx, queue, message.GetData()); err != nil {
					switch {
					case errors.Is(err, ErrUnsupportedMessage),
						errors.Is(err, ErrInvalidMessage):
						if err := delivery.Discard(ctx, nil); err != nil {
							errCh <- err
							return
						}

					default:
						if err := delivery.Requeue(ctx); err != nil {
							errCh <- err
							return
						}
					}

					continue
				}

				if err := delivery.Accept(ctx); err != nil {
					errCh <- err
					return
				}
			}
		}()
	}

	select {
	case <-ctx.Done():
		return nil

	case err := <-errCh:
		return err
	}
}

func (s *rabbitMQConsumer) Stop(ctx context.Context) error {
	var errs []error

	for name, consumer := range s.consumers {
		if consumer == nil {
			delete(s.consumers, name)
			continue
		}

		if err := consumer.Close(ctx); err != nil {
			errs = append(errs, err)
		}

		delete(s.consumers, name)
	}

	if s.conn != nil {
		if err := s.conn.Close(ctx); err != nil {
			errs = append(errs, err)
		}

		s.conn = nil
	}

	return errors.Join(errs...)
}
