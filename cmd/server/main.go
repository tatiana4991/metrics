package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/go-chi/chi/v5"

	"github.com/tatiana4991/metrics/internal/config"
	models "github.com/tatiana4991/metrics/internal/model"
	"github.com/tatiana4991/metrics/internal/storage"
)

type TemplateData struct {
	Gauges   []models.Metrics
	Counters []models.Metrics
}

func main() {
	cfg := config.Load()
	store := storage.NewMemStorage()

	tmpl, err := template.ParseFiles("./cmd/server/tmpl/index.html")
	if err != nil {
		log.Fatal("Failed to load template:", err)
	}

	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		update(w, r, store)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		getAll(w, store, tmpl)
	})

	r.Get("/value/{type}/{name}", func(w http.ResponseWriter, r *http.Request) {
		getOne(w, r, store)
	})

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}

func update(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	valueStr := chi.URLParam(r, "value")

	if metricType != models.Gauge && metricType != models.Counter {
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}

	switch metricType {
	case models.Gauge:
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			http.Error(w, "Invalid gauge value", http.StatusBadRequest)
			return
		}
		store.SetGauge(metricName, value)

	case models.Counter:
		delta, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid counter value", http.StatusBadRequest)
			return
		}
		store.IncCounter(metricName, delta)
	}

	log.Printf("Updated metric: type=%s, name=%s, value=%s", metricType, metricName, valueStr)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func getAll(w http.ResponseWriter, store storage.Storage, tmpl *template.Template) {
	metrics := store.GetAll()
	var gauges, counters []models.Metrics

	for _, m := range metrics {
		switch m.MType {
		case models.Gauge:
			gauges = append(gauges, m)
		case models.Counter:
			counters = append(counters, m)
		default:
			log.Printf("Unknown metric type: %s", m.MType)
		}
	}

	data := TemplateData{
		Gauges:   gauges,
		Counters: counters,
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	response := []byte(buf.String())
	w.Header().Set("Content-Length", strconv.Itoa(len(response)))
	w.Write(response)
}

func getOne(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	if metricType != models.Gauge && metricType != models.Counter {
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}

	var valueStr string
	found := false

	for _, m := range store.GetAll() {
		if m.ID == metricName && m.MType == metricType {
			found = true
			switch m.MType {
			case models.Gauge:
				if m.Value != nil {
					valueStr = strconv.FormatFloat(*m.Value, 'f', -1, 64)
				} else {
					valueStr = "0"
				}
			case models.Counter:
				if m.Delta != nil {
					valueStr = strconv.FormatInt(*m.Delta, 10)
				} else {
					valueStr = "0"
				}
			}
			break
		}
	}

	if !found {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	response := []byte(valueStr)
	w.Header().Set("Content-Length", strconv.Itoa(len(response)))
	w.Write(response)
}
