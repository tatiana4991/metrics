package storage

import models "github.com/tatiana4991/metrics/internal/model"

type Storage interface {
	SetGauge(name string, value float64)
	IncCounter(name string, delta int64)
	GetAll() []models.Metrics
}
