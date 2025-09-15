package sender

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/tatiana4991/metrics/internal/config"
	models "github.com/tatiana4991/metrics/internal/model"
	"github.com/tatiana4991/metrics/internal/storage"
)

type Sender struct {
	config *config.Config
	store  storage.Storage
	client *http.Client
}

func NewSender(cfg *config.Config, store storage.Storage) *Sender {
	return &Sender{
		config: cfg,
		store:  store,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *Sender) SendAll() {
	metricsList := s.store.GetAll()

	for _, m := range metricsList {
		var valueStr string

		switch m.MType {
		case models.Gauge:
			if m.Value != nil {
				valueStr = strconv.FormatFloat(*m.Value, 'f', -1, 64)
			} else {
				log.Printf("Gauge %s has nil value", m.ID)
				continue
			}
		case models.Counter:
			if m.Delta != nil {
				valueStr = strconv.FormatInt(*m.Delta, 10)
			} else {
				log.Printf("Counter %s has nil delta", m.ID)
				continue
			}
		default:
			log.Printf("Unknown metric type: %s", m.MType)
			continue
		}

		url := fmt.Sprintf("%s/update/%s/%s/%s",
			s.config.ServerAddress,
			m.MType,
			m.ID,
			valueStr)

		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			log.Printf("Failed to create request for %s: %v", m.ID, err)
			continue
		}
		req.Header.Set("Content-Type", "text/plain")

		resp, err := s.client.Do(req)
		if err != nil {
			log.Printf("Failed to send metric %s: %v", m.ID, err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Server responded with %d for metric %s", resp.StatusCode, m.ID)
		} else {
			log.Printf("Sent %s %s = %s", m.MType, m.ID, valueStr)
		}
	}
}
