package fias

import (
	"testing"
	"time"
)

func TestGetLastLocalArchivePath(t *testing.T) {
	y := time.Now().AddDate(0, 0, -1)

	path, isLast := GetLastLocalArchivePath(y)

	if path == "" {
		t.Error("")
	}
	if isLast == false {
		t.Error("")
	}
}
