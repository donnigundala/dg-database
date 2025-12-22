package dgdatabase

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	instrumentationName = "github.com/donnigundala/dg-database"
)

// RegisterMetrics registers database connection pool metrics with OpenTelemetry.
// This should be called after the manager is started.
func (m *Manager) RegisterMetrics() error {
	meter := otel.GetMeterProvider().Meter(instrumentationName)

	// Callback to collect metrics
	callback := func(ctx context.Context, o metric.Observer) error {
		statsMap := m.AllStats()

		for name, stats := range statsMap {
			attrs := metric.WithAttributes(
				attribute.String("db.connection.name", name),
			)

			o.ObserveInt64(m.metricOpen, int64(stats.OpenConnections), attrs)
			o.ObserveInt64(m.metricInUse, int64(stats.InUse), attrs)
			o.ObserveInt64(m.metricIdle, int64(stats.Idle), attrs)
			o.ObserveInt64(m.metricWaitCount, stats.WaitCount, attrs)
			o.ObserveFloat64(m.metricWaitDuration, float64(stats.WaitDuration.Milliseconds()), attrs)
		}
		return nil
	}

	// Create instruments
	var err error
	m.metricOpen, err = meter.Int64ObservableGauge(
		"db.client.connections.open",
		metric.WithDescription("The number of established connections"),
	)
	if err != nil {
		return err
	}

	m.metricInUse, err = meter.Int64ObservableGauge(
		"db.client.connections.in_use",
		metric.WithDescription("The number of connections currently in use"),
	)
	if err != nil {
		return err
	}

	m.metricIdle, err = meter.Int64ObservableGauge(
		"db.client.connections.idle",
		metric.WithDescription("The number of idle connections"),
	)
	if err != nil {
		return err
	}

	m.metricWaitCount, err = meter.Int64ObservableGauge(
		"db.client.connections.wait_count",
		metric.WithDescription("The total number of connections waited for"),
	)
	if err != nil {
		return err
	}

	m.metricWaitDuration, err = meter.Float64ObservableGauge(
		"db.client.connections.wait_duration",
		metric.WithDescription("The total time blocked waiting for a new connection (ms)"),
	)
	if err != nil {
		return err
	}

	// Register callback
	_, err = meter.RegisterCallback(callback, m.metricOpen, m.metricInUse, m.metricIdle, m.metricWaitCount, m.metricWaitDuration)
	return err
}
