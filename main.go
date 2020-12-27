package main

import (
	"github.com/cppdebug/windev"
	"math/rand"
	"os"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	time1 = 700 * time.Millisecond
)

var (
	//加载驱动动态库
	DD, _ = syscall.LoadDLL("DD.64.dll")
	//1：暂停；0：运行
	pause    int32
	skill1Ch = make(chan int, 1)
	ticker1  = time.NewTicker(time1)
)

func main() {
	go ticker()
	go command()
	run()
}

func command() {
	for {
		time.Sleep(time.Duration(randDelay()) * time.Millisecond)
		//开始
		if windev.KeyDownUp(windev.VK_CHARF) == 1 {
			reset()
			atomic.CompareAndSwapInt32(&pause, 0, 1)
		}
		//暂停
		if windev.KeyDownUp(windev.VK_CHARG) == 1 {
			stop()
			atomic.CompareAndSwapInt32(&pause, 1, 0)
		}
		//结束
		if windev.KeyDownUp(windev.VK_F1) == 1 {
			exit()
		}
	}
}

func run() {
	defer exit()
	keyProc, _ := DD.FindProc("DD_key")
	mouseProc, _ := DD.FindProc("DD_btn")
	for {
		time.Sleep(time.Duration(randDelay()) * time.Millisecond)
		if atomic.LoadInt32(&pause) == 0 {
			continue
		}
		select {
		case <-skill1Ch:
			skill1(keyProc)
		default:
			skill2(keyProc, mouseProc)
		}
	}
}

//数字1
func skill1(keyProc *syscall.Proc) {
	keyProc.Call(uintptr(201), 1)
	randDelayMin()
	keyProc.Call(uintptr(201), 2)
	randDelayMin()
}

//shift+鼠标左键
func skill2(keyProc *syscall.Proc, mouseProc *syscall.Proc) {
	keyProc.Call(uintptr(500), 1)
	randDelayMin()
	mouseProc.Call(uintptr(1))
	randDelayMin()
	mouseProc.Call(uintptr(2))
	randDelayMin()
	keyProc.Call(uintptr(500), 2)
	randDelayMin()
}

//定时任务
func ticker() {
	go func() {
		for range ticker1.C {
			skill1Ch <- 0
		}
	}()
}

//重置定时器
func reset() {
	ticker1.Reset(time1)
}

//暂停取消定时器
func stop() {
	ticker1.Stop()
}

//退出释放资源
func exit() {
	stop()
	DD.Release()
	os.Exit(-1)
}

//随机延迟时间
func randDelay() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(80) + 80
}

//随机延迟时间
func randDelayMin() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10) + 10
}
