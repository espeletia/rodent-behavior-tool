package config

import (
	"ghiaccio/config"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"time"
)

type NatsConfig struct {
	Host           string
	MaxRetry       int
	AckWait        time.Duration
	RequestTimeout time.Duration
	MaxMsgs        int64
	MaxAge         time.Duration
	MaxAckPending  int

	JobVideoConsumer string
	JobVideoSubject  string

	Streams StreamsConfig
}

type StreamsConfig struct {
	Encoder config.StreamConfig
}

func loadNatsConfig() NatsConfig {
	natsConfig := &NatsConfig{}
	v := configViper("nats")
	err := v.BindEnv("Host", "NATS_URL", "DEMETER_NATS_URL")
	if err != nil {
		panic(err)
	}
	err = v.BindEnv("StreamName", "NATS_STREAM_NAME")
	if err != nil {
		panic(err)
	}
	err = v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = v.Unmarshal(&natsConfig, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			mapstructure.TextUnmarshallerHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
		),
	))
	if err != nil {
		panic(err)
	}
	return *natsConfig
}
