package config

import "fmt"

type LoggingConfig struct {
	EnableDebugLogger bool
	EnableFileLogger  bool
	FileLogLevel      string
	FileLogOutput     string
	LoggerConfig      LoggerConfig
}

type LoggerConfig struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Sampling          SamplingConfig
	SamplingEnable    bool
	InitialFields     map[string]interface{}
}

type SamplingConfig struct {
	Initial    int
	Thereafter int
}

func loadLoggingConfig(serviceName string) LoggingConfig {
	loggingConfig := &LoggingConfig{}
	v := configViper("logging", serviceName)
	err := v.BindEnv("EnableFileLogger", "ENABLE_FILE_LOGGER")
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	err = v.BindEnv("EnableDebugLogger", "ENABLE_DEBUG_LOGGER")
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	err = v.BindEnv("FileLogOutput", "FILE_LOG_OUTPUT")
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	err = v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	err = v.Unmarshal(loggingConfig)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	return *loggingConfig
}
