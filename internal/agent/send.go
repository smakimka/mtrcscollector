package agent

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/smakimka/mtrcscollector/internal/agent/config"
	"github.com/smakimka/mtrcscollector/internal/model"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func SendMetrics(cfg *config.Config, s storage.Storage, logger *log.Logger, client *resty.Client, c chan error) {
	gaugeMetrics, err := s.GetAllGaugeMetrics()
	if err != nil {
		c <- err
		return
	}

	counterMetrics, err := s.GetAllCounterMetrics()
	if err != nil {
		c <- err
		return
	}

	for _, gaugeMetric := range gaugeMetrics {
		go sendGaugeMetric(s, client, gaugeMetric, logger, c)
	}

	for _, counterMetric := range counterMetrics {
		if counterMetric.Name == "PollCount" {
			go sendPollCount(s, client, counterMetric, logger, c)
			continue
		}
		go sendCounterMetric(s, client, counterMetric, logger, c)
	}
}

func sendPollCount(s storage.Storage, client *resty.Client, m model.CounterMetric, logger *log.Logger, c chan error) {
	pollCount, err := s.GetCounterMetric("PollCount")
	if err != nil {
		c <- err
		return
	}

	lastPollCount, err := s.GetGaugeMetric("LastPollCount")
	if err != nil {
		c <- err
		return
	}

	reqURL := fmt.Sprintf("/update/%s/%s/%s",
		m.GetType(),
		m.GetName(),
		fmt.Sprint(pollCount.Value-int64(lastPollCount.Value)),
	)
	s.UpdateGaugeMetric(model.GaugeMetric{Name: "LastPollCount", Value: float64(pollCount.Value)})

	logger.Printf("sending update poll count request (%s)", reqURL)
	sendRequest(reqURL, client, logger, c)
}

func sendGaugeMetric(s storage.Storage, client *resty.Client, m model.GaugeMetric, logger *log.Logger, c chan error) {
	reqURL := fmt.Sprintf("/update/%s/%s/%s", m.GetType(), m.GetName(), m.GetStringValue())

	logger.Printf("sending update gauge metric request (%s)", reqURL)
	sendRequest(reqURL, client, logger, c)
}

func sendCounterMetric(s storage.Storage, client *resty.Client, m model.CounterMetric, logger *log.Logger, c chan error) {
	reqURL := fmt.Sprintf("/update/%s/%s/%s", m.GetType(), m.GetName(), m.GetStringValue())

	logger.Printf("sending update counter metric request (%s)", reqURL)
	sendRequest(reqURL, client, logger, c)
}

func sendRequest(reqURL string, client *resty.Client, logger *log.Logger, c chan error) {
	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Post(reqURL)

	if err != nil {
		c <- err
		return
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Printf("got not ok status (%d)", resp.StatusCode())
	}
}
