package shutdown

import (
	zapTool "github.com/justdomepaul/toolbox/zap"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Option interface
type Option interface {
	Apply(*Shutdown)
}

// WithQuit method
func WithQuit(quit chan os.Signal) Option {
	return withQuit{quit: quit}
}

type withQuit struct {
	quit chan os.Signal
}

// Apply method
func (w withQuit) Apply(c *Shutdown) {
	c.quit = w.quit
}

// WithDone method
func WithDone(done chan bool) Option {
	return withDone{done: done}
}

type withDone struct {
	done chan bool
}

// Apply method
func (w withDone) Apply(c *Shutdown) {
	c.done = w.done
}

// WithServerTimeout method
func WithServerTimeout(duration time.Duration) Option {
	return withServerTimeout{server: duration}
}

type withServerTimeout struct {
	server time.Duration
}

// Apply method
func (w withServerTimeout) Apply(c *Shutdown) {
	c.serverTimeout = w.server
}

// WithEndTask method
func WithEndTask(fn func()) Option {
	return withEndTask{fn: fn}
}

type withEndTask struct {
	fn func()
}

// Apply method
func (w withEndTask) Apply(c *Shutdown) {
	c.endTask = w.fn
}

// Shutdown type
type Shutdown struct {
	quit          chan os.Signal
	done          chan bool
	serverTimeout time.Duration
	endTask       func()
}

// Shutdown method
func (s Shutdown) Shutdown() {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	signal.Notify(s.quit, syscall.SIGINT, syscall.SIGTERM)
	<-s.quit
	zapTool.Logger.Info("Start Shutdown server ...")

	if s.endTask != nil {
		zapTool.Logger.Warn("shutdown endTask ...", zap.String("system", "Shutdown"))
		s.endTask()
	}
	zapTool.Logger.Info("system Successfully Stop", zap.String("system", "Shutdown"))
	if s.done != nil {
		s.done <- true
	}
}

// NewShutdown method
func NewShutdown(options ...Option) *Shutdown {
	shutdown := &Shutdown{}
	for _, option := range options {
		option.Apply(shutdown)
	}

	if shutdown.quit == nil {
		shutdown.quit = make(chan os.Signal)
	}
	if shutdown.serverTimeout == 0 {
		shutdown.serverTimeout = 5 * time.Second
	}
	return shutdown
}
