package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"google.golang.org/grpc/metadata"

	"github.com/smakimka/mtrcscollector/internal/agent/config"
	"github.com/smakimka/mtrcscollector/internal/auth"
	"github.com/smakimka/mtrcscollector/internal/logger"
	"github.com/smakimka/mtrcscollector/internal/model"
	"github.com/smakimka/mtrcscollector/internal/storage"
	pb "github.com/smakimka/mtrcscollector/protobuf/server"
)

func Worker(ctx context.Context, cfg config.Config, client *resty.Client, id int, jobs <-chan model.MetricsData, errs chan<- error) {
	for metricsData := range jobs {
		logger.Log.Debug().Msg(fmt.Sprintf("worker %d started work", id))

		err := sendRequest(ctx, &cfg, metricsData, client)
		if err != nil {
			errs <- err
		}
		logger.Log.Debug().Msg(fmt.Sprintf("worker %d finished work", id))
	}
}

func GRPCWorker(ctx context.Context, cfg config.Config, client pb.MetricsCollectorClient, id int, jobs <-chan model.MetricsData, errs chan<- error) {
	for metricsData := range jobs {
		logger.Log.Debug().Msg(fmt.Sprintf("worker %d started work", id))

		err := sendGRPCRequest(ctx, &cfg, metricsData, client)
		if err != nil {
			errs <- err
		}
		logger.Log.Debug().Msg(fmt.Sprintf("worker %d finished work", id))
	}
}

func SendMetrics(ctx context.Context, _ *config.Config, s storage.Storage, jobs chan<- model.MetricsData, errs chan<- error) {
	gaugeMetrics, err := s.GetAllGaugeMetrics(ctx)
	if err != nil {
		errs <- err
		return
	}

	counterMetrics, err := s.GetAllCounterMetrics(ctx)
	if err != nil {
		errs <- err
		return
	}

	metricsData := model.MetricsData{}

	for i := range gaugeMetrics {
		if gaugeMetrics[i].Name == "LastPollCount" {
			continue
		}

		metricsData = append(metricsData, model.MetricData{
			Name:  gaugeMetrics[i].Name,
			Kind:  model.Gauge,
			Value: &gaugeMetrics[i].Value,
		})
	}

	for i := range counterMetrics {
		if counterMetrics[i].Name == "PollCount" {
			pollCountData, err := getPollCountData(ctx, s, counterMetrics[i])
			if err != nil {
				errs <- err
				return
			}
			metricsData = append(metricsData, pollCountData)
			continue
		}
		metricsData = append(metricsData, model.MetricData{
			Name:  counterMetrics[i].Name,
			Kind:  model.Counter,
			Delta: &counterMetrics[i].Value,
		})
	}

	logger.Log.Debug().Msg("sending job to workers")
	jobs <- metricsData
}

func getPollCountData(ctx context.Context, s storage.Storage, m model.CounterMetric) (model.MetricData, error) {
	data := model.MetricData{}

	pollCount, err := s.GetCounterMetric(ctx, "PollCount")
	if err != nil {
		return data, err
	}

	lastPollCount, err := s.GetGaugeMetric(ctx, "LastPollCount")
	if err != nil {
		return data, err
	}
	inc := pollCount.Value - int64(lastPollCount.Value)

	if err = s.UpdateGaugeMetric(ctx, model.GaugeMetric{Name: "LastPollCount", Value: float64(pollCount.Value)}); err != nil {
		return data, err
	}

	return model.MetricData{
		Name:  m.Name,
		Kind:  model.Counter,
		Delta: &inc,
	}, nil
}

func sendRequest(_ context.Context, cfg *config.Config, data model.MetricsData, client *resty.Client) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if cfg.CryptoKey != nil {
		encBody, err := rsa.EncryptPKCS1v15(rand.Reader, cfg.CryptoKey, body)
		if err != nil {
			return err
		}
		body = encBody
	}

	zipBody := bytes.NewBuffer([]byte{})
	zw := gzip.NewWriter(zipBody)
	zw.Write(body)
	zw.Close()

	req := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("X-Real-IP", cfg.MyIP)

	if cfg.CryptoKey != nil {
		req.SetHeader("Encryption", "crypto-key")
	}

	if auth.Enabled() {
		sign := auth.Sign(zipBody.Bytes())
		req.SetHeader("HashSHA256", hex.EncodeToString(sign))
	}

	// Можно просто SetBody со структурой, которая сюда передается, но надо чтобы в импортах был хоть где-то json, будет тут
	resp, err := req.
		SetBody(zipBody).
		Post("/updates/")

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Log.Warn().Msg(fmt.Sprintf("got not ok status (%d)", resp.StatusCode()))
	}

	return nil
}

func sendGRPCRequest(ctx context.Context, cfg *config.Config, data model.MetricsData, client pb.MetricsCollectorClient) error {
	in := &pb.UpdateMetrics{Metrics: []*pb.Metric{}}
	for _, metric := range data {
		in.Metrics = append(in.Metrics, &pb.Metric{
			Delta: *metric.Delta,
			Value: *metric.Value,
			Name:  metric.Name,
			Kind:  metric.Kind,
		})
	}

	md := metadata.New(map[string]string{"X-Real-IP": cfg.MyIP})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.Update(ctx, in)
	if err != nil {
		return err
	}

	if !resp.Ok {
		logger.Log.Warn().Msg(fmt.Sprintf("got error (%s)", resp.Detail))
	}

	return nil

}
