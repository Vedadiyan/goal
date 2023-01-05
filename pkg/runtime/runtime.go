package runtime

import (
	"os"
	"os/signal"
)

func WaitForInterrupt(callback func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	callback()
}
