package out

import (
	"fmt"
	"sync"

	"github.com/sta-golang/filefinder/conf"
	"github.com/sta-golang/filefinder/result"
)

type outResult struct {
	buff []*result.Result
	mu   sync.Mutex
	cnt  int
}

var outFmt = "%s %d %v"

func init() {
	if !conf.NoColor {
		outFmt = "\033[1;36m%s %d\033[0m %v"
	}
}

var baseOut = &outResult{
	buff: make([]*result.Result, 0, 8192),
	mu:   sync.Mutex{},
	cnt:  1,
}

var stepOut = []*outResult{
	baseOut,
}

func Put(res *result.Result) {
	put(res, stepOut[conf.Step])
}

func InitStep() {
	for conf.Step >= len(stepOut) {
		stepOut = append(stepOut, &outResult{
			buff: make([]*result.Result, 0, len(stepOut[len(stepOut)-1].buff)),
			cnt:  1,
			mu:   sync.Mutex{},
		})
	}
}

func PutWithStep(res *result.Result) {
}

func OutResult() {
	putBatch(stepOut[conf.Step-1].buff)
}

func putBatch(arr []*result.Result) {
	for i := 1; i <= len(arr); i++ {
		fmt.Println(fmt.Sprintf(outFmt, conf.GloabalConfig().LoggerConf.Prefix, i, arr[i-1]))
	}
}

func Get(index int) *result.Result {
	if index < 1 || index > len(stepOut[conf.Step-1].buff) {
		return nil
	}
	return stepOut[conf.Step-1].buff[index-1]
}

func ResultSize() int {
	return len(stepOut[conf.Step-1].buff)
}

func GetAllResult() []*result.Result {
	return stepOut[conf.Step-1].buff
}

func put(res *result.Result, o *outResult) {
	o.mu.Lock()
	defer o.mu.Unlock()
	fmt.Println(fmt.Sprintf(outFmt, conf.GloabalConfig().LoggerConf.Prefix, o.cnt, res))
	o.buff = append(o.buff, res)
	o.cnt++

}
