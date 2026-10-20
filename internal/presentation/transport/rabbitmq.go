package transport

import (
	"context"
	"errors"
	"sync"

	v1 "github.com/velonyapp/asset/gen/job/v1"
	"github.com/velonyapp/asset/internal/conf"
	"github.com/velonyapp/asset/internal/presentation/service"

	"github.com/Azure/go-amqp"
	"github.com/go-kratos/kratos/v3/transport"
	"github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	ErrUnsupportedMessage = errors.New("unsupported message")
	ErrInvalidMessage     = errors.New("invalid message")
)

const (
	createImageQueue  = "asset.image.create"
	confirmImageQueue = "asset.image.confirm"
	processImageQueue = "asset.image.process"
	deleteImageQueue  = "asset.image.delete"
)

var _ transport.Server = (*RabbitMQConsumer)(nil)

type RabbitMQConsumer struct {
	c          *conf.Transport
	jobService *service.Job

	conn      *rabbitmqamqp.AmqpConnection
	consumers map[string]*rabbitmqamqp.Consumer

	wg sync.WaitGroup
}

func NewRabbitMQConsumer(
	c *conf.Transport,
	jobService *service.Job,
) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		c:          c,
		jobService: jobService,
	}
}

func (rc *RabbitMQConsumer) registerAllConsumers(ctx context.Context) error {
	var errs []error

	if err := rc.registerConsumer(ctx, createImageQueue); err != nil {
		errs = append(errs, err)
	}
	if err := rc.registerConsumer(ctx, confirmImageQueue); err != nil {
		errs = append(errs, err)
	}
	if err := rc.registerConsumer(ctx, processImageQueue); err != nil {
		errs = append(errs, err)
	}
	if err := rc.registerConsumer(ctx, deleteImageQueue); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (rc *RabbitMQConsumer) handleMessage(ctx context.Context, queue string, data []byte) error {
	switch queue {
	case createImageQueue:
		payload := new(v1.CreateImagePayload)

		if err := protojson.Unmarshal(data, payload); err != nil {
			return ErrInvalidMessage
		}

		return rc.jobService.CreateImage(ctx, payload)

	case confirmImageQueue:
		payload := new(v1.ConfirmImagePayload)

		if err := protojson.Unmarshal(data, payload); err != nil {
			return ErrInvalidMessage
		}

		return rc.jobService.ConfirmImage(ctx, payload)

	case processImageQueue:
		payload := new(v1.ProcessImagePayload)

		if err := protojson.Unmarshal(data, payload); err != nil {
			return ErrInvalidMessage
		}

		return rc.jobService.ProcessImage(ctx, payload)

	case deleteImageQueue:
		payload := new(v1.DeleteImagePayload)

		if err := protojson.Unmarshal(data, payload); err != nil {
			return ErrInvalidMessage
		}

		return rc.jobService.DeleteImage(ctx, payload)

	default:
		return ErrUnsupportedMessage
	}
}

func (rc *RabbitMQConsumer) registerConsumer(ctx context.Context, queue string) error {
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

func (rc *RabbitMQConsumer) Start(ctx context.Context) error {
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
						if err := delivery.RequeueWithAnnotationsAndDeliveryFailed(ctx, nil, true); err != nil {
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

func (rc *RabbitMQConsumer) Stop(ctx context.Context) error {
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
