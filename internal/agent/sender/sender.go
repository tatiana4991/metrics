package sender

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tatiana4991/metrics/internal/config"
	models "github.com/tatiana4991/metrics/internal/model"
	"github.com/tatiana4991/metrics/internal/storage"
)

type Sender struct {
	config *config.Config
	store  storage.Storage
	client *resty.Client
}

func NewSender(cfg *config.Config, store storage.Storage) *Sender {
	return &Sender{
		config: cfg,
		store:  store,
		client: resty.New().SetTimeout(5 * time.Second),
	}
}

func (s *Sender) SendAll() {
	metricsList := s.store.GetAll()

	for _, m := range metricsList {
		var valueStr string

		switch m.MType {
		case models.Gauge:
			if m.Value != nil {
				valueStr = fmt.Sprintf("%.f", *m.Value)
			} else {
				log.Printf("Gauge %s has nil value", m.ID)
				continue
			}
		case models.Counter:
			if m.Delta != nil {
				valueStr = fmt.Sprintf("%d", *m.Delta)
			} else {
				log.Printf("Counter %s has nil delta", m.ID)
				continue
			}
		default:
			log.Printf("Unknown metric type: %s", m.MType)
			continue
		}

		serverAddr := s.config.ServerAddress
		if !strings.HasPrefix(serverAddr, "http://") && !strings.HasPrefix(serverAddr, "https://") {
			serverAddr = "http://" + serverAddr
		}

		url := fmt.Sprintf("%s/update/%s/%s/%s",
			serverAddr,
			m.MType,
			m.ID,
			valueStr)

		resp, err := s.client.R().
			SetHeader("Content-Type", "text/plain").
			Post(url)

		if err != nil {
			log.Printf("Failed to send metric %s: %v", m.ID, err)
			continue
		}

		if resp.StatusCode() == 200 {
			log.Printf("Sent %s %s = %s", m.MType, m.ID, valueStr)
		} else {
			log.Printf("Server responded with %d for metric %s", resp.StatusCode(), m.ID)
		}
	}
}
