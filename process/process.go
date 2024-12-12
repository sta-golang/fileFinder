package process

import (
	"github.com/sta-golang/filefinder/conf"
	"github.com/sta-golang/filefinder/find"
	"github.com/sta-golang/filefinder/out"
)

func Search(keyword string) {
	conf.SetKeyword(keyword)
	res := out.GetAllResult()
	out.AddOuts()
	find.Search(res, keyword)
}

func Undo() bool {
	if !conf.UndoKeyword() {
		return false
	}
	res := out.GetAllResult()
	out.AddOuts()
	find.Search(res, conf.GetKeyword())
	return true
}

func Redo() bool {
	if !conf.RedoKeyword() {
		return false
	}
	out.Pop()
	out.OutResult()
	return true
}
