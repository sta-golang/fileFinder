package find

import (
	"strconv"
	"strings"

	"github.com/sta-golang/filefinder/conf"
	"github.com/sta-golang/filefinder/out"
	"github.com/sta-golang/filefinder/result"
	"github.com/sta-golang/filefinder/utils"
	"github.com/sta-golang/go-lib-utils/async/asyncgroup"
)

func Search(keyword string) {
	conf.HistoryKeyword = append(conf.HistoryKeyword, keyword)
	conf.KEYWORD = keyword
	out.InitStep()
	arr := out.GetAllResult()
	limit := len(arr) / conf.GNum
	if limit == 0 {
		search(arr)
		return
	}
	ag := asyncgroup.New(asyncgroup.WithTaskSize(conf.GNum))
	defer ag.Shutdown()
	for i := 0; i*limit < len(arr); i++ {
		start := i * limit
		currentArr := arr[start:utils.MinInt(len(arr), start+limit)]
		_ = ag.Add(strconv.Itoa(i), func() (interface{}, error) {
			search(currentArr)
			return nil, nil
		})
	}
	ag.Wait()
	conf.Step += 1
}

func search(arr []*result.Result) {
	for _, res := range arr {
		filename := res.DirPath + res.FileName
		if index := strings.Index(filename, conf.KEYWORD); index != -1 {
			out.Put(result.New(filename))
		}
	}
}
