package event

/*
定时事件

每个定时事件都对应一个Executor

CHECK: 检查锁相关实现是否正确
OPT: 实际上更符合信号量语义
	每个时刻只能有一个正在运行的executor goroutine，可以看作一种资源。
	Stop时，释放资源，Start时，申请资源。
	使用信号量重写。 
*/

import (
	"sync"

	"github.com/robfig/cron"
)

type TimerExecutor struct {
	CronString string
	// 定时器
	cron *cron.Cron
	// 关闭Executor的通道
	closeCh chan struct{}
	// 是否正在运行
	isRunning bool
	// 保护isRunning的锁
	runningLock sync.Mutex
}

func (e *TimerExecutor)Init(CronString string, Func func()) {
	e.cron = cron.New()
	e.cron.AddFunc(CronString, Func)
	e.CronString = CronString
	e.closeCh = make(chan struct{})
}

func (e *TimerExecutor)execute() {

	e.runningLock.Lock()
	if e.isRunning {
		e.runningLock.Unlock()
		return
	}
	e.isRunning = true
	e.runningLock.Unlock()

	e.cron.Start()
	<-e.closeCh
	e.cron.Stop()

	e.runningLock.Lock()
	e.isRunning = false
	e.runningLock.Unlock()
}

func (e *TimerExecutor)Start() {
	go e.execute()
}

func (e *TimerExecutor)Stop() {
	e.runningLock.Lock()
	defer e.runningLock.Unlock()
	if !e.isRunning {
		return
	}
	e.closeCh <- struct{}{}
}

func (e *TimerExecutor)IsRunning() bool {
	return e.isRunning
}
