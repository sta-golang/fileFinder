package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sta-golang/filefinder/conf"
	"github.com/sta-golang/filefinder/find"
	"github.com/sta-golang/filefinder/out"
	"github.com/sta-golang/filefinder/utils"
	"github.com/sta-golang/go-lib-utils/log"
	tm "github.com/sta-golang/go-lib-utils/time"
)

var help = flag.Bool("help", false, "请求帮助")

func init() {
	home := os.Getenv("HOME")
	conf.InitConfig(home + "/.config/filefinder/conf.yaml")
}

func main() {
	flag.Parse()
	if *help {
		fmt.Print("Maybe you can look at the following example ：")
		fmt.Println("./filefinder test /Users true")
		fmt.Println("\t1. ==> test is frist args (must!) This is keyword")
		fmt.Println("\t2. ==> /Users is second args (not must!) This is find path default : ./")
		fmt.Println("\t3. ==> true is third args (not must!) This is keyword ignore case default true:")
		return
	}
	args := os.Args[utils.MinInt(1, len(os.Args)):]
	rootDir, start := parseArgs(args)
	if !start {
		log.Warn("Unable to start lookup ! ")
		return
	}
	log.Infof("find keyword : \033[3;34m%s\033[0m IgnoreCase : \033[3;34m%v\033[0m from path : \033[3;34m%s\033[0m",
		conf.KEYWORD, conf.IgnoreCase, rootDir)
	timing := tm.FuncTiming(func() {
		find.Do(rootDir)
	})
	conf.Step = 1
	log.Infof("Find finished! timing : %v ms", timing.Milliseconds())
	log.Infof("Find Dir Total : %v File Total : %v You need File Total : %v", conf.FindDirTotal, conf.FindFileTotal, out.ResultSize())
	if out.ResultSize() <= 0 {
		return
	}
	interactive()
}
