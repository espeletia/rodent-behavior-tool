package natsqueue

import (
	"context"
	"fmt"
	commonPorts "ghiaccio/ports"

	"encoding/json"
	"tusk/internal/config"

	commonDomain "ghiaccio/domain"

	"github.com/nats-io/nats.go"
)

type NatsQueue struct {
	conn           *nats.Conn
	client         nats.JetStreamContext
	analystStream  string
	analystSubject string

	analystJetstream *commonPorts.NatsJetstream[commonDomain.AnalystJobMessage]
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

	return &NatsQueue{
		conn:           conn,
		client:         js,
		analystStream:  cfg.Streams.Analyst.Name,
		analystSubject: analystSubject,

		analystJetstream: commonPorts.NewNatsJetstream[commonDomain.AnalystJobMessage](analystSubject, cfg.JobAnalystConsumer, 1, js),
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
