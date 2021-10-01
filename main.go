package main

import (
	"github.com/sta-golang/go-lib-utils/log"
)

func init() {
	logger := log.NewConsoleLog(log.DEBUG, "😊FileFinder")
	log.SetGlobalLogger(logger)
}

func main() {
	log.Info("Finder File Task Begin")
}
