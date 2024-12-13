package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/sta-golang/filefinder/conf"
	"github.com/sta-golang/filefinder/out"
	"github.com/sta-golang/filefinder/process"
	"github.com/sta-golang/go-lib-utils/cmd"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/str"
	"github.com/sta-golang/go-lib-utils/time"
)

func parseArgs(args []string) (string, bool) {
	if len(args) <= 0 {
		return "", false
	}
	if len(args) > 2 {
		initIgnoreCase(args[2])
	}
	keyword := args[0]
	if conf.IgnoreCase {
		keyword = strings.ToLower(keyword)
	}
	conf.SetKeyword(keyword)
	currentDir := "./"
	command, err := cmd.ExecCmd("pwd")
	if err != nil {
		log.Warnf("bash pwd have err : :%v", err)
	} else {
		currentDir = str.BytesToString(command.OutMessage)
		currentDir = str.Trim(currentDir)
	}
	if len(args) == 1 {
		return currentDir, true
	}
	if !checkDirPath(args[1]) {
		log.Errorf("path : %s is not dir", args[1])
		return "", false
	}
	return args[1], true
}

func initIgnoreCase(ig string) {
	if ig == "0" || ig == "false" || ig == "no" || ig == "n" {
		conf.IgnoreCase = false
	}
}

func checkDirPath(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		log.Errorf("path : %s have error : %v", path, err)
		os.Exit(1)
	}
	return stat.IsDir()
}

func doInteractive() bool {
	msgFmt := "%s\t%s"
	if !conf.NoColor {
		msgFmt = "%s\t\033[1;2;32m%s\033[0m"
	}
	stepMsg := fmt.Sprintf(msgFmt, conf.GloabalConfig().LoggerConf.Prefix,
		fmt.Sprintf("Please Input : (1-%d) or n (next search) or r (redo) or u (undo) or q (quit)", out.ResultSize()))
	fmt.Println(stepMsg)
	var command string
	if _, err := fmt.Scanf("%s", &command); err != nil {
		return true
	}
	if command == "quit" || command == "q" || command == "Q" || command == "Quit" || command == "exit" {
		return false
	}
	if command == "n" || command == "s" || command == "next" || command == "search" {
		keyword := ""
		for {
			fmt.Println(fmt.Sprintf(msgFmt, conf.GloabalConfig().LoggerConf.Prefix, "Please Input Next Search Keyword : "))
			if _, err := fmt.Scanf("%s", &keyword); err != nil || keyword == "" {
				continue
			}
			process.Search(keyword)
			break
		}
	}
	if command == "u" || command == "undo" {
		if !process.Undo() {
			stepMsg := fmt.Sprintf(msgFmt, conf.GloabalConfig().LoggerConf.Prefix,
				fmt.Sprintf("no undo result"))
			fmt.Println(stepMsg)
			return true
		}
	}
	if command == "r" || command == "redo" {
		if !process.Redo() {
			stepMsg := fmt.Sprintf(msgFmt, conf.GloabalConfig().LoggerConf.Prefix,
				fmt.Sprintf("no redo result"))
			fmt.Println(stepMsg)
			return true
		}
	}
	return true
}

func showErrDirInfo() {
	if conf.GloabalConfig().NotShowWarn {
		return
	}
	msgFmt := "%s %s"
	if !conf.NoColor {
		msgFmt = "%s \033[1;2;38m%s\033[0m"
	}
	fileLogMsg := ""
	if conf.GloabalConfig().LoggerConf.FileLogConf != nil {
		time.GetNowDateStr()
		logName := "sta"
		if conf.GloabalConfig().LoggerConf.FileLogConf.FileName != "" {
			logName = conf.GloabalConfig().LoggerConf.FileLogConf.FileName
		}
		fileLogMsg = fmt.Sprintf("detail please view :  %s/%s.log.%s.*", conf.GloabalConfig().LoggerConf.FileLogConf.FileDir, logName, time.GetNowDateStr())
	}
	if conf.ErrDirTotal > 0 {
		warnMsg := fmt.Sprintf(msgFmt, conf.GloabalConfig().LoggerConf.Prefix,
			fmt.Sprintf(" ==> You has %d search Dir err %s", conf.ErrDirTotal, fileLogMsg))
		fmt.Println(warnMsg)
	}
}

func interactive() {
	showErrDirInfo()
	for {
		if !doInteractive() {
			break
		}
	}
}
