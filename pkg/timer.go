package pkg

import (
	"fmt"
	"time"
)

type timerCtxValue string

const (
	TimerCtxValue timerCtxValue = "timerCtxValue"
)

type timerInstance struct {
	startTime time.Time
}

func StartNewTimer() *timerInstance {
	return &timerInstance{
		startTime: time.Now(),
	}
}

func (ti *timerInstance) Since() time.Duration {
	return time.Since(ti.startTime)
}

func (ti *timerInstance) SinceStringInMS() string {
	return fmt.Sprintf("%vms", float64(ti.Since())/float64(time.Millisecond))
}
