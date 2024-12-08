package natsqueue

import (
	"context"
	"encoding/json"
	"fmt"

	commonDomain "ghiaccio/domain"
	commonPorts "ghiaccio/ports"
	"tusk/internal/config"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type NatsQueue struct {
	conn                 *nats.Conn
	client               nats.JetStreamContext
	analystStream        string
	analystSubject       string
	analystResultSubject string

	encodingStream  string
	encodingSubject string

	analystResultJetstream  *commonPorts.NatsJetstream[commonDomain.AnalystJobResultMessage]
	encodingResultJetstream *commonPorts.NatsJetstream[commonDomain.VideoEncodingResultMessage]
}

func NewNatsQueue(cfg config.NatsConfig) (*NatsQueue, error) {
	conn, err := nats.Connect(cfg.Host)
	if err != nil {
		return nil, err
	}

	js, err := conn.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	analystSubject := fmt.Sprintf("%s.%s", cfg.Streams.Analyst.Name, cfg.JobAnalystSubject)
	analystResultSubject := fmt.Sprintf("%s.%s", cfg.Streams.Analyst.Name, cfg.JobAnalystResultSubject)

	encodingSubject := fmt.Sprintf("%s.%s", cfg.Streams.Encoder.Name, cfg.JobEncoderSubject)
	encodingResultSubject := fmt.Sprintf("%s.%s", cfg.Streams.Encoder.Name, cfg.JobEncoderResultSubject)

	return &NatsQueue{
		conn:                 conn,
		client:               js,
		analystStream:        cfg.Streams.Analyst.Name,
		analystSubject:       analystSubject,
		analystResultSubject: analystResultSubject,

		encodingStream:          cfg.Streams.Encoder.Name,
		encodingSubject:         encodingSubject,
		analystResultJetstream:  commonPorts.NewNatsJetstream[commonDomain.AnalystJobResultMessage](analystResultSubject, cfg.JobAnalystResultConsumer, 1, js),
		encodingResultJetstream: commonPorts.NewNatsJetstream[commonDomain.VideoEncodingResultMessage](encodingResultSubject, cfg.JobAnalystResultConsumer, 1, js),
	}, nil
}

func (nq *NatsQueue) AddAnalystJob(ctx context.Context, job commonDomain.AnalystJobMessage) error {
	msgBytes, err := prepareMessage(ctx, job, nil)
	if err != nil {
		return err
	}

	_, err = nq.client.Publish(nq.analystSubject, *msgBytes)
	if err != nil {
		return err
	}
	return nil
}

func (nq *NatsQueue) AddEncoderJob(ctx context.Context, job commonDomain.VideoEncodingMessage) error {
	msgBytes, err := prepareMessage(ctx, job, nil)
	if err != nil {
		return err
	}

	_, err = nq.client.Publish(nq.encodingSubject, *msgBytes)
	if err != nil {
		return err
	}
	return nil
}

func (nq *NatsQueue) HandleAnalystJobResult(ctx context.Context, handler func(ctx context.Context, job commonDomain.AnalystJobResultMessage) error, errChan chan error) error {
	return nq.analystResultJetstream.GenericHandler(ctx, handler, errChan)
}

func (nq *NatsQueue) HandleEncodingJobResult(ctx context.Context, handler func(ctx context.Context, job commonDomain.VideoEncodingResultMessage) error, errChan chan error) error {
	return nq.encodingResultJetstream.GenericHandler(ctx, handler, errChan)
}

func prepareMessage[Message any](ctx context.Context, message Message, error *string) (*[]byte, error) {

	msg := commonDomain.MessageWrapper[Message]{
		Message: message,
		Err:     error,
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return &msgBytes, nil
}

func (nq *NatsQueue) Close(ctx context.Context) error {
	zap.L().Info("Closing NATs connection")
	nq.conn.Close()
	return nil
}
