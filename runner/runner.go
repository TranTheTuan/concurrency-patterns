package runner

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

var (
	InterruptError = errors.New("interrupt error")
	TimeoutError   = errors.New("timeout error")
)

type Runner struct {
	interruptChan chan os.Signal
	completeChan  chan error
	timeoutChan   <-chan time.Time
	tasks         []func(int)
}

func NewRunner(d time.Duration) *Runner {
	return &Runner{
		interruptChan: make(chan os.Signal, 1),
		completeChan:  make(chan error),
		timeoutChan:   time.After(d),
	}
}

func (r *Runner) AddTasks(tasks ...func(int)) {
	for _, task := range tasks {
		r.tasks = append(r.tasks, task)
	}
}

func (r *Runner) run() error {
	for i, task := range r.tasks {
		if r.gotInterrupt() {
			return InterruptError
		}
		task(i)
	}
	return nil
}

func (r *Runner) gotInterrupt() bool {
	select {
	case <-r.interruptChan:
		return true
	default:
		return false
	}
}

func (r *Runner) Start() error {
	signal.Notify(r.interruptChan, os.Interrupt)

	go func() {
		r.completeChan <- r.run()
	}()

	select {
	case err := <-r.completeChan:
		return err
	case <-r.timeoutChan:
		return TimeoutError
	}
}
