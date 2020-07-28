package game

import (
	"container/list"
	"time"

	"github.com/yenkeia/yams/game/cm"
)

type delayedAction struct {
	delayedType cm.DelayedType
	actionTime  time.Time
	callback    delayedCallback
}

type delayedCallback func()

type actionList struct {
	ls *list.List
}

func newActionList() *actionList {
	return &actionList{
		ls: list.New(),
	}
}

func (al *actionList) pushDelayAction(typ cm.DelayedType, delay int, cb delayedCallback) {
	al.ls.PushBack(&delayedAction{
		delayedType: typ,
		actionTime:  time.Now().Add(time.Millisecond * time.Duration(delay)),
		callback:    cb,
	})
}

func (al *actionList) execute() {
	now := time.Now()
	for it := al.ls.Front(); it != nil; {
		action := it.Value.(*delayedAction)
		if now.Before(action.actionTime) {
			it = it.Next()
			continue
		}
		action.callback()
		tmp := it
		it = it.Next()
		al.ls.Remove(tmp)
	}
}