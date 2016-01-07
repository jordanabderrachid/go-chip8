package timer

import (
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	var b byte = 60 // Should take 1s to decrease
	var timer Timer
	go timer.Run(&b)

	time.Sleep(1*time.Second + 100*time.Millisecond)
	if b != 0 {
		t.Errorf("byte should be 0, got %d", b)
	}
}
