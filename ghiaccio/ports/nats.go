package ports

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"ghiaccio/domain"
	"strings"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v3/json"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func NewNatsJetstream[T any](subject string, group string, batchSize int, client nats.JetStreamContext) *NatsJetstream[T] {
	return &NatsJetstream[T]{
		subject:     subject,
		group:       group,
		client:      client,
		batchSize:   batchSize,
		ExitChan:    make(chan bool, 1),
		wg:          sync.WaitGroup{},
		isEphemeral: false,
	}
}

type NatsJetstream[T any] struct {
	subject     string
	group       string
	maxRetry    int
	batchSize   int
	client      nats.JetStreamContext
	wg          sync.WaitGroup
	ExitChan    chan bool
	Sub         *nats.Subscription
	isEphemeral bool
}

func (nj *NatsJetstream[T]) SetEphemeral() *NatsJetstream[T] {
	nj.isEphemeral = true
	return nj
}

func (nj *NatsJetstream[T]) GenericBatchHandler(subscriberCtx context.Context, handler func(ctx context.Context, message []T) error, errChan chan error) error {
	durable := nj.group
	options := []nats.SubOpt{
		nats.Context(subscriberCtx),
	}
	if nj.isEphemeral {
		durable = ""
		options = append(options, nats.Bind(strings.Split(nj.subject, ".")[0], nj.group))
	}

	sub, err := nj.client.PullSubscribe(nj.subject, durable, options...)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to connect to stream %s %s", nj.subject, nj.group), zap.Error(err))
		return err
	}
	nj.Sub = sub
	defer func() {
		if sub != nil && sub.IsValid() {
			err := sub.Unsubscribe()
			if err != nil {
				zap.L().Error("Error unsubscribing", zap.Error(err))
			}
		}
	}()

	for {
		select {
		case <-subscriberCtx.Done():
			return fmt.Errorf("Context deadline exceeded")
		case <-nj.ExitChan:
			return nil
		default:
			break
		}

		pending, _, err := nj.Sub.Pending()
		if err != nil {
			return nil
		}

		if !nj.Sub.IsValid() && pending <= 0 {
			return nil
		}

		msgs, err := sub.Fetch(nj.batchSize)
		if err != nil {
			if errors.Is(err, nats.ErrConnectionClosed) || errors.Is(err, nats.ErrConsumerDeleted) {
				return err
			}
			if !errors.Is(err, nats.ErrTimeout) && !errors.Is(err, nats.ErrBadSubscription) && !strings.Contains(err.Error(), "Exceeded MaxWaiting") {
				zap.L().Info(err.Error())
			}
			continue
		}

		nj.wg.Add(1)
		ctx := context.Background()

		batch := make([]T, len(msgs))
		for i, m := range msgs {
			_, parsedMessage, err := prepareHandlerMessage[T](ctx, m)
			if err != nil {
				zap.L().Error(err.Error())
				nakErr := m.Nak()
				if nakErr != nil {
					zap.L().Error(nakErr.Error())
				}
				errChan <- err
				continue
			}

			batch[i] = parsedMessage
		}
		err = handler(ctx, batch)

		for _, m := range msgs {
			if err != nil {
				zap.L().Error(err.Error())
				nakErr := m.Nak()
				if nakErr != nil {
					zap.L().Error(nakErr.Error())
				}
				errChan <- err
				continue
			}

			err = m.AckSync()
			if err != nil {
				zap.L().Error(err.Error())
				nakErr := m.Nak()
				if nakErr != nil {
					zap.L().Error(nakErr.Error())
				}
				errChan <- err
			}
		}
		nj.wg.Done()
	}
}

func (nj *NatsJetstream[T]) GenericHandler(ctx context.Context, handler func(ctx context.Context, message T) error, errChan chan error) error {
	durable := nj.group
	options := []nats.SubOpt{
		nats.Context(ctx),
	}
	if nj.isEphemeral {
		durable = ""
		options = append(options, nats.Bind(strings.Split(nj.subject, ".")[0], nj.group))
	}

	sub, err := nj.client.PullSubscribe(nj.subject, durable, options...)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	nj.Sub = sub
	defer func() {
		if sub != nil && sub.IsValid() {
			err := sub.Unsubscribe()
			if err != nil {
				zap.L().Error("Error unsubscribing", zap.Error(err))
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("Context deadline exceeded")
		case <-nj.ExitChan:
			return nil
		default:
			break
		}

		pending, _, err := nj.Sub.Pending()
		if err != nil {
			return nil
		}

		if !nj.Sub.IsValid() && pending <= 0 {
			return nil
		}

		msgs, err := sub.Fetch(nj.batchSize)
		if err != nil {
			if errors.Is(err, nats.ErrConnectionClosed) || errors.Is(err, nats.ErrConsumerDeleted) {
				return err
			}
			if !errors.Is(err, nats.ErrTimeout) && !errors.Is(err, nats.ErrBadSubscription) && !strings.Contains(err.Error(), "Exceeded MaxWaiting") {
				zap.L().Info(err.Error())
			}
			continue
		}
		for _, m := range msgs {
			nj.wg.Add(1)
			processGeneric(m, nj.subject, handler, errChan)
			nj.wg.Done()
		}
	}
}

func (nj *NatsJetstream[T]) GenericAsyncHandler(ctx context.Context, handler func(ctx context.Context, jobResult T) error, errChan chan error) error {
	durable := nj.group
	options := []nats.SubOpt{
		nats.Context(ctx),
	}
	if nj.isEphemeral {
		durable = ""
		options = append(options, nats.Bind(strings.Split(nj.subject, ".")[0], nj.group))
	}

	sub, err := nj.client.PullSubscribe(nj.subject, durable, options...)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	nj.Sub = sub
	defer func() {
		if sub != nil && sub.IsValid() {
			err := sub.Unsubscribe()
			if err != nil {
				zap.L().Error("Error unsubscribing", zap.Error(err))
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("Context deadline exceeded")
		case <-nj.ExitChan:
			return nil
		default:
			break
		}

		pending, _, err := nj.Sub.Pending()
		if err != nil {
			return nil
		}

		if !nj.Sub.IsValid() && pending <= 0 {
			return nil
		}

		msgs, err := sub.Fetch(nj.batchSize)
		if err != nil {
			if errors.Is(err, nats.ErrConnectionClosed) || errors.Is(err, nats.ErrConsumerDeleted) {
				return err
			}
			if !errors.Is(err, nats.ErrTimeout) && !errors.Is(err, nats.ErrBadSubscription) && !strings.Contains(err.Error(), "Exceeded MaxWaiting") {
				zap.L().Info(err.Error())
			}
			continue
		}
		nj.wg.Add(1)
		handlerWaitGroup := sync.WaitGroup{}

		for _, m := range msgs {
			m := m
			go func() {
				handlerWaitGroup.Add(1)
				processGeneric(m, nj.subject, handler, errChan)
				handlerWaitGroup.Done()
			}()
		}
		handlerWaitGroup.Wait()
		nj.wg.Done()
	}
}

func prepareHandlerMessage[T any](ctx context.Context, m *nats.Msg) (context.Context, T, error) {
	var result T
	msg := domain.MessageWrapper[T]{}
	decoder := json.NewDecoder(bytes.NewBuffer(m.Data))
	decoder.SetNumberType(json.UnmarshalIntOrFloat)
	err := decoder.Decode(&msg)
	if err != nil {
		return ctx, result, err
	}
	return ctx, msg.Message, nil
}

func processGeneric[T any](m *nats.Msg, subedSubject string, handler func(ctx context.Context, jobResult T) error, errChan chan error) {
	ctx := context.Background()
	zap.L().Info(fmt.Sprintf("Received message %v on %v", m.Subject, subedSubject))

	msg := domain.MessageWrapper[T]{}
	decoder := json.NewDecoder(bytes.NewBuffer(m.Data))
	decoder.SetNumberType(json.UnmarshalIntOrFloat)
	err := decoder.Decode(&msg)
	if err != nil {
		zap.L().Error("Failed decoding message", zap.String("Subject", subedSubject), zap.String("Message.Subject", m.Subject), zap.Error(err))
		nakErr := m.Nak()
		if nakErr != nil {
			zap.L().Error(nakErr.Error())
		}
		errChan <- err
		return
	}
	zap.L().Info("message", zap.Any("msg", msg), zap.Any("raw data", m.Data))

	err = handler(ctx, msg.Message)
	if err != nil {
		zap.L().Error(err.Error())
		nakErr := m.Nak()
		if nakErr != nil {
			zap.L().Error(nakErr.Error())
		}
		errChan <- err
		return
	}
	attempt := 0
	for {
		err = m.AckSync()
		if err != nil {
			attempt += 1
			if attempt >= 5 {
				return
			}

			err := fmt.Errorf("Failed to ack msg %w", err)
			zap.L().Error(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		return
	}
}

func (nj *NatsJetstream[T]) Close() error {
	zap.L().Info("Closing Jetstream queue")
	if nj.Sub == nil {
		return nil
	}

	err := nj.Sub.Drain()
	if err != nil {
		return err
	}

	nj.wg.Wait()
	return nil
}
