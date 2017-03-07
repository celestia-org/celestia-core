package eventmeter

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	metrics "github.com/rcrowley/go-metrics"
	events "github.com/tendermint/go-events"
	client "github.com/tendermint/go-rpc/client"
	log15 "github.com/tendermint/log15"
)

// Log allows you to set your own logger.
var Log log15.Logger

//------------------------------------------------------
// Generic system to subscribe to events and record their frequency
//------------------------------------------------------

//------------------------------------------------------
// Meter for a particular event

// Closure to enable side effects from receiving an event
type EventCallbackFunc func(em *EventMetric, data events.EventData)

// Metrics for a given event
type EventMetric struct {
	ID          string    `json:"id"`
	Started     time.Time `json:"start_time"`
	LastHeard   time.Time `json:"last_heard"`
	MinDuration int64     `json:"min_duration"`
	MaxDuration int64     `json:"max_duration"`

	// tracks event count and rate
	meter metrics.Meter

	// filled in from the Meter
	Count    int64   `json:"count"`
	Rate1    float64 `json:"rate_1" wire:"unsafe"`
	Rate5    float64 `json:"rate_5" wire:"unsafe"`
	Rate15   float64 `json:"rate_15" wire:"unsafe"`
	RateMean float64 `json:"rate_mean" wire:"unsafe"`

	// so the event can have effects in the event-meter's consumer.
	// runs in a go routine
	callback EventCallbackFunc
}

func (metric *EventMetric) Copy() *EventMetric {
	metric2 := *metric
	metric2.meter = metric.meter.Snapshot()
	return &metric2
}

// called on GetMetric
func (metric *EventMetric) fillMetric() *EventMetric {
	metric.Count = metric.meter.Count()
	metric.Rate1 = metric.meter.Rate1()
	metric.Rate5 = metric.meter.Rate5()
	metric.Rate15 = metric.meter.Rate15()
	metric.RateMean = metric.meter.RateMean()
	return metric
}

//------------------------------------------------------
// Websocket client and event meter for many events

const maxPingsPerPong = 30 // if we haven't received a pong in this many attempted pings we kill the conn

// Get the eventID and data out of the raw json received over the go-rpc websocket
type EventUnmarshalFunc func(b json.RawMessage) (string, events.EventData, error)

// Closure to enable side effects from receiving a pong
type LatencyCallbackFunc func(meanLatencyNanoSeconds float64)

// Closure to notify consumer that the connection died
type DisconnectCallbackFunc func()

// Each node gets an event meter to track events for that node
type EventMeter struct {
	wsc *client.WSClient

	mtx    sync.Mutex
	events map[string]*EventMetric

	// to record ws latency
	timer              metrics.Timer
	lastPing           time.Time
	receivedPong       bool
	latencyCallback    LatencyCallbackFunc
	disconnectCallback DisconnectCallbackFunc

	unmarshalEvent EventUnmarshalFunc

	quit chan struct{}
}

func NewEventMeter(addr string, unmarshalEvent EventUnmarshalFunc) *EventMeter {
	em := &EventMeter{
		wsc:            client.NewWSClient(addr, "/websocket"),
		events:         make(map[string]*EventMetric),
		timer:          metrics.NewTimer(),
		receivedPong:   true,
		unmarshalEvent: unmarshalEvent,
		quit:           make(chan struct{}),
	}
	return em
}

func (em *EventMeter) String() string {
	return em.wsc.Address
}

func (em *EventMeter) Start() error {
	if _, err := em.wsc.Start(); err != nil {
		return err
	}

	em.wsc.Conn.SetPongHandler(func(m string) error {
		// NOTE: https://github.com/gorilla/websocket/issues/97
		em.mtx.Lock()
		defer em.mtx.Unlock()
		em.receivedPong = true
		em.timer.UpdateSince(em.lastPing)
		if em.latencyCallback != nil {
			go em.latencyCallback(em.timer.Mean())
		}
		return nil
	})
	go em.receiveRoutine()
	return nil
}

func (em *EventMeter) Stop() {
	<-em.quit

	em.RegisterDisconnectCallback(nil) // so we don't try and reconnect
	em.wsc.Stop()                      // close(wsc.Quit)
}

func (em *EventMeter) StopAndReconnect() {
	em.wsc.Stop()

	em.mtx.Lock()
	defer em.mtx.Unlock()
	if em.disconnectCallback != nil {
		go em.disconnectCallback()
	}
}

func (em *EventMeter) Subscribe(eventID string, cb EventCallbackFunc) error {
	em.mtx.Lock()
	defer em.mtx.Unlock()

	if _, ok := em.events[eventID]; ok {
		return fmt.Errorf("subscribtion already exists")
	}
	if err := em.wsc.Subscribe(eventID); err != nil {
		return err
	}

	metric := &EventMetric{
		ID:          eventID,
		Started:     time.Now(),
		MinDuration: 1 << 62,
		meter:       metrics.NewMeter(),
		callback:    cb,
	}
	em.events[eventID] = metric
	return nil
}

func (em *EventMeter) Unsubscribe(eventID string) error {
	em.mtx.Lock()
	defer em.mtx.Unlock()
	if err := em.wsc.Unsubscribe(eventID); err != nil {
		return err
	}
	// XXX: should we persist or save this info first?
	delete(em.events, eventID)
	return nil
}

// Fill in the latest data for an event and return a copy
func (em *EventMeter) GetMetric(eventID string) (*EventMetric, error) {
	em.mtx.Lock()
	defer em.mtx.Unlock()
	metric, ok := em.events[eventID]
	if !ok {
		return nil, fmt.Errorf("Unknown event %s", eventID)
	}
	return metric.fillMetric().Copy(), nil
}

// Return the average latency over the websocket
func (em *EventMeter) Latency() float64 {
	em.mtx.Lock()
	defer em.mtx.Unlock()
	return em.timer.Mean()
}

func (em *EventMeter) RegisterLatencyCallback(f LatencyCallbackFunc) {
	em.mtx.Lock()
	defer em.mtx.Unlock()
	em.latencyCallback = f
}

func (em *EventMeter) RegisterDisconnectCallback(f DisconnectCallbackFunc) {
	em.mtx.Lock()
	defer em.mtx.Unlock()
	em.disconnectCallback = f
}

//------------------------------------------------------

func (em *EventMeter) receiveRoutine() {
	pingTime := time.Second * 1
	pingTicker := time.NewTicker(pingTime)
	pingAttempts := 0 // if this hits maxPingsPerPong we kill the conn
	var err error
	for {
		select {
		case <-pingTicker.C:
			if pingAttempts, err = em.pingForLatency(pingAttempts); err != nil {
				Log.Error("Failed to write ping message on websocket", err)
				em.StopAndReconnect()
				return
			} else if pingAttempts >= maxPingsPerPong {
				Log.Error(fmt.Sprintf("Have not received a pong in %v", time.Duration(pingAttempts)*pingTime))
				em.StopAndReconnect()
				return
			}
		case r := <-em.wsc.ResultsCh:
			if r == nil {
				em.StopAndReconnect()
				return
			}
			eventID, data, err := em.unmarshalEvent(r)
			if err != nil {
				Log.Error(err.Error())
				continue
			}
			if eventID != "" {
				em.updateMetric(eventID, data)
			}
		case <-em.wsc.Quit:
			em.StopAndReconnect()
			return
		case <-em.quit:
			return
		}
	}
}

func (em *EventMeter) pingForLatency(pingAttempts int) (int, error) {
	em.mtx.Lock()
	defer em.mtx.Unlock()

	// ping to record latency
	if !em.receivedPong {
		return pingAttempts + 1, nil
	}

	em.lastPing = time.Now()
	em.receivedPong = false
	err := em.wsc.Conn.WriteMessage(websocket.PingMessage, []byte{})
	if err != nil {
		return pingAttempts, err
	}
	return 0, nil
}

func (em *EventMeter) updateMetric(eventID string, data events.EventData) {
	em.mtx.Lock()
	defer em.mtx.Unlock()

	metric, ok := em.events[eventID]
	if !ok {
		// we already unsubscribed, or got an unexpected event
		return
	}

	last := metric.LastHeard
	metric.LastHeard = time.Now()
	metric.meter.Mark(1)
	dur := int64(metric.LastHeard.Sub(last))
	if dur < metric.MinDuration {
		metric.MinDuration = dur
	}
	if !last.IsZero() && dur > metric.MaxDuration {
		metric.MaxDuration = dur
	}

	if metric.callback != nil {
		go metric.callback(metric.Copy(), data)
	}
}
