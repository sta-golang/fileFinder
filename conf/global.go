package conf

import "runtime"

var KEYWORD = ""
var HistoryKeyword = []string{}
var IgnoreCase = true
var FindDirTotal = int32(0)
var FindFileTotal = int32(0)
var GNum = 32
var Step = 0

var NoColor bool

func init() {
	NoColor = (runtime.GOOS == "windows")
}

const (
	Byte = int64(1)
	KB   = Byte * 1024
	MB   = KB * 1024
	GB   = MB * 1024
	TB   = GB * 1024
)
