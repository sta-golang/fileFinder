package conf

import (
	"fmt"
	"runtime"
)

type keywordInfo struct {
	index   int
	history []string
	size    int
}

var globalKeywordInfo keywordInfo

var IgnoreCase = true
var FindDirTotal = int32(0)
var FindFileTotal = int32(0)
var GNum = 32

var NoColor bool

func init() {
	NoColor = (runtime.GOOS == "windows")
	GNum = runtime.NumCPU() * 2
	globalKeywordInfo = keywordInfo{
		index:   0,
		history: make([]string, 4),
	}
}

func (k *keywordInfo) Add(keyword string) {
	fmt.Println("key word : ", keyword)
	if k.index >= cap(k.history) {
		newArr := make([]string, len(k.history)*2)
		for i := 0; i < len(k.history); i++ {
			newArr[i] = k.history[i]
		}
		k.history = newArr
	}
	k.history[k.index] = keyword
	k.size++
	k.index += 1
}

func (k *keywordInfo) GetKeyword() string {
	return k.history[k.index-1]
}
func (k *keywordInfo) AddIndex() bool {
	if k.size <= k.index {
		return false
	}
	k.index += 1
	return true
}
func (k *keywordInfo) SubIndex() bool {
	if k.index <= 1 {
		return false
	}
	k.index -= 1
	return true
}

func SetKeyword(keyword string) {
	globalKeywordInfo.Add(keyword)
}

func RedoKeyword() bool {
	return globalKeywordInfo.SubIndex()
}

func UndoKeyword() bool {
	return globalKeywordInfo.AddIndex()
}

func Debug() {
	fmt.Println(globalKeywordInfo)
}

func GetKeyword() string {
	return globalKeywordInfo.GetKeyword()
}

const (
	Byte = int64(1)
	KB   = Byte * 1024
	MB   = KB * 1024
	GB   = MB * 1024
	TB   = GB * 1024
)
