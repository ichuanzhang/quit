package quit

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var defaultQuit = &Quit{}

type HandlerFunc func()

type Quit struct {
	m sync.Map
}

// Handle register handle
func Handle(name string, handler HandlerFunc) {
	defaultQuit.Handle(name, handler)
}

// DeleteHandle delete handle
func DeleteHandle(name string) {
	defaultQuit.DeleteHandle(name)
}

// Wait capture program exit signal
func Wait() {
	defaultQuit.Wait()
}

func New() *Quit {
	return &Quit{}
}

// Handle register handle
func (q *Quit) Handle(name string, handler HandlerFunc) {
	if _, ok := q.m.Load(name); ok {
		return
	}
	q.m.Store(name, handler)
}

// DeleteHandle delete handle
func (q *Quit) DeleteHandle(name string) {
	if _, ok := q.m.Load(name); ok {
		q.m.Delete(name)
	}
}

// Wait capture program exit signal
func (q *Quit) Wait() {
	sig := make(chan os.Signal)

	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	for s := range sig {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
			q.m.Range(func(key, value interface{}) bool {
				handler, ok := value.(HandlerFunc)
				if ok {
					handler()
				}
				return true
			})
			return
		}
	}
}
