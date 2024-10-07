package config

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

type Policy nats.DeliverPolicy

func (p *Policy) UnmarshalText(data []byte) error {
	switch string(data) {
	case string("all"), string("undefined"):
		*p = Policy(nats.DeliverAllPolicy)
	case string("last"):
		*p = Policy(nats.DeliverLastPolicy)
	case string("new"):
		*p = Policy(nats.DeliverNewPolicy)
	case string("by_start_sequence"):
		*p = Policy(nats.DeliverByStartSequencePolicy)
	case string("by_start_time"):
		*p = Policy(nats.DeliverByStartTimePolicy)
	case string("last_per_subject"):
		*p = Policy(nats.DeliverLastPerSubjectPolicy)
	}

	return nil
}

func (p Policy) MarshalText() ([]byte, error) {
	switch p {
	case Policy(nats.DeliverAllPolicy):
		return json.Marshal("all")
	case Policy(nats.DeliverLastPolicy):
		return json.Marshal("last")
	case Policy(nats.DeliverNewPolicy):
		return json.Marshal("new")
	case Policy(nats.DeliverByStartSequencePolicy):
		return json.Marshal("by_start_sequence")
	case Policy(nats.DeliverByStartTimePolicy):
		return json.Marshal("by_start_time")
	case Policy(nats.DeliverLastPerSubjectPolicy):
		return json.Marshal("last_per_subject")
	default:
		return nil, fmt.Errorf("nats: unknown deliver policy %v", p)
	}
}

type ReplyPolicy nats.ReplayPolicy

func (p *ReplyPolicy) UnmarshalText(data []byte) error {
	switch string(data) {
	case string("instant"):
		*p = ReplyPolicy(nats.ReplayInstantPolicy)
	case string("original"):
		*p = ReplyPolicy(nats.ReplayOriginalPolicy)
	default:
		*p = ReplyPolicy(nats.ReplayInstantPolicy)
	}

	return nil
}

func (p ReplyPolicy) MarshalText() ([]byte, error) {
	switch p {
	case ReplyPolicy(nats.DeliverAllPolicy):
		return json.Marshal("instant")
	case ReplyPolicy(nats.DeliverLastPolicy):
		return json.Marshal("original")
	default:
		return nil, fmt.Errorf("nats: unknown deliver policy %v", p)
	}
}

type NatsConfig struct {
	Host    string
	Streams map[string]StreamConfig
}

type StreamConfig struct {
	Name            string
	Subjects        []string
	Description     string
	Replicas        int32
	MaxMsg          int64
	MaxRetry        int32
	MaxAge          time.Duration
	AckWait         time.Duration
	Consumers       []ConsumerConfig
	DuplicateWindow time.Duration
}

type ConsumerConfig struct {
	Name         string
	AckWait      time.Duration
	AckPending   int32
	Subject      string
	BackOff      []time.Duration
	ReplayPolicy *ReplyPolicy
	// TODO maybe change to something more generic
	DeliverPolicy *Policy
}

func LoadNatsConfig(serviceName string) NatsConfig {
	natsConfig := NatsConfig{}
	u := configViper("nats", serviceName)
	err := u.BindEnv("Host", "NATS_URL")
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	err = u.ReadInConfig()
	if err != nil {
		// We can ignore error since nats is optional config
		return natsConfig
	}
	err = u.Unmarshal(&natsConfig, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			mapstructure.TextUnmarshallerHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
		),
	))
	if err != nil {
		// We can ignore error since nats is optional config
		return natsConfig
	}
	return natsConfig
}
