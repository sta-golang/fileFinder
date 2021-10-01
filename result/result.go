package result

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/sta-golang/filefinder/conf"
	tm "github.com/sta-golang/go-lib-utils/time"
)

var NoColor bool

func init() {
	NoColor = (runtime.GOOS == "windows")
}

type Result struct {
	IsDir          bool
	DirPath        string
	FileName       string
	FileSize       int64
	FileChangeTime time.Time
	KeywordIndex   []int
}

var New = func(filename string) *Result {
	ret := &Result{}
	stat, _ := os.Stat(filename)
	ret.IsDir = stat.IsDir()
	ret.FileSize = stat.Size()
	ret.FileChangeTime = stat.ModTime()
	ret.DirPath = filename
	if !ret.IsDir {
		ret.FileName = stat.Name()
		index := strings.LastIndex(filename, ret.FileName)
		ret.DirPath = filename[:index]
	}
	index := -1
	lastIndex := -len(conf.KEYWORD)
	if conf.IgnoreCase {
		filename = strings.ToLower(filename)
	}
	for {
		index = strings.Index(filename, conf.KEYWORD)
		if index == -1 {
			break
		}
		filename = filename[index+len(conf.KEYWORD):]
		ret.KeywordIndex = append(ret.KeywordIndex, index+lastIndex+len(conf.KEYWORD))
		lastIndex = lastIndex + index + len(conf.KEYWORD)
	}
	return ret
}

func (r *Result) String() string {
	var logFormat = "%-120s\t%-4s\t%-20s\t%s"
	filename := r.FileDirPath()
	if !NoColor {
		textFmt := "\033[38m%s\033[0m"
		if r.IsDir {
			textFmt = "\033[35m%s\033[0m"
		}
		keywordFmt := "\033[1;3;4;31m%s\033[0m"
		buff := bytes.Buffer{}
		lastIndex := 0
		for _, index := range r.KeywordIndex {
			buff.WriteString(fmt.Sprintf(textFmt,
				filename[lastIndex:index]) + fmt.Sprintf(keywordFmt, filename[index:index+len(conf.KEYWORD)]))
			lastIndex = index + len(conf.KEYWORD)
		}
		buff.WriteString(fmt.Sprintf(textFmt, filename[lastIndex:]))
		filename = buff.String()
	}
	return fmt.Sprintf(logFormat, filename, r.FileType(), r.Size(), r.ModTime())
}

func (r *Result) FileDirPath() string {
	return r.DirPath + r.FileName
}

func (r *Result) FileType() string {
	if r.IsDir {
		return "\033[1;38mDir\033[0m"
	}
	return "\033[1;38mFile\033[0m"
}

func (r *Result) ModTime() string {
	return fmt.Sprintf("\033[36m%s\033[0m", tm.ParseDataTimeToStr(r.FileChangeTime))
}

func (r *Result) Size() string {
	unit := conf.Byte
	unitStr := "Byte"
	if r.FileSize > 5*conf.KB {
		unit = conf.KB
		unitStr = "KB"
	}
	if r.FileSize > 10*conf.MB {
		unit = conf.MB
		unitStr = "MB"
	}
	if r.FileSize > 1*conf.GB {
		unit = conf.GB
		unitStr = "GB"
	}
	if r.FileSize > conf.TB {
		unit = conf.TB
		unitStr = "TB"
	}
	size := 0.0
	if unit >= conf.GB {
		temp := r.FileSize / 1000
		size = float64(temp) / float64(unit)
		size *= 1000
		return fmt.Sprintf("\033[1;38m%.2f %-10s\033[0m", size, unitStr)
	}
	return fmt.Sprintf("\033[1;38m%.2f %-5s\033[0m", float64(r.FileSize)/float64(unit), unitStr)
}
