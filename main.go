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
	// 2 20 5 21 7
	time1 = 1333 * time.Millisecond
	time2 = 20333 * time.Millisecond
	time3 = 5333 * time.Millisecond
	time5 = 21333 * time.Millisecond
	time6 = 7333 * time.Millisecond
)

var (
	//加载驱动动态库
	DD, _ = syscall.LoadDLL("DD.64.dll")
	//1：暂停；0：运行
	pause    int32
	skill1Ch = make(chan int, 1)
	skill2Ch = make(chan int, 1)
	skill3Ch = make(chan int, 1)
	skill5Ch = make(chan int, 1)
	skill6Ch = make(chan int, 1)
	ticker1  = time.NewTicker(time1)
	ticker2  = time.NewTicker(time2)
	ticker3  = time.NewTicker(time3)
	ticker5  = time.NewTicker(time5)
	ticker6  = time.NewTicker(time6)
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
		if windev.KeyDownUp(windev.VK_CHARZ) == 1 {
			reset()
			atomic.CompareAndSwapInt32(&pause, 0, 1)
		}
		//暂停
		if windev.KeyDownUp(windev.VK_CHARX) == 1 {
			stop()
			atomic.CompareAndSwapInt32(&pause, 1, 0)
		}
		//结束
		if windev.KeyDownUp(windev.VK_CHARB) == 1 {
			exit()
		}
	}
}

func run() {
	defer exit()
	keyProc, _ := DD.FindProc("DD_key")
	mouseProc, _ := DD.FindProc("DD_btn")
	skill5Ch <- 0
	for {
		time.Sleep(time.Duration(randDelay()) * time.Millisecond)
		if atomic.LoadInt32(&pause) == 0 {
			continue
		}
		select {
		case <-skill1Ch:
			skill1(keyProc)
		case <-skill2Ch:
			skill2(keyProc)
		case <-skill3Ch:
			skill3(keyProc)
		case <-skill5Ch:
			skill5(keyProc, mouseProc)
		case <-skill6Ch:
			skill6(keyProc)
		default:
			skill4(keyProc)
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

//数字2
func skill2(keyProc *syscall.Proc) {
	keyProc.Call(uintptr(202), 1)
	randDelayMin()
	keyProc.Call(uintptr(202), 2)
	randDelayMin()
}

//数字3
func skill3(keyProc *syscall.Proc) {
	keyProc.Call(uintptr(203), 1)
	randDelayMin()
	keyProc.Call(uintptr(203), 2)
	randDelayMin()
}

//数字4
func skill4(keyProc *syscall.Proc) {
	keyProc.Call(uintptr(204), 1)
	randDelayMin()
	keyProc.Call(uintptr(204), 2)
	randDelayMin()
}

//shift+鼠标左键
func skill5(keyProc *syscall.Proc, mouseProc *syscall.Proc) {
	keyProc.Call(uintptr(500), 1)
	randDelayMin()
	mouseProc.Call(uintptr(1))
	randDelayMin()
	mouseProc.Call(uintptr(2))
	randDelayMin()
	keyProc.Call(uintptr(500), 2)
	randDelayMin()
}

//鼠标右键
func skill6(mouseProc *syscall.Proc) {
	mouseProc.Call(uintptr(4))
	randDelayMin()
	mouseProc.Call(uintptr(8))
	randDelayMin()
}

//定时任务
func ticker() {
	go func() {
		for range ticker1.C {
			skill1Ch <- 0
		}
	}()
	go func() {
		for range ticker2.C {
			skill2Ch <- 0
		}
	}()

	go func() {
		for range ticker3.C {
			skill3Ch <- 0
		}
	}()

	go func() {
		for range ticker5.C {
			skill5Ch <- 0
		}
	}()
	go func() {
		for range ticker6.C {
			skill6Ch <- 0
		}
	}()
}

//重置定时器
func reset() {
	ticker1.Reset(time1)
	ticker2.Reset(time2)
	ticker3.Reset(time3)
	ticker5.Reset(time5)
	ticker6.Reset(time6)
}

//暂停取消定时器
func stop() {
	ticker1.Stop()
	ticker2.Stop()
	ticker3.Stop()
	ticker5.Stop()
	ticker6.Stop()
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
	return rand.Intn(60) + 60
}

//随机延迟时间
func randDelayMin() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10) + 10
}
