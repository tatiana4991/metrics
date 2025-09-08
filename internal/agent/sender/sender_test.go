package sender

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tatiana4991/metrics/internal/config"
	models "github.com/tatiana4991/metrics/internal/model"
)

type MockStorageForSender struct {
	metrics []models.Metrics
}

func (m *MockStorageForSender) SetGauge(name string, value float64) {}
func (m *MockStorageForSender) IncCounter(name string, delta int64) {}
func (m *MockStorageForSender) GetAll() []models.Metrics            { return m.metrics }

func TestSender_SendAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "text/plain", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	gaugeVal := 123.456
	counterVal := int64(42)
	mockStore := &MockStorageForSender{
		metrics: []models.Metrics{
			{ID: "TestGauge", MType: models.Gauge, Value: &gaugeVal},
			{ID: "TestCounter", MType: models.Counter, Delta: &counterVal},
		},
	}

	cfg := &config.Config{
		ServerAddress: server.URL,
	}

	sender := NewSender(cfg, mockStore)
	sender.SendAll()
}

func TestSender_SendAll_NilValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Should not be called")
	}))
	defer server.Close()

	mockStore := &MockStorageForSender{
		metrics: []models.Metrics{
			{ID: "NilGauge", MType: models.Gauge},
			{ID: "NilCounter", MType: models.Counter},
		},
	}

	cfg := &config.Config{ServerAddress: server.URL}
	sender := NewSender(cfg, mockStore)

	sender.SendAll()
}

func TestSender_SendAll_UnknownMetricType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Should not be called")
	}))
	defer server.Close()

	mockStore := &MockStorageForSender{
		metrics: []models.Metrics{
			{ID: "Unknown", MType: "unknown-type"},
		},
	}

	cfg := &config.Config{ServerAddress: server.URL}
	sender := NewSender(cfg, mockStore)

	sender.SendAll()
}

func TestSender_SendAll_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	val := 123.456
	mockStore := &MockStorageForSender{
		metrics: []models.Metrics{
			{ID: "TestGauge", MType: models.Gauge, Value: &val},
		},
	}

	cfg := &config.Config{ServerAddress: server.URL}
	sender := NewSender(cfg, mockStore)

	sender.SendAll()
}

func TestSender_SendAll_InvalidURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Should not be called")
	}))
	defer server.Close()

	val := 123.456
	mockStore := &MockStorageForSender{
		metrics: []models.Metrics{
			{ID: "Invalid\nName", MType: models.Gauge, Value: &val},
		},
	}

	cfg := &config.Config{ServerAddress: server.URL}
	sender := NewSender(cfg, mockStore)

	sender.SendAll()
}

func TestSender_SendAll_HTTPClientError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second)
	}))
	defer server.Close()

	val := 123.456
	mockStore := &MockStorageForSender{
		metrics: []models.Metrics{
			{ID: "TestGauge", MType: models.Gauge, Value: &val},
		},
	}

	client := &http.Client{Timeout: 100 * time.Millisecond}
	sender := &Sender{
		config: &config.Config{ServerAddress: server.URL},
		store:  mockStore,
		client: client,
	}

	sender.SendAll()
}
