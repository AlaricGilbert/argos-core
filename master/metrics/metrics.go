package metrics

import (
	"container/list"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
)

const UpdateInterval = time.Minute

var (
	ReportMetrics = NewHistoricalMeter("report")
)

type HistoricalMeter struct {
	meter    metrics.Meter
	history  *list.List
	interval time.Duration
	maximum  int
	mu       sync.Mutex
}

func NewHistoricalMeter(name string) *HistoricalMeter {
	return &HistoricalMeter{
		meter:    metrics.GetOrRegisterMeter(name, nil),
		history:  list.New(),
		interval: UpdateInterval,
		maximum:  60,
		mu:       sync.Mutex{},
	}
}

func (m *HistoricalMeter) Mark(n int64) {
	m.meter.Mark(n)
}

func (m *HistoricalMeter) Host() {
	tick := time.Tick(m.interval)
	for range tick {
		m.mu.Lock()
		if m.history.Len() == 60 {
			m.history.Remove(m.history.Front())
		}
		m.history.PushBack(m.meter.Rate1())
		m.mu.Unlock()
	}
}

func (m *HistoricalMeter) GetMetrics() []float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	history := make([]float64, m.history.Len())
	i := 0
	for e := m.history.Front(); e != nil; e = e.Next() {
		history[i] = e.Value.(float64)
		i++
	}
	return history
}
