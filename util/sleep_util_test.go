package util

import (
	"testing"
)

func TestIsSleepMode(t *testing.T) {
	start := "22:00"
	end := "08:00"
	sleepMode := IsSleepMode(start, end)
	t.Log("睡眠模式:", sleepMode)
}
