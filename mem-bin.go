package main

import (
	"fmt"
	"time"
	"syscall"
	"github.com/shirou/gopsutil/mem"
)

func getunixtime() int64 {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic(err)
	}
	now := time.Now()
	t := now.In(loc)
	a := t.UnixNano() / 1000000
	return a
}

// 다른 메모리 정보와는 다르게, 캐시 메모리 정보는 syscall 라이브러리를 통해서 가져올 수 없다.
// 따라서 gopsutil/mem 라이브러리를 통해 데이터를 가져온다.
// 두 라이브러리 모두 /proc/meminfo 파일을 읽어와 메모리 정보를 반환한다.
// 따라서 프로그램이 실행되며 파일 I/O를 수행하므로, 파일 I/O 캐시 데이터가 반영된 결과가 반영된다.
func getCachedMemory() (uint64, error) {
	m, err := mem.VirtualMemory()
	if err != nil {
		return uint64(0), err
	}
	cached := m.Cached

	return cached, nil
}

func memoryInfo() ([]string, error) {
	subs := []string{}

	var memInfo syscall.Sysinfo_t
	if err := syscall.Sysinfo(&memInfo); err != nil {
		fmt.Print("Failed to syscall info: %v", err)
	}

	// total memory
	total := memInfo.Totalram * uint64(memInfo.Unit)

	// free memory
	free := memInfo.Freeram * uint64(memInfo.Unit)

	// buffered memory
	buffered := memInfo.Bufferram * uint64(memInfo.Unit)

	// cached memory
	cached, err := getCachedMemory()
	if err != nil {
		fmt.Print("Failed to get cached memory: %v", err)
	}

	nowtime := getunixtime()
	fmt.Printf("system_mem_free{} %f %d\n", float64(free) / float64(total) * 100, nowtime)
	fmt.Printf("system_mem_free_bytes{} %d %d\n", free, nowtime)
	fmt.Printf("system_mem_buffered{} %f %d\n", float64(buffered) / float64(total) * 100, nowtime)
	fmt.Printf("system_mem_buffered_bytes{} %d %d\n", buffered, nowtime)
	fmt.Printf("system_mem_cached{} %f %d\n", float64(cached) / float64(total) * 100, nowtime)
	fmt.Printf("system_mem_cached_bytes{} %d %d\n", cached, nowtime)

	return subs, nil
}

func main() {
	memoryInfo()
}
