package conf

import "github.com/sta-golang/go-lib-utils/log"

type Config struct {
	LoggerConf LogConfig `yaml:"log"`
}

type LogConfig struct {
	Prefix string `yaml:"prefix"`
	Level  int    `yaml:"Level"`
}

var defaultConfig = func() *Config {
	return &Config{
		LoggerConf: LogConfig{
			Prefix: "【fileFinder】",
			Level:  int(log.INFO),
		},
	}
}

func InitConfig(confPath string) {

}
