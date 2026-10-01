package transport

import (
	"context"
	"errors"
	"sync"

	v1 "github.com/velonyapp/asset/gen/api/v1"
	"github.com/velonyapp/asset/internal/application/usecase"
	"github.com/velonyapp/asset/internal/conf"
	"go.einride.tech/aip/resourcename"

	"github.com/Azure/go-amqp"
	"github.com/go-kratos/kratos/v3/transport"
	"github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
	"google.golang.org/protobuf/proto"
)

var (
	ErrUnsupportedMessage = errors.New("unsupported message")
	ErrInvalidMessage     = errors.New("invalid message")
)

var _ transport.Server = (*rabbitMQConsumer)(nil)

type rabbitMQConsumer struct {
	c *conf.Transport

	deleteImageHandler *usecase.DeleteImageHandler

	conn      *rabbitmqamqp.AmqpConnection
	consumers map[string]*rabbitmqamqp.Consumer

	wg sync.WaitGroup
}

func NewRabbitMQConsumer(
	c *conf.Transport,
	deleteImageHandler *usecase.DeleteImageHandler,
) *rabbitMQConsumer {
	return &rabbitMQConsumer{
		c:                  c,
		deleteImageHandler: deleteImageHandler,
	}
}

func (rc *rabbitMQConsumer) registerAllConsumers(ctx context.Context) error {
	var errs []error

	if err := rc.registerConsumer(ctx, rc.c.Rabbitmq.Queue.DeleteImage); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (rc *rabbitMQConsumer) handleMessage(ctx context.Context, queue string, data []byte) error {
	switch queue {
	case rc.c.Rabbitmq.Queue.DeleteImage:
		req := new(v1.DeleteImageRequest)

		if err := proto.Unmarshal(data, req); err != nil {
			return ErrInvalidMessage
		}

		var imageID string
		if err := resourcename.Sscan(req.GetName(), "images/{image}", &imageID); err != nil {
			return err
		}

		_, err := rc.deleteImageHandler.Execute(ctx, &usecase.DeleteImage{
			ImageID: imageID,
		})
		return err

	default:
		return ErrUnsupportedMessage
	}
}

func (rc *rabbitMQConsumer) registerConsumer(ctx context.Context, queue string) error {
	consumer, err := rc.conn.NewConsumer(ctx, queue, nil)
	if err != nil {
		return err
	}

	if rc.consumers == nil {
		rc.consumers = make(map[string]*rabbitmqamqp.Consumer)
	}

	rc.consumers[queue] = consumer

	return nil
}

func (rc *rabbitMQConsumer) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := rabbitmqamqp.Dial(
		ctx,
		rc.c.Rabbitmq.Address,
		&rabbitmqamqp.AmqpConnOptions{
			SASLType: amqp.SASLTypePlain(
				rc.c.Rabbitmq.Username,
				rc.c.Rabbitmq.Password,
			),
		},
	)
	if err != nil {
		return err
	}

	rc.conn = conn

	if err := rc.registerAllConsumers(ctx); err != nil {
		return err
	}

	errCh := make(chan error, len(rc.consumers))

	for queue, consumer := range rc.consumers {
		rc.wg.Go(func() {
			for {
				delivery, err := consumer.Receive(ctx)
				if err != nil {
					if ctx.Err() != nil {
						return
					}

					continue
				}

				message := delivery.Message()

				if err := rc.handleMessage(ctx, queue, message.GetData()); err != nil {
					switch {
					case errors.Is(err, ErrUnsupportedMessage),
						errors.Is(err, ErrInvalidMessage):
						delivery.Discard(ctx, nil)

					default:
						if err := delivery.Requeue(ctx); err != nil {
							errCh <- err
							return
						}
					}

					continue
				}

				delivery.Accept(ctx)
			}
		})
	}

	select {
	case <-ctx.Done():
		return nil

	case err := <-errCh:
		return err
	}
}

func (rc *rabbitMQConsumer) Stop(ctx context.Context) error {
	var errs []error

	for name, consumer := range rc.consumers {
		if consumer == nil {
			delete(rc.consumers, name)
			continue
		}

		if err := consumer.Close(ctx); err != nil {
			errs = append(errs, err)
		}

		delete(rc.consumers, name)
	}

	if rc.conn != nil {
		if err := rc.conn.Close(ctx); err != nil {
			errs = append(errs, err)
		}

		rc.conn = nil
	}

	rc.wg.Wait()

	return errors.Join(errs...)
}
