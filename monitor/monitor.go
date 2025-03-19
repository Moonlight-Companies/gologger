package monitor

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/Moonlight-Companies/gologger/logger"
)

var (
	instance *Monitor
	once     sync.Once
	stopChan chan struct{}
	mu       sync.Mutex
)

type Monitor struct {
	interval time.Duration
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func (m *Monitor) run() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			logger.Log.Info("stats Goroutines: %d, Allocated Memory: %s, Total Alloc: %s, Sys Memory: %s, Num GC: %d",
				runtime.NumGoroutine(),
				formatBytes(memStats.Alloc),
				formatBytes(memStats.TotalAlloc),
				formatBytes(memStats.Sys),
				memStats.NumGC,
			)
		case <-stopChan:
			return
		}
	}
}

func Start(interval time.Duration) {
	mu.Lock()
	defer mu.Unlock()

	if instance != nil {
		logger.Log.Warn("Monitor is already running")
		return
	}

	once.Do(func() {
		instance = &Monitor{interval: interval}
		stopChan = make(chan struct{})
		go instance.run()
	})
}

func Stop() {
	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		logger.Log.Warn("Monitor is not running")
		return
	}

	close(stopChan)
	instance = nil
	once = sync.Once{}
}
