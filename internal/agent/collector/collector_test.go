package collector

import (
	"testing"

	"github.com/stretchr/testify/mock"
	models "github.com/tatiana4991/metrics/internal/model"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) SetGauge(name string, value float64) {
	m.Called(name, value)
}

func (m *MockStorage) IncCounter(name string, delta int64) {
	m.Called(name, delta)
}

func (m *MockStorage) GetAll() []models.Metrics {
	args := m.Called()
	if result, ok := args.Get(0).([]models.Metrics); ok {
		return result
	}
	return nil
}

func TestCollector_Collect(t *testing.T) {
	mockStore := new(MockStorage)

	mockStore.On("SetGauge", mock.AnythingOfType("string"), mock.AnythingOfType("float64")).Times(27)

	mockStore.On("SetGauge", "RandomValue", mock.AnythingOfType("float64")).Once()

	mockStore.On("IncCounter", "PollCount", int64(1)).Once()

	collector := NewCollector(mockStore)
	collector.Collect()

	mockStore.AssertExpectations(t)
}
