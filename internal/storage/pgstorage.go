package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/smakimka/mtrcscollector/internal/model"
	"github.com/smakimka/mtrcscollector/internal/retry"
)

type PGStorage struct {
	p *pgxpool.Pool
}

func NewPGStorage(ctx context.Context, p *pgxpool.Pool) (PGStorage, error) {
	s := PGStorage{
		p: p,
	}

	err := s.CreateSchemaIFNotExists(ctx)
	if err != nil {
		return s, err
	}

	return s, nil
}

func (s PGStorage) Ping(ctx context.Context) error {
	return s.p.Ping(ctx)
}

func (s PGStorage) CreateSchemaIFNotExists(ctx context.Context) error {
	_, err := retry.Exec(s.p.Exec, ctx, `create table if not exists counter_metrics (
		id serial primary key,
		name text,
		value integer
	)`)
	if err != nil {
		return err
	}

	_, err = retry.Exec(s.p.Exec, ctx, `create table if not exists gauge_metrics (
		id serial primary key,
		name text,
		value double precision
	)`)
	if err != nil {
		return err
	}

	return nil
}

func (s PGStorage) UpdateCounterMetric(ctx context.Context, m model.CounterMetric) (int64, error) {
	row := s.p.QueryRow(ctx, "select id, value from counter_metrics where name like $1", m.Name)
	var id int16
	var value int64

	err := row.Scan(&id, &value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_, err := retry.Exec(s.p.Exec, ctx, "insert into counter_metrics (name, value) values ($1, $2)", m.Name, m.Value)
			if err != nil {
				return 0, err
			}
			return m.Value, nil
		} else {
			return 0, err
		}
	}

	newVal := m.Value + value
	_, err = retry.Exec(s.p.Exec, ctx, "update counter_metrics set value = $1 where id = $2", newVal, id)
	if err != nil {
		return 0, err
	}

	return newVal, nil
}

func (s PGStorage) UpdateGaugeMetric(ctx context.Context, m model.GaugeMetric) error {
	res, err := retry.Exec(s.p.Exec, ctx, "update gauge_metrics set value = $1 where name like $2", m.Value, m.Name)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		_, err := retry.Exec(s.p.Exec, ctx, "insert into gauge_metrics (name, value) values ($1, $2)", m.Name, m.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s PGStorage) GetGaugeMetric(ctx context.Context, name string) (model.GaugeMetric, error) {
	var m model.GaugeMetric
	row := s.p.QueryRow(ctx, "select name, value from gauge_metrics where name like $1", name)

	err := row.Scan(&m.Name, &m.Value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return m, ErrNoSuchMetric
		} else {
			return m, err
		}
	}

	return m, nil
}

func (s PGStorage) GetCounterMetric(ctx context.Context, name string) (model.CounterMetric, error) {
	var m model.CounterMetric
	row := s.p.QueryRow(ctx, "select name, value from counter_metrics where name like $1", name)

	err := row.Scan(&m.Name, &m.Value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return m, ErrNoSuchMetric
		} else {
			return m, err
		}
	}

	return m, nil
}

func (s PGStorage) GetAllGaugeMetrics(ctx context.Context) ([]model.GaugeMetric, error) {
	rows, err := retry.Query(s.p.Query, ctx, "select name, value from gauge_metrics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []model.GaugeMetric
	for rows.Next() {
		var m model.GaugeMetric
		err := rows.Scan(&m.Name, &m.Value)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, m)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return metrics, nil
}

func (s PGStorage) GetAllCounterMetrics(ctx context.Context) ([]model.CounterMetric, error) {
	rows, err := retry.Query(s.p.Query, ctx, "select name, value from counter_metrics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []model.CounterMetric
	for rows.Next() {
		var m model.CounterMetric
		err := rows.Scan(&m.Name, &m.Value)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, m)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return metrics, nil
}

func (s PGStorage) UpdateMetrics(ctx context.Context, metricsData model.MetricsData) error {
	tx, err := s.p.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, metricData := range metricsData {
		switch metricData.Kind {
		case model.Gauge:
			err := txUpdateGaugeMetric(ctx, tx, metricData)
			if err != nil {
				return err
			}
		case model.Counter:
			_, err := txUpdateCounterMetric(ctx, tx, metricData)
			if err != nil {
				return err
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func txUpdateGaugeMetric(ctx context.Context, tx pgx.Tx, metricData model.MetricData) error {
	res, err := retry.Exec(tx.Exec, ctx, "update gauge_metrics set value = $1 where name like $2", metricData.Value, metricData.Name)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		_, err := retry.Exec(tx.Exec, ctx, "insert into gauge_metrics (name, value) values ($1, $2)", metricData.Name, metricData.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

func txUpdateCounterMetric(ctx context.Context, tx pgx.Tx, metricData model.MetricData) (int64, error) {
	row := tx.QueryRow(ctx, "select id, value from counter_metrics where name like $1", metricData.Name)
	var id int16
	var value int64

	err := row.Scan(&id, &value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_, err := retry.Exec(tx.Exec, ctx, "insert into counter_metrics (name, value) values ($1, $2)", metricData.Name, *metricData.Delta)
			if err != nil {
				return 0, err
			}
			return *metricData.Delta, nil
		} else {
			return 0, err
		}
	}

	newVal := *metricData.Delta + value
	_, err = retry.Exec(tx.Exec, ctx, "update counter_metrics set value = $1 where id = $2", newVal, id)
	if err != nil {
		return 0, err
	}

	return newVal, nil
}
