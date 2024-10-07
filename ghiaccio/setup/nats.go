package setup

import (
	"fmt"

	"ghiaccio/config"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func MigrateNatsStreams(cfg config.NatsConfig) error {
	zap.L().Info("Starting Nats migration")
	conn, err := nats.Connect(cfg.Host)
	if err != nil {
		return err
	}

	js, err := conn.JetStream()
	if err != nil {
		return err
	}

	for _, s := range cfg.Streams {
		zap.L().Info(fmt.Sprintf("Migrating %s stream", s.Name))
		_, err := js.StreamInfo(s.Name)
		if err != nil {
			_, err := js.AddStream(&nats.StreamConfig{
				Name:              s.Name,
				Description:       s.Description,
				Subjects:          s.Subjects,
				Retention:         nats.InterestPolicy,
				MaxMsgs:           s.MaxMsg,
				Discard:           nats.DiscardOld,
				MaxAge:            s.MaxAge,
				MaxMsgsPerSubject: s.MaxMsg,
				Replicas:          int(s.Replicas),
				Duplicates:        s.DuplicateWindow,
			})
			if err != nil {
				return fmt.Errorf("Failed to create stream %s err: %w", s.Name, err)
			}
		} else {
			_, err := js.UpdateStream(&nats.StreamConfig{
				Name:              s.Name,
				Description:       s.Description,
				Subjects:          s.Subjects,
				Retention:         nats.InterestPolicy,
				MaxMsgs:           s.MaxMsg,
				Discard:           nats.DiscardOld,
				MaxAge:            s.MaxAge,
				MaxMsgsPerSubject: s.MaxMsg,
				Replicas:          int(s.Replicas),
			})
			if err != nil {
				return fmt.Errorf("Failed to update stream %s err: %w", s.Name, err)
			}
		}

		for _, c := range s.Consumers {
			filterSubject := ""
			if c.Subject != "" {
				filterSubject = fmt.Sprintf("%s.%s", s.Name, c.Subject)
			}
			deliverPolicy := nats.DeliverLastPolicy
			if c.DeliverPolicy != nil {
				deliverPolicy = nats.DeliverPolicy(*c.DeliverPolicy)
			}
			replayPolicy := nats.ReplayInstantPolicy
			if c.ReplayPolicy != nil {
				replayPolicy = nats.ReplayPolicy(*c.ReplayPolicy)
			}

			_, err := js.ConsumerInfo(s.Name, c.Name)
			if err != nil {
				_, err := js.AddConsumer(s.Name, &nats.ConsumerConfig{
					Durable:       c.Name,
					Name:          c.Name,
					DeliverPolicy: deliverPolicy,
					DeliverGroup:  c.Name,
					MaxDeliver:    int(s.MaxRetry),
					AckPolicy:     nats.AckExplicitPolicy,
					AckWait:       c.AckWait,
					FilterSubject: filterSubject,
					MaxAckPending: int(c.AckPending),
					BackOff:       c.BackOff,
					ReplayPolicy:  replayPolicy,
					Replicas:      int(s.Replicas),
				})
				if err != nil {
					return fmt.Errorf("Failed to create consumer %s %s err: %w", s.Name, c.Name, err)
				}
			} else {
				_, err := js.UpdateConsumer(s.Name, &nats.ConsumerConfig{
					Durable:       c.Name,
					Name:          c.Name,
					DeliverGroup:  c.Name,
					DeliverPolicy: deliverPolicy,
					MaxDeliver:    int(s.MaxRetry),
					AckPolicy:     nats.AckExplicitPolicy,
					AckWait:       c.AckWait,
					FilterSubject: filterSubject,
					MaxAckPending: int(c.AckPending),
					BackOff:       c.BackOff,
					ReplayPolicy:  replayPolicy,
					Replicas:      int(s.Replicas),
				})
				if err != nil {
					return fmt.Errorf("Failed to update consumer %s %s err: %w", s.Name, c.Name, err)
				}

			}
		}
	}
	return nil
}
