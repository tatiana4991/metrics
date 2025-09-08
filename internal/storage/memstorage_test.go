package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	models "github.com/tatiana4991/metrics/internal/model"
)

func TestMemStorage_SetGauge(t *testing.T) {
	store := NewMemStorage()
	value := 123.456
	store.SetGauge("TestGauge", value)

	metrics := store.GetAll()
	assert.Len(t, metrics, 1)
	assert.Equal(t, "TestGauge", metrics[0].ID)
	assert.Equal(t, models.Gauge, metrics[0].MType)
	assert.NotNil(t, metrics[0].Value)
	assert.Equal(t, value, *metrics[0].Value)
}

func TestMemStorage_IncCounter(t *testing.T) {
	store := NewMemStorage()
	delta := int64(5)
	store.IncCounter("TestCounter", delta)

	metrics := store.GetAll()
	assert.Len(t, metrics, 1)
	assert.Equal(t, "TestCounter", metrics[0].ID)
	assert.Equal(t, models.Counter, metrics[0].MType)
	assert.NotNil(t, metrics[0].Delta)
	assert.Equal(t, delta, *metrics[0].Delta)

	store.IncCounter("TestCounter", 3)
	metrics = store.GetAll()
	assert.Equal(t, int64(8), *metrics[0].Delta)
}

func TestMemStorage_IncCounter_ExistingWithNilDelta(t *testing.T) {
	store := NewMemStorage()

	store.IncCounter("TestCounter", 0)

	metrics := store.GetAll()
	if len(metrics) > 0 {
		metrics[0].Delta = nil
	}

	store.IncCounter("TestCounter", 5)

	metrics = store.GetAll()
	assert.Len(t, metrics, 1)
	assert.Equal(t, int64(5), *metrics[0].Delta)
}

func TestMemStorage_IncCounter_ExistingCounterWithNilDelta(t *testing.T) {
	store := NewMemStorage()

	store.IncCounter("TestCounter", 0)

	metrics := store.GetAll()
	if len(metrics) > 0 {
		metrics[0].Delta = nil
	}

	store.IncCounter("TestCounter", 5)

	metrics = store.GetAll()
	assert.Len(t, metrics, 1)
	assert.Equal(t, int64(5), *metrics[0].Delta)
}
