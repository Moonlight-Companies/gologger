package monitor

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestMonitorStartStop(t *testing.T) {
	// Ensure monitor is stopped before test
	Stop()

	// Start the monitor
	Start(100 * time.Millisecond)

	// Wait for a short period
	time.Sleep(150 * time.Millisecond)

	// Check if monitor is running (indirectly by checking if instance is not nil)
	if instance == nil {
		t.Error("Monitor should be running")
	}

	// Stop the monitor
	Stop()

	// Check if monitor has stopped (instance should be nil)
	if instance != nil {
		t.Error("Monitor should have stopped")
	}
}

func TestMonitorOnlyOneInstance(t *testing.T) {
	// Ensure monitor is stopped before test
	Stop()

	// Start the monitor
	Start(100 * time.Millisecond)

	// Try to start it again
	Start(200 * time.Millisecond)

	// Check if the interval is still the original one
	if instance.interval != 100*time.Millisecond {
		t.Error("Monitor interval should not have changed")
	}

	// Stop the monitor
	Stop()

	// Try to stop it again (this should not panic)
	defer func() {
		if r := recover(); r != nil {
			t.Error("Stopping an already stopped monitor should not panic")
		}
	}()
	Stop()
}

func TestMonitorInterval(t *testing.T) {
	// Ensure monitor is stopped before test
	Stop()

	// Start the monitor with a short interval
	Start(50 * time.Millisecond)

	// Wait for a short period
	time.Sleep(175 * time.Millisecond)

	// Stop the monitor
	Stop()

	// We can't directly check the number of logs, but we can check if the monitor ran without errors
	if instance != nil {
		t.Error("Monitor should have stopped")
	}
}

func TestConcurrentAccess(t *testing.T) {
	// Ensure monitor is stopped before test
	Stop()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Start(100 * time.Millisecond)
			time.Sleep(10 * time.Millisecond)
			Stop()
		}()
	}

	wg.Wait()

	// After all goroutines complete, the monitor should be stopped
	if instance != nil {
		t.Error("Monitor should be stopped after concurrent access")
	}
}

func TestMonitorImpact(t *testing.T) {
	// Ensure monitor is stopped before test
	Stop()

	// Record the number of goroutines before starting the monitor
	goroutinesBefore := runtime.NumGoroutine()

	// Start the monitor
	Start(100 * time.Millisecond)

	// Wait a short while for the monitor to start
	time.Sleep(150 * time.Millisecond)

	// Record the number of goroutines after starting the monitor
	goroutinesAfter := runtime.NumGoroutine()

	// Stop the monitor
	Stop()

	// Check if only one additional goroutine was created
	if goroutinesAfter != goroutinesBefore+1 {
		t.Errorf("Expected one additional goroutine, got %d (before: %d, after: %d)",
			goroutinesAfter-goroutinesBefore, goroutinesBefore, goroutinesAfter)
	}

	// Wait a bit to ensure the goroutine has stopped
	time.Sleep(150 * time.Millisecond)

	// Check if the number of goroutines has returned to the original count
	goroutinesAfterStop := runtime.NumGoroutine()
	if goroutinesAfterStop != goroutinesBefore {
		t.Errorf("Number of goroutines should return to original count. Expected: %d, Got: %d",
			goroutinesBefore, goroutinesAfterStop)
	}
}
