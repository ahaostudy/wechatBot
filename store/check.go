package store

import (
	"sync"
	"time"
)

var preTime int64

// 单位转换 1s == 1000000us
// 检查周期
const cycle int64 = 600 * 1000000

// 有效时长
const duration int64 = 600 * 1000000

func Check(wg *sync.WaitGroup) {
	defer wg.Done()
	now := time.Now().UnixMicro()
	// 未到检测周期
	if now-preTime < cycle {
		return
	}
	for k, v := range CS {
		if now-v > duration {
			println("清楚记录:", k, v)
			delete(MS, k)
		}
	}
}
