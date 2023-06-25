package exit

import (
	"os"
	"os/signal"
	"syscall"
)

var c chan os.Signal

func OnExit(fn func()) {
	c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fn()
	}()
}
