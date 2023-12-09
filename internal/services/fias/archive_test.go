package fias

import (
	"testing"
)

func TestGetLastLocalArchivePath(t *testing.T) {
	path := GetLastLocalArchivePath()

	if path == "" {
		t.Error("")
	}
}
