package hpecxi

import (
	"testing"
)

func hasHPECXI(t *testing.T) bool {
	devices := GetDevices()
	if len(devices) <= 0 {
		return false
	}
	return true
}
