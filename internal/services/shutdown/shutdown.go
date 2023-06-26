package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

var c chan os.Signal

func OnShutdown(fn func()) {
	c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fn()
	}()
}
