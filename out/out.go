package out

import (
	"fmt"
	"sync"

	"github.com/sta-golang/filefinder/conf"
	"github.com/sta-golang/filefinder/result"
	"github.com/sta-golang/go-lib-utils/algorithm/data_structure"
)

type outResult struct {
	buff []*result.Result
	mu   sync.Mutex
}

var globalOuts *data_structure.Stack
var currentOuts *outResult
var baseOuts *outResult

var outFmt = "%s %d %v"

func Init() {
	if !conf.NoColor {
		outFmt = "\033[1;36m%s %d\033[0m %v"
	}
	baseOuts = &outResult{
		buff: make([]*result.Result, 0, 8192),
		mu:   sync.Mutex{},
	}
	globalOuts = data_structure.NewStack()
	globalOuts.Push(baseOuts)
	currentOuts = baseOuts
}

func AddOuts() {
	peek := globalOuts.Peek().(*outResult)
	globalOuts.Push(&outResult{
		buff: make([]*result.Result, 0, len(peek.buff)),
		mu:   sync.Mutex{},
	})
	currentOuts = globalOuts.Peek().(*outResult)
}

func Pop() {
	if currentOuts == baseOuts {
		return
	}
	_ = globalOuts.Pop()
	currentOuts = globalOuts.Peek().(*outResult)
}

var baseOut = &outResult{
	buff: make([]*result.Result, 0, 8192),
	mu:   sync.Mutex{},
}

func Put(res *result.Result) {
	put(res, currentOuts)
}

func PutWithStep(res *result.Result) {
}

func OutResult() {
	putBatch(currentOuts.buff)
}

func putBatch(arr []*result.Result) {
	for i := 1; i <= len(arr); i++ {
		fmt.Println(fmt.Sprintf(outFmt, conf.GloabalConfig().LoggerConf.Prefix, i, arr[i-1]))
	}
}

func Get(index int) *result.Result {
	if index < 1 || index > ResultSize() {
		return nil
	}
	return currentOuts.buff[index-1]
}

func ResultSize() int {
	return len(currentOuts.buff)
}

func GetAllResult() []*result.Result {
	return currentOuts.buff
}

func put(res *result.Result, o *outResult) {
	o.mu.Lock()
	defer o.mu.Unlock()
	fmt.Println(fmt.Sprintf(outFmt, conf.GloabalConfig().LoggerConf.Prefix, len(o.buff)+1, res))
	o.buff = append(o.buff, res)
}
