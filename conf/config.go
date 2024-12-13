package conf

import (
	"io/ioutil"
	"runtime"

	"github.com/sta-golang/go-lib-utils/codec"
	"github.com/sta-golang/go-lib-utils/log"
)

type Config struct {
	WorkerConf  WorkerpoolConfig `yaml:"worker"`
	LoggerConf  LogConfig        `yaml:"log"`
	NotShowWarn bool             `yaml:"not_show_warn"`
}

type WorkerpoolConfig struct {
	GNum   int `yaml:"gnum"`
	ChSize int `yaml:"ch_size"`
}

var globalConfig Config

func GloabalConfig() *Config {
	return &globalConfig
}

type LogConfig struct {
	Prefix      string             `yaml:"prefix"`
	Level       int                `yaml:"Level"`
	FileLogConf *log.FileLogConfig `yaml:"file_log"`
}

var defaultConfig = func() Config {
	return Config{
		LoggerConf: LogConfig{
			Prefix: "【fileFinder】",
			Level:  int(log.INFO),
		},
		WorkerConf: WorkerpoolConfig{
			GNum:   (runtime.NumCPU() << 1) + 1,
			ChSize: 8192 << 4,
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
	initWorker(&globalConfig)
}

func initWorker(config *Config) {
	if config.WorkerConf.GNum == 0 {
		config.WorkerConf.GNum = (runtime.NumCPU() << 1) + 1
	}
	if config.WorkerConf.ChSize == 0 {
		config.WorkerConf.ChSize = 8192 << 5
	}
	if config.WorkerConf.GNum <= 1 {
		config.WorkerConf.GNum = 4
	}
	if config.WorkerConf.ChSize < 8192 {
		config.WorkerConf.ChSize = 8192
	}
}

func initLog(config Config) {
	if config.LoggerConf.FileLogConf != nil {
		logger := log.NewFileLog(config.LoggerConf.FileLogConf)
		log.SetGlobalLogger(logger)
		return
	}
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
