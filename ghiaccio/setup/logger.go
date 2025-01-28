package setup

import (
	"ghiaccio/config"
	"os"
	"sort"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(configuration config.Config) *zap.Logger {
	cores := []zapcore.Core{}
	if configuration.LoggingConfig.EnableDebugLogger {
		lowCore, highCore := setupDebugLogger(configuration.LoggingConfig)
		cores = append(cores, lowCore)
		cores = append(cores, highCore)
	}
	if configuration.LoggingConfig.EnableFileLogger {
		cores = append(cores, setupFileLogger(configuration.LoggingConfig))
	}
	core := zapcore.NewTee(cores...)
	logger := zap.New(core, buildOptions(configuration)...)
	zap.ReplaceGlobals(logger)
	return logger
}

func setupFileLogger(logConfig config.LoggingConfig) zapcore.Core {
	var level zapcore.Level
	err := level.Set(logConfig.FileLogLevel)
	if err != nil {
		// We want to panic here, since logger is important
		panic(err)
	}

	encoderConfig := ecszap.NewDefaultEncoderConfig()
	if logConfig.FileLogOutput == "stdout" {
		sink, _, err := zap.Open(logConfig.FileLogOutput)
		if err != nil {
			// We want to panic here, since logger is important
			panic(err)
		}
		return ecszap.NewCore(encoderConfig, sink, level)
	} else {
		// rolling file appender
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logConfig.FileLogOutput,
			MaxSize:    100, // megabytes
			MaxBackups: 2,
			MaxAge:     2, // days
		})

		return ecszap.NewCore(encoderConfig, w, level)
	}
}

func setupDebugLogger(logConfig config.LoggingConfig) (zapcore.Core, zapcore.Core) {
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	errorPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > zap.ErrorLevel
	})

	priority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zap.ErrorLevel
	})

	consoleWriter := zapcore.Lock(os.Stdout)
	consoleError := zapcore.Lock(os.Stderr)

	return zapcore.NewCore(consoleEncoder, consoleWriter, priority), zapcore.NewCore(consoleEncoder, consoleError, errorPriority)
}

func buildOptions(cfg config.Config) []zap.Option {
	opts := []zap.Option{}

	if cfg.LoggingConfig.LoggerConfig.Development {
		opts = append(opts, zap.Development())
	}

	if !cfg.LoggingConfig.LoggerConfig.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	stackLevel := zap.ErrorLevel
	if cfg.LoggingConfig.LoggerConfig.Development {
		stackLevel = zap.WarnLevel
	}
	if !cfg.LoggingConfig.LoggerConfig.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	if cfg.LoggingConfig.LoggerConfig.SamplingEnable {
		scfg := cfg.LoggingConfig.LoggerConfig.Sampling
		opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			var samplerOpts []zapcore.SamplerOption
			return zapcore.NewSamplerWithOptions(
				core,
				time.Second,
				scfg.Initial,
				scfg.Thereafter,
				samplerOpts...,
			)
		}))
	}

	if len(cfg.LoggingConfig.LoggerConfig.InitialFields) > 0 {
		fs := make([]zap.Field, 0, len(cfg.LoggingConfig.LoggerConfig.InitialFields))
		keys := make([]string, 0, len(cfg.LoggingConfig.LoggerConfig.InitialFields))
		for k := range cfg.LoggingConfig.LoggerConfig.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, zap.Any(k, cfg.LoggingConfig.LoggerConfig.InitialFields[k]))
		}
		opts = append(opts, zap.Fields(fs...))
	}

	serviceFields := []zapcore.Field{
		zap.String("name", "tusk"),
		zap.String("environment", "tusk"),
	}

	opts = append(opts, zap.Fields(serviceFields...))

	return opts
}
