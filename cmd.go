package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sta-golang/filefinder/conf"
	"github.com/sta-golang/filefinder/find"
	"github.com/sta-golang/filefinder/out"
	"github.com/sta-golang/filefinder/result"
	"github.com/sta-golang/go-lib-utils/cmd"
	"github.com/sta-golang/go-lib-utils/log"
	"github.com/sta-golang/go-lib-utils/str"
)

func parseArgs(args []string) (string, bool) {
	if len(args) <= 0 {
		return "", false
	}
	if len(args) > 2 {
		initIgnoreCase(args[2])
	}
	conf.KEYWORD = args[0]
	if conf.IgnoreCase {
		conf.KEYWORD = strings.ToLower(conf.KEYWORD)
	}
	conf.HistoryKeyword = append(conf.HistoryKeyword, conf.KEYWORD)
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

func interactive() {
	msgFmt := "%s\t%s"
	if !conf.NoColor {
		msgFmt = "%s\t\033[1;2;32m%s\033[0m"
	}
	for {
		stepMsg := fmt.Sprintf(msgFmt, conf.GloabalConfig().LoggerConf.Prefix,
			fmt.Sprintf("Please Input : (1-%d) or n (next search) or r (redo) or q (quit)", out.ResultSize()))
		fmt.Println(stepMsg)
		var command string
		if _, err := fmt.Scanf("%s", &command); err != nil {
			continue
		}
		if command == "quit" || command == "q" || command == "Q" || command == "Quit" || command == "exit" {
			break
		}
		if command == "n" || command == "s" || command == "next" || command == "search" {
			keyword := ""
			for {
				fmt.Println(fmt.Sprintf(msgFmt, conf.GloabalConfig().LoggerConf.Prefix, "Please Input Next Search Keyword : "))
				if _, err := fmt.Scanf("%s", &keyword); err != nil || keyword == "" {
					continue
				}
				find.Search(keyword)
				break
			}
		}
		if command == "r" || command == "redo" {
			conf.Step -= 1
			if conf.Step <= 0 {
				conf.Step = 1
			}
			out.OutResult()
		}

		index, err := strconv.Atoi(command)
		if err != nil {
			continue
		}
		res := out.Get(index)
		if res == nil {
			continue
		}
		if nextStep(res) {
			break
		}
	}
}

func nextStep(res *result.Result) bool {
	fmt.Println(conf.GloabalConfig().LoggerConf.Prefix, " current File :", res)
	msgFmt := "%s\t%s"
	if !conf.NoColor {
		msgFmt = "%s\t\033[1;2;32m%s\033[0m"
	}
	nextStepMsg := fmt.Sprintf(msgFmt, conf.GloabalConfig().LoggerConf.Prefix,
		fmt.Sprintf("Please Input : p (print) or r (reselect) or q (quit)"))
	for {
		fmt.Println(nextStepMsg)
		var command string
		if _, err := fmt.Scanf("%s", &command); err != nil {
			continue
		}
		if command == "r" {
			return false
		}
		if command == "quit" || command == "q" || command == "Q" || command == "Quit" || command == "exit" {
			break
		}
		if command == "o" || command == "e" || command == "p" {
			fmt.Println(res.DirPath + res.FileName)
			break
		}
	}
	return true
}
