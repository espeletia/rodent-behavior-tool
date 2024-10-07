package config

import "strings"

type encodingConfigRaw struct {
	FfmpegPath     string
	FfprobePath    string
	GrMagicPath    string
	MaxClipLength  int
	AllowedClients string
}

type EncodingConfig struct {
	FfmpegPath     string
	FfprobePath    string
	GrMagicPath    string
	MaxClipLength  int
	AllowedClients []string
}

func loadEncodingConfig() EncodingConfig {
	encodingConfig := &encodingConfigRaw{}
	v := configViper("encoding")
	err := v.BindEnv("FfmpegPath", "FFMPEG_PATH")
	if err != nil {
		panic(err)
	}
	err = v.BindEnv("FfprobePath", "FFPROBE_PATH")
	if err != nil {
		panic(err)
	}
	err = v.BindEnv("GrMagicPath", "GR_MAGIC_PATH")
	if err != nil {
		panic(err)
	}
	err = v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = v.Unmarshal(encodingConfig)
	if err != nil {
		panic(err)
	}
	return EncodingConfig{
		FfmpegPath:     encodingConfig.FfmpegPath,
		FfprobePath:    encodingConfig.FfprobePath,
		GrMagicPath:    encodingConfig.GrMagicPath,
		AllowedClients: strings.Split(encodingConfig.AllowedClients, ","),
	}
}
