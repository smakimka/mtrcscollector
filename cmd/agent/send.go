package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/smakimka/mtrcscollector/internal/mtrcs"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func sendMetrics(wg *sync.WaitGroup, s storage.Storage, client *http.Client, logger *log.Logger) {
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
			go sendMetric(&sendWg, client, metric, logger)
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

func sendMetric(wg *sync.WaitGroup, client *http.Client, m mtrcs.Metric, logger *log.Logger) {
	defer wg.Done()
	logger.Printf("sending metric %s", m)

	reqURL := fmt.Sprintf("%s/update/%s/%s/%s", serverAddr, m.GetType(), m.GetName(), m.GetStringValue())
	req, err := http.NewRequest("POST", reqURL, nil)
	if err != nil {
		logger.Printf("error creating request (%s)", err.Error())
		return
	}
	req.Header.Add("Content-Type", "text/plain")

	logger.Printf("sending update metric request (%s)", reqURL)
	resp, err := client.Do(req)
	if err != nil {
		logger.Printf("error sending request (%s)", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Printf("got not ok status (%d)", resp.StatusCode)
	}
}
