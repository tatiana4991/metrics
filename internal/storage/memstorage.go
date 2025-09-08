package storage

import (
	"sync"

	models "github.com/tatiana4991/metrics/internal/model"
)

type MemStorage struct {
	mu      sync.RWMutex
	metrics map[string]*models.Metrics
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]*models.Metrics),
	}
}

func (s *MemStorage) SetGauge(name string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	metric := &models.Metrics{
		ID:    name,
		MType: models.Gauge,
		Value: &value,
	}
	s.metrics[name] = metric
}

func (s *MemStorage) IncCounter(name string, delta int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if m, exists := s.metrics[name]; exists && m.MType == models.Counter {
		if m.Delta != nil {
			*m.Delta += delta
		} else {
			m.Delta = &delta
		}
	} else {
		s.metrics[name] = &models.Metrics{
			ID:    name,
			MType: models.Counter,
			Delta: &delta,
		}
	}
}

func (s *MemStorage) GetAll() []models.Metrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.Metrics
	for _, m := range s.metrics {
		result = append(result, *m)
	}
	return result
}
