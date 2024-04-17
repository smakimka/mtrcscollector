package storage

import (
	"context"

	"github.com/smakimka/mtrcscollector/internal/model"
)

// SyncMemStorage Реализация интерфейса storage для памяти с сохранением данных после каждой записи, использует MemStorage.
type SyncMemStorage struct {
	s        *MemStorage
	syncFile string
}

func NewSyncMemStorage(syncFile string) *SyncMemStorage {
	return &SyncMemStorage{
		syncFile: syncFile,
		s:        NewMemStorage(),
	}
}

func (s *SyncMemStorage) Restore(filePath string) error {
	return s.s.Restore(filePath)
}

func (s *SyncMemStorage) Save(filePath string) error {
	return s.s.Save(filePath)
}

func (s *SyncMemStorage) UpdateCounterMetric(ctx context.Context, m model.CounterMetric) (int64, error) {
	res, err := s.s.UpdateCounterMetric(ctx, m)
	if err != nil {
		return 0, err
	}

	if err = s.s.Save(s.syncFile); err != nil {
		return 0, err
	}

	return res, nil
}

func (s *SyncMemStorage) UpdateGaugeMetric(ctx context.Context, m model.GaugeMetric) error {
	err := s.s.UpdateGaugeMetric(ctx, m)
	if err != nil {
		return err
	}

	return s.s.Save(s.syncFile)
}

func (s *SyncMemStorage) GetGaugeMetric(ctx context.Context, name string) (model.GaugeMetric, error) {
	return s.s.GetGaugeMetric(ctx, name)
}

func (s *SyncMemStorage) GetCounterMetric(ctx context.Context, name string) (model.CounterMetric, error) {
	return s.s.GetCounterMetric(ctx, name)
}

func (s *SyncMemStorage) GetAllGaugeMetrics(ctx context.Context) ([]model.GaugeMetric, error) {
	return s.s.GetAllGaugeMetrics(ctx)
}

func (s *SyncMemStorage) GetAllCounterMetrics(ctx context.Context) ([]model.CounterMetric, error) {
	return s.s.GetAllCounterMetrics(ctx)
}

func (s *SyncMemStorage) UpdateMetrics(ctx context.Context, metricsData model.MetricsData) error {
	err := s.s.UpdateMetrics(ctx, metricsData)
	if err != nil {
		return err
	}

	if err = s.s.Save(s.syncFile); err != nil {
		return err
	}

	return err
}
