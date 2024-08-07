package utils

import (
	"fmt"
	"runtime"
	"syscall"
	"time"
)

// GetAlloc 查看内存使用
/*
 memBefore := getAlloc()
 // do something
 memAfter := getAlloc()
 memUsed := memAfter - memBefore
 fmt.Printf("Memory used: %dKB\n", memUsed/1024)
*/
func GetAlloc() uint64 {
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	return m.Alloc
}

// PrintRuntimeSyscallInfo 打印运行时系统调用信息
/*
	runtime.GC()
	startTime := time.Now()
	// do something
	PrintRuntimeSyscallInfo("descInfo",startTime)
*/
func PrintRuntimeSyscallInfo(descInfo string, startTime time.Time) {

	var memStatus runtime.MemStats // 内存分配器的统计信息
	var rUsage syscall.Rusage      // 通过系统调用获取到的资源值
	var byteToMb = func(b uint64) uint64 {
		return b / 1024 / 1024
	}

	runtime.ReadMemStats(&memStatus)
	err := syscall.Getrusage(syscall.RUSAGE_SELF, &rUsage)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Desc:%s\nCost = %s\nRSS=%v\nAlloc = %v MB\nTotalAlloc = %v MB\nSys = %v MB\nNumGc = %v\n",
		descInfo,
		time.Since(startTime),
		byteToMb(uint64(rUsage.Maxrss)),
		byteToMb(memStatus.Alloc),
		byteToMb(memStatus.TotalAlloc),
		byteToMb(memStatus.Sys),
		memStatus.NumGC,
	)
}

// GetCallInfo 调用的文件所在的行
func GetCallInfo() string {
	_, file, line, _ := runtime.Caller(2)
	return fmt.Sprintf("%s:%d", file, line)
}
