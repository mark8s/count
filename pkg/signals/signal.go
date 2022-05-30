package signals

import (
	"os"
	"os/signal"
)

var onlyOneSignalHandler = make(chan struct{})

// SetupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func SetupSignalHandler() (stopCh <-chan struct{}) {
	// 关闭channel
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	// 容量为2的channel
	c := make(chan os.Signal, 2)
	// 指定监听信号，监听Ctrl+C
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
