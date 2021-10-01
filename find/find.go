package find

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"strings"
	"sync/atomic"
	"time"

	"github.com/sta-golang/filefinder/conf"
	"github.com/sta-golang/filefinder/out"
	"github.com/sta-golang/filefinder/result"
	"github.com/sta-golang/go-lib-utils/async/asyncgroup"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/pool/workerpool"
	"github.com/sta-golang/go-lib-utils/str"
)

func Do(rootDir string) error {
	infos, err := ioutil.ReadDir(rootDir)
	if err != nil {
		log.Error(err)
		return err
	}
	wp := workerpool.NewWithQueueSize(conf.GNum, 8192<<4)
	defer wp.Stop()
	ag := asyncgroup.New(asyncgroup.WithConcurrentSecurity(),
		asyncgroup.WithTaskSize(8192), asyncgroup.WithWorkPool(wp))
	defer ag.Shutdown()
	_ = ag.Add(getID(), func() (interface{}, error) {
		do(ag, infos, rootDir)
		return nil, nil
	})
	ag.Wait()
	return nil
}

func do(ag *asyncgroup.Group, infos []fs.FileInfo, parentDir string) {
	for _, info := range infos {
		filename := parentDir + "/" + info.Name()
		tempFilename := filename
		if conf.IgnoreCase {
			tempFilename = strings.ToLower(filename)
		}
		if index := strings.Index(tempFilename, conf.KEYWORD); index != -1 {
			out.Put(result.New(filename))
		}
		if info.IsDir() {
			atomic.AddInt32(&conf.FindDirTotal, 1)
			dirInfos, err := ioutil.ReadDir(filename)
			if err != nil {
				log.Error(err)
				continue
			}
			currentFilename := filename
			if err := ag.Add(getID(), func() (interface{}, error) {
				do(ag, dirInfos, currentFilename)
				return nil, nil
			}); err != nil {
				log.Error(err)
			}
			continue
		}
		atomic.AddInt32(&conf.FindFileTotal, 1)
	}
}

func getID() string {
	return fmt.Sprint(time.Now().UnixNano()) + str.XID()
}
