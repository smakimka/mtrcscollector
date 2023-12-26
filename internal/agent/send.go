package agent

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/smakimka/mtrcscollector/internal/agent/config"
	"github.com/smakimka/mtrcscollector/internal/logger"
	"github.com/smakimka/mtrcscollector/internal/model"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func SendMetrics(cfg *config.Config, s storage.Storage, client *resty.Client, c chan error) {
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
		go sendGaugeMetric(s, client, gaugeMetric, c)
	}

	for _, counterMetric := range counterMetrics {
		if counterMetric.Name == "PollCount" {
			go sendPollCount(s, client, counterMetric, c)
			continue
		}
		go sendCounterMetric(s, client, counterMetric, c)
	}
}

func sendPollCount(s storage.Storage, client *resty.Client, m model.CounterMetric, c chan error) {
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
	inc := pollCount.Value - int64(lastPollCount.Value)
	reqData := &model.MetricsData{Name: m.Name, Kind: m.GetType(), Delta: &inc}

	s.UpdateGaugeMetric(model.GaugeMetric{Name: "LastPollCount", Value: float64(pollCount.Value)})
	logger.Log.Debug().Msg(fmt.Sprintf("sending update poll count request (%s)", fmt.Sprint(reqData)))

	sendRequest(reqData, client, c)
}

func sendGaugeMetric(s storage.Storage, client *resty.Client, m model.GaugeMetric, c chan error) {
	reqData := &model.MetricsData{Name: m.Name, Kind: m.GetType(), Value: &m.Value}
	logger.Log.Debug().Msg(fmt.Sprintf("sending update gauge metric request (%s)", fmt.Sprint(reqData)))
	sendRequest(reqData, client, c)
}

func sendCounterMetric(s storage.Storage, client *resty.Client, m model.CounterMetric, c chan error) {
	reqData := &model.MetricsData{Name: m.Name, Kind: m.GetType(), Delta: &m.Value}
	logger.Log.Debug().Msg(fmt.Sprintf("sending update counter metric request (%s)", fmt.Sprint(reqData)))
	sendRequest(reqData, client, c)
}

func sendRequest(data *model.MetricsData, client *resty.Client, c chan error) {
	resp, err := client.R().
		SetBody(data).
		Post("/update")

	if err != nil {
		c <- err
		return
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Log.Warn().Msg(fmt.Sprintf("got not ok status (%d)", resp.StatusCode()))
	}
}
