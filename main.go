package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sta-golang/filefinder/conf"
	"github.com/sta-golang/filefinder/find"
	"github.com/sta-golang/filefinder/out"
	"github.com/sta-golang/filefinder/utils"
	"github.com/sta-golang/go-lib-utils/log"
	tm "github.com/sta-golang/go-lib-utils/time"
)

var help = flag.Bool("help", false, "请求帮助")
var noColor = flag.Bool("color", false, "没有color")

func init() {
	if *noColor {
		conf.NoColor = true
	}
	home := os.Getenv("HOME")
	conf.InitConfig(home + "/.config/filefinder/conf.yaml")
}

func signalHandler(masterCh chan bool) chan bool {
	// 创建一个通道用来接收信号
	sigs := make(chan os.Signal, 1)

	// 注册要监听的信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// 创建一个通道用来接收程序退出通知
	done := make(chan bool, 1)

	// 启动一个 goroutine 来处理信号
	go func() {
		select {
		case <-sigs:
			done <- true
		case <-masterCh:
			done <- true
		}
	}()
	return done
}

func main() {
	flag.Parse()
	if *help || len(os.Args) <= 2 {
		fmt.Print("Maybe you can look at the following example ：")
		fmt.Println("./filefinder test /Users true")
		fmt.Println("\t1. ==> test is frist args (must!) This is keyword")
		fmt.Println("\t2. ==> /Users is second args (not must!) This is find path default : ./")
		fmt.Println("\t3. ==> true is third args (not must!) This is keyword ignore case default true:")
		fmt.Println("example : ", os.Args[0], "<.go[require]> <./[option]> <true[option]>")
		return
	}
	startIndex := 1
	if *noColor {
		conf.NoColor = true
		startIndex = 2
	}
	args := os.Args[utils.MinInt(startIndex, len(os.Args)):]
	rootDir, start := parseArgs(args)
	if !start {
		log.Warn("Unable to start lookup ! ")
		return
	}
	find.Init()
	out.Init()
	masterCh := make(chan bool, 1)
	done := signalHandler(masterCh)
	go func() {
		if err := recover(); err != nil {
			log.Error("panic : ", err)
			masterCh <- true
		}
		log.Debugf("find keyword : \033[3;34m%s\033[0m IgnoreCase : \033[3;34m%v\033[0m from path : \033[3;34m%s\033[0m",
			conf.GetKeyword(), conf.IgnoreCase, rootDir)
		timing := tm.FuncTiming(func() {
			find.Do(rootDir)
		})
		log.Infof("Find finished! timing : %v ms", timing.Milliseconds())
		log.Infof("Find Dir Total : %v File Total : %v You need File Total : %v", conf.FindDirTotal, conf.FindFileTotal, out.ResultSize())
		if out.ResultSize() <= 0 {
			return
		}
		interactive()
		masterCh <- true
	}()
	<-done
	close(done)
	close(masterCh)
}
