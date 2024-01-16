package storage

import "github.com/smakimka/mtrcscollector/internal/model"

type SyncMemStorage struct {
	syncFile string
	s        *MemStorage
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

func (s *SyncMemStorage) UpdateCounterMetric(m model.CounterMetric) (int64, error) {
	res, err := s.s.UpdateCounterMetric(m)
	if err != nil {
		return 0, err
	}

	if err = s.s.Save(s.syncFile); err != nil {
		return 0, err
	}

	return res, nil
}

func (s *SyncMemStorage) UpdateGaugeMetric(m model.GaugeMetric) error {
	err := s.s.UpdateGaugeMetric(m)
	if err != nil {
		return err
	}

	return s.s.Save(s.syncFile)
}

func (s *SyncMemStorage) GetGaugeMetric(name string) (model.GaugeMetric, error) {
	return s.s.GetGaugeMetric(name)
}

func (s *SyncMemStorage) GetCounterMetric(name string) (model.CounterMetric, error) {
	return s.s.GetCounterMetric(name)
}

func (s *SyncMemStorage) GetAllGaugeMetrics() ([]model.GaugeMetric, error) {
	return s.s.GetAllGaugeMetrics()
}

func (s *SyncMemStorage) GetAllCounterMetrics() ([]model.CounterMetric, error) {
	return s.s.GetAllCounterMetrics()
}
