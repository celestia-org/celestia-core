package v2

import (
	"fmt"
	"time"

	"github.com/tendermint/tendermint/libs/log"
)

type timeCheck struct {
	priorityHigh
	time time.Time
}

func schedulerHandle(event Event) (Event, error) {
	switch event.(type) {
	case timeCheck:
		fmt.Println("scheduler handle timeCheck")
	}
	return noOp, nil
}

func processorHandle(event Event) (Event, error) {
	switch event.(type) {
	case timeCheck:
		fmt.Println("processor handle timeCheck")
	}
	return noOp, nil

}

type Reactor struct {
	events    chan Event
	stopDemux chan struct{}
	scheduler *Routine
	processor *Routine
	ticker    *time.Ticker
	logger    log.Logger
}

var bufferSize int = 10

func NewReactor() *Reactor {
	return &Reactor{
		events:    make(chan Event, bufferSize),
		stopDemux: make(chan struct{}),
		scheduler: newRoutine("scheduler", schedulerHandle),
		processor: newRoutine("processor", processorHandle),
		ticker:    time.NewTicker(1 * time.Second),
		logger:    log.NewNopLogger(),
	}
}

// nolint:unused
func (r *Reactor) setLogger(logger log.Logger) {
	r.logger = logger
	r.scheduler.setLogger(logger)
	r.processor.setLogger(logger)
}

func (r *Reactor) Start() {
	go r.scheduler.start()
	go r.processor.start()
	go r.demux()

	<-r.scheduler.ready()
	<-r.processor.ready()

	go func() {
		for t := range r.ticker.C {
			r.events <- timeCheck{time: t}
		}
	}()
}

// Would it be possible here to provide some kind of type safety for the types
// of events that each routine can produce and consume?
func (r *Reactor) demux() {
	for {
		select {
		case event := <-r.events:
			// XXX: check for backpressure
			r.scheduler.trySend(event)
			r.processor.trySend(event)
		case _ = <-r.stopDemux:
			r.logger.Info("demuxing stopped")
			return
		case event := <-r.scheduler.next():
			r.processor.trySend(event)
		case event := <-r.processor.next():
			r.scheduler.trySend(event)
		case err := <-r.scheduler.final():
			r.logger.Info(fmt.Sprintf("scheduler final %s", err))
		case err := <-r.processor.final():
			r.logger.Info(fmt.Sprintf("processor final %s", err))
			// XXX: switch to consensus
		}
	}
}

func (r *Reactor) Stop() {
	r.logger.Info("reactor stopping")

	r.ticker.Stop()
	r.scheduler.stop()
	r.processor.stop()
	close(r.stopDemux)
	close(r.events)

	r.logger.Info("reactor stopped")
}

func (r *Reactor) Receive(event Event) {
	// XXX: decode and serialize write events
	r.events <- event
}

func (r *Reactor) AddPeer() {
	// TODO: add peer event and send to demuxer
}
