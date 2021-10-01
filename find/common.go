package find

import (
	"github.com/sta-golang/filefinder/conf"
	"github.com/sta-golang/go-lib-utils/pool/workerpool"
)

var wp workerpool.Executor

func Init() {
	wp = workerpool.NewWithQueueSize(conf.GloabalConfig().WorkerConf.GNum, conf.GloabalConfig().WorkerConf.ChSize)
}
