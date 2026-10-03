package transport

import (
	"context"
	"errors"
	"sync"

	v1 "github.com/velonyapp/asset/gen/api/v1"
	"github.com/velonyapp/asset/internal/application/command"
	"github.com/velonyapp/asset/internal/application/port"
	"github.com/velonyapp/asset/internal/conf"

	"github.com/Azure/go-amqp"
	"github.com/go-kratos/kratos/v3/transport"
	"github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
	"go.einride.tech/aip/resourcename"
	"google.golang.org/protobuf/proto"
)

var (
	ErrUnsupportedMessage = errors.New("unsupported message")
	ErrInvalidMessage     = errors.New("invalid message")
)

const (
	imageResourcePattern = "images/{image}"
)

var _ transport.Server = (*RabbitMQConsumer)(nil)

type RabbitMQConsumer struct {
	c *conf.Transport

	conn      *rabbitmqamqp.AmqpConnection
	consumers map[string]*rabbitmqamqp.Consumer

	wg sync.WaitGroup
}

func NewRabbitMQConsumer(
	_ *command.HandlerRegistry,
	c *conf.Transport,
) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		c: c,
	}
}

func (rc *RabbitMQConsumer) registerAllConsumers(ctx context.Context) error {
	var errs []error

	if err := rc.registerConsumer(ctx, rc.c.Rabbitmq.Queues.CreateImage); err != nil {
		errs = append(errs, err)
	}
	if err := rc.registerConsumer(ctx, rc.c.Rabbitmq.Queues.ProcessImage); err != nil {
		errs = append(errs, err)
	}
	if err := rc.registerConsumer(ctx, rc.c.Rabbitmq.Queues.ReconcileImage); err != nil {
		errs = append(errs, err)
	}
	if err := rc.registerConsumer(ctx, rc.c.Rabbitmq.Queues.DeleteImage); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (rc *RabbitMQConsumer) handleMessage(ctx context.Context, queue string, data []byte) error {
	switch queue {
	case rc.c.Rabbitmq.Queues.CreateImage:
		req := new(v1.CreateImageRequest)

		if err := proto.Unmarshal(data, req); err != nil {
			return ErrInvalidMessage
		}

		if _, err := command.Send(ctx, &command.CreateImage{
			Tags:      req.Image.Tags,
			ObjectKey: req.Image.ObjectKey,
		}); err != nil {
			return err
		}

	case rc.c.Rabbitmq.Queues.ProcessImage:
		req := new(v1.ProcessImageRequest)

		if err := proto.Unmarshal(data, req); err != nil {
			return ErrInvalidMessage
		}

		var imageID string
		if err := resourcename.Sscan(req.GetName(), imageResourcePattern, &imageID); err != nil {
			return err
		}

		var resize *port.ImageResize
		if req.Resize != nil {
			var fit port.ImageResizeFit
			switch req.Resize.Fit {
			case v1.ImageResizeFit_IMAGE_RESIZE_FIT_UNSPECIFIED:
			case v1.ImageResizeFit_IMAGE_RESIZE_FIT_CONTAIN:
				fit = port.ImageResizeFitContain
			case v1.ImageResizeFit_IMAGE_RESIZE_FIT_COVER:
				fit = port.ImageResizeFitCover
			case v1.ImageResizeFit_IMAGE_RESIZE_FIT_PAD:
				fit = port.ImageResizeFitPad
			case v1.ImageResizeFit_IMAGE_RESIZE_FIT_STRETCH:
				fit = port.ImageResizeFitStretch
			default:
				return port.ErrUnsupportedResizeFit
			}

			var gravity port.ImageGravity
			switch req.Resize.Gravity {
			case v1.ImageGravity_IMAGE_GRAVITY_UNSPECIFIED:
			case v1.ImageGravity_IMAGE_GRAVITY_CENTER:
				gravity = port.ImageGravityCenter
			case v1.ImageGravity_IMAGE_GRAVITY_TOP:
				gravity = port.ImageGravityTop
			case v1.ImageGravity_IMAGE_GRAVITY_TOP_RIGHT:
				gravity = port.ImageGravityTopRight
			case v1.ImageGravity_IMAGE_GRAVITY_RIGHT:
				gravity = port.ImageGravityRight
			case v1.ImageGravity_IMAGE_GRAVITY_BOTTOM_RIGHT:
				gravity = port.ImageGravityBottomRight
			case v1.ImageGravity_IMAGE_GRAVITY_BOTTOM:
				gravity = port.ImageGravityBottom
			case v1.ImageGravity_IMAGE_GRAVITY_BOTTOM_LEFT:
				gravity = port.ImageGravityBottomLeft
			case v1.ImageGravity_IMAGE_GRAVITY_LEFT:
				gravity = port.ImageGravityLeft
			case v1.ImageGravity_IMAGE_GRAVITY_TOP_LEFT:
				gravity = port.ImageGravityTopLeft
			default:
				return port.ErrUnsupportedImageGravity
			}

			resize = &port.ImageResize{
				Width:           req.Resize.Width,
				Height:          req.Resize.Height,
				Fit:             fit,
				Gravity:         gravity,
				BackgroundColor: req.Resize.BackgroundColor,
				AllowUpscale:    req.Resize.AllowUpscale,
			}
		}

		var encoding *port.ImageEncoding
		if req.Encoding != nil {
			var format port.ImageFormat
			switch req.Encoding.Format {
			case v1.ImageFormat_IMAGE_FORMAT_UNSPECIFIED:
			case v1.ImageFormat_IMAGE_FORMAT_JPEG:
				format = port.ImageFormatJPEG
			case v1.ImageFormat_IMAGE_FORMAT_PNG:
				format = port.ImageFormatPNG
			case v1.ImageFormat_IMAGE_FORMAT_WEBP:
				format = port.ImageFormatWebP
			case v1.ImageFormat_IMAGE_FORMAT_AVIF:
				format = port.ImageFormatAVIF
			default:
				return port.ErrUnsupportedImageFormat
			}

			encoding = &port.ImageEncoding{
				Format:  format,
				Quality: req.Encoding.Quality,
			}
		}

		if _, err := command.Send(ctx, &command.ProcessImage{
			ImageID: imageID,
			Options: port.ImageProcessOptions{
				Resize:         resize,
				Encoding:       encoding,
				AutoRotate:     req.AutoRotate,
				RemoveMetadata: req.RemoveMetadata,
			},
		}); err != nil {
			return err
		}

	case rc.c.Rabbitmq.Queues.ReconcileImage:
		req := new(v1.ReconcileImageRequest)

		if err := proto.Unmarshal(data, req); err != nil {
			return ErrInvalidMessage
		}

		var imageID string
		if err := resourcename.Sscan(req.GetName(), imageResourcePattern, &imageID); err != nil {
			return err
		}

		if _, err := command.Send(ctx, &command.ReconcileImage{
			ImageID: imageID,
		}); err != nil {
			return err
		}

	case rc.c.Rabbitmq.Queues.DeleteImage:
		req := new(v1.DeleteImageRequest)

		if err := proto.Unmarshal(data, req); err != nil {
			return ErrInvalidMessage
		}

		var imageID string
		if err := resourcename.Sscan(req.GetName(), imageResourcePattern, &imageID); err != nil {
			return err
		}

		if _, err := command.Send(ctx, &command.DeleteImage{
			ImageID: imageID,
		}); err != nil {
			return err
		}

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
