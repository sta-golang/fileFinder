package conf

import (
	"io/ioutil"

	"github.com/sta-golang/go-lib-utils/codec"
	"github.com/sta-golang/go-lib-utils/log"
)

type Config struct {
	LoggerConf LogConfig `yaml:"log"`
}

var globalConfig Config

func GloabalConfig() *Config {
	return &globalConfig
}

type LogConfig struct {
	Prefix string `yaml:"prefix"`
	Level  int    `yaml:"Level"`
}

var defaultConfig = func() Config {
	return Config{
		LoggerConf: LogConfig{
			Prefix: "【fileFinder】",
			Level:  int(log.INFO),
		},
	}
}

func InitConfig(confPath string) {
	var config Config
	bys, err := ioutil.ReadFile(confPath)
	if err != nil {
		config = defaultConfig()
		initLog(config)
		return
	}
	err = codec.API.YamlAPI.UnMarshal(bys, &config)
	if err != nil {
		config = defaultConfig()
		initLog(config)
		return
	}
	globalConfig = config
	initLog(config)
}

func initLog(config Config) {
	if config.LoggerConf.Level < int(log.DEBUG) {
		config.LoggerConf.Level = int(log.DEBUG)
	}
	if config.LoggerConf.Level > int(log.FATAL) {
		config.LoggerConf.Level = int(log.FATAL)
	}
	logger := log.NewConsoleLog(log.Level(config.LoggerConf.Level), config.LoggerConf.Prefix)
	log.ConsoleLogger = logger
	log.SetGlobalLogger(logger)
}
