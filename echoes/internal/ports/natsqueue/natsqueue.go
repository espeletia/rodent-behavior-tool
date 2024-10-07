package natsqueue

import (
	"context"
	// "encoding/json"
	"fmt"
	commonPorts "ghiaccio/ports"

	"echoes/internal/config"

	commonDomain "ghiaccio/domain"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type NatsQueue struct {
	conn                 *nats.Conn
	client               nats.JetStreamContext
	streamName           string
	encodingVideoSubject string

	encodingVideoJetstream *commonPorts.NatsJetstream[commonDomain.VideoEncodingMessage]
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
	encodingVideoJobSubject := fmt.Sprintf("%s.%s", cfg.Streams.Encoder.Name, cfg.JobVideoSubject)

	return &NatsQueue{
		conn:                 conn,
		client:               js,
		streamName:           cfg.Streams.Encoder.Name,
		encodingVideoSubject: encodingVideoJobSubject,

		encodingVideoJetstream: commonPorts.NewNatsJetstream[commonDomain.VideoEncodingMessage](encodingVideoJobSubject, cfg.JobVideoConsumer, 1, js),
	}, nil
}

func (nq *NatsQueue) HandleVideoJob(ctx context.Context, handler func(ctx context.Context, job commonDomain.VideoEncodingMessage) error, errChan chan error) error {
	return nq.encodingVideoJetstream.GenericHandler(ctx, handler, errChan)
}

func (nq *NatsQueue) Close(ctx context.Context) error {
	zap.L().Info("Closing NATs connection")
	nq.conn.Close()
	return nil
}
