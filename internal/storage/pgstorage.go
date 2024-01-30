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

	err := s.CreateSchemaIfNotExists(ctx)
	if err != nil {
		return s, err
	}

	return s, nil
}

func (s PGStorage) Ping(ctx context.Context) error {
	return s.p.Ping(ctx)
}

func (s PGStorage) CreateSchemaIfNotExists(ctx context.Context) error {
	_, err := retry.Exec(s.p.Exec, ctx, `create table if not exists counter_metrics (
		id serial primary key,
		name text,
		value bigint,
		constraint c_name_uq unique (name)
	)`)
	if err != nil {
		return err
	}

	_, err = retry.Exec(s.p.Exec, ctx, `create table if not exists gauge_metrics (
		id serial primary key,
		name text,
		value double precision,
		constraint g_name_uq unique (name)
	)`)
	if err != nil {
		return err
	}

	return nil
}

func (s PGStorage) UpdateCounterMetric(ctx context.Context, m model.CounterMetric) (int64, error) {
	tx, err := s.p.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, `insert into counter_metrics as cm (name, value) values ($1, $2) 
							on conflict on constraint c_name_uq do update set value = cm.value + $2
							returning cm.value`, m.Name, m.Value)

	var value int64
	err = row.Scan(&value)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return value, nil
}

func (s PGStorage) UpdateGaugeMetric(ctx context.Context, m model.GaugeMetric) error {
	tx, err := s.p.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = retry.Exec(tx.Exec, ctx, `insert into gauge_metrics (name, value) values ($1, $2) 
										  on conflict on constraint g_name_uq do update set value = $2`, m.Name, m.Value)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
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
	_, err := retry.Exec(tx.Exec, ctx, `insert into gauge_metrics (name, value) values ($1, $2) 
										  on conflict on constraint g_name_uq do update set value = $2`, metricData.Name, metricData.Value)
	if err != nil {
		return err
	}

	return nil
}

func txUpdateCounterMetric(ctx context.Context, tx pgx.Tx, metricData model.MetricData) (int64, error) {
	row := tx.QueryRow(ctx, `insert into counter_metrics as cm (name, value) values ($1, $2) 
							on conflict on constraint c_name_uq do update set value = cm.value + $2
							returning cm.value`, metricData.Name, metricData.Delta)

	var value int64
	err := row.Scan(&value)
	if err != nil {
		return 0, err
	}

	return value, nil
}
