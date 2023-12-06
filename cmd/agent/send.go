package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/smakimka/mtrcscollector/internal/mtrcs"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func sendMetrics(wg *sync.WaitGroup, s storage.Storage, client *resty.Client, logger *log.Logger) {
	defer wg.Done()

	for {
		metrics, err := s.GetAllMetrics()
		if err != nil {
			logger.Printf("error getting metrics from storage: %s", err.Error())
		}

		count := 0
		var sendWg sync.WaitGroup
		for _, metric := range metrics {
			sendWg.Add(1)
			go sendMetric(&sendWg, s, client, metric, logger)
			count++
			if count == concurrentMetricsSendCount {
				sendWg.Wait()
				count = 0
			}
		}
		sendWg.Wait()

		time.Sleep(reportInterval)
	}
}

func sendMetric(wg *sync.WaitGroup, s storage.Storage, client *resty.Client, m mtrcs.Metric, logger *log.Logger) {
	defer wg.Done()

	var reqURL string
	switch m.GetName() {
	case "PollCount":
		pollCount, err := s.GetMetric("counter", "PollCount")
		if err != nil {
			logger.Printf("error getting poll count (%s)", err.Error())
			return
		}

		LastpollCount, err := s.GetMetric("gauge", "LastPollCount")
		if err != nil {
			logger.Printf("error getting last poll count (%s)", err.Error())
			return
		}

		count, countOK := pollCount.(mtrcs.CounterMetric)
		lastCount, lastCountOK := LastpollCount.(mtrcs.GaugeMetric)
		if countOK && lastCountOK {
			reqURL = fmt.Sprintf("/update/%s/%s/%s",
				m.GetType(),
				m.GetName(),
				fmt.Sprint(count.Value-int64(lastCount.Value)),
			)
			s.UpdateMetric(mtrcs.GaugeMetric{Name: "LastPollCount", Value: float64(count.Value)})
		}

	default:
		reqURL = fmt.Sprintf("/update/%s/%s/%s", m.GetType(), m.GetName(), m.GetStringValue())
	}

	logger.Printf("sending update metric request (%s)", reqURL)
	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Post(reqURL)

	if err != nil {
		logger.Printf("error sending request %s", err.Error())
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Printf("got not ok status (%d)", resp.StatusCode())
	}
}
