package telemetry

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	once    sync.Once
	metrics *PipelineMetrics
)

// PipelineMetrics encapsulates all metrics for the neural pipeline
type PipelineMetrics struct {
	// Allocator metrics
	memoryAllocationHistogram *prometheus.HistogramVec
	allocationErrorsTotal     prometheus.Counter

	// Pipeline performance metrics
	batchLatencyHistogram *prometheus.HistogramVec
	batchSizeGauge        *prometheus.GaugeVec
	throughputCounter     *prometheus.CounterVec

	// Neural graph metrics
	graphNodeCount       *prometheus.GaugeVec
	layerDimensionsGauge *prometheus.GaugeVec

	// Plugin metrics
	pluginStatusGauge   *prometheus.GaugeVec
	pluginErrorsCounter *prometheus.CounterVec
	pluginMetricsGauge  *prometheus.GaugeVec
}

func GetMetrics() *PipelineMetrics {
	once.Do(func() {
		metrics = &PipelineMetrics{
			memoryAllocationHistogram: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "neural_memory_allocation_bytes",
					Help:    "Memory allocation size in bytes",
					Buckets: prometheus.ExponentialBuckets(1024, 2, 20),
				},
				[]string{"allocator_id", "tensor_type"},
			),

			allocationErrorsTotal: promauto.NewCounter(prometheus.CounterOpts{
				Name: "neural_allocation_errors_total",
				Help: "Total number of memory allocation errors",
			}),

			batchLatencyHistogram: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "neural_batch_latency_seconds",
					Help:    "Processing latency per batch",
					Buckets: prometheus.LinearBuckets(0.001, 0.005, 20),
				},
				[]string{"pipeline_id"},
			),

			batchSizeGauge: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "neural_batch_size",
					Help: "Current batch size",
				},
				[]string{"pipeline_id"},
			),

			throughputCounter: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "neural_batches_processed_total",
					Help: "Total number of batches processed",
				},
				[]string{"pipeline_id", "status"},
			),

			graphNodeCount: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "neural_graph_nodes",
					Help: "Number of nodes in neural graph",
				},
				[]string{"graph_id", "node_type"},
			),

			layerDimensionsGauge: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "neural_layer_dimensions",
					Help: "Dimensions of neural network layers",
				},
				[]string{"layer_id", "dimension_type"},
			),

			// New plugin metrics
			pluginStatusGauge: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "plugin_status",
					Help: "Current status of plugins (1=healthy/active, 0=unhealthy/inactive)",
				},
				[]string{"plugin_id", "status"},
			),

			pluginErrorsCounter: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "plugin_errors_total",
					Help: "Total number of plugin errors by type",
				},
				[]string{"plugin_id", "error_type"},
			),

			pluginMetricsGauge: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "plugin_metrics",
					Help: "Custom plugin metrics",
				},
				[]string{"plugin_id", "metric_name"},
			),
		}
	})
	return metrics
}

// ObserveAllocation records a memory allocation
func (m *PipelineMetrics) ObserveAllocation(allocatorID string, tensorType string, bytes float64) {
	m.memoryAllocationHistogram.WithLabelValues(allocatorID, tensorType).Observe(bytes)
}

// RecordAllocationError increments allocation error counter
func (m *PipelineMetrics) RecordAllocationError() {
	m.allocationErrorsTotal.Inc()
}

// ObserveBatchLatency records batch processing latency
func (m *PipelineMetrics) ObserveBatchLatency(pipelineID string, duration time.Duration) {
	m.batchLatencyHistogram.WithLabelValues(pipelineID).Observe(duration.Seconds())
}

// SetBatchSize updates the current batch size gauge
func (m *PipelineMetrics) SetBatchSize(pipelineID string, size float64) {
	m.batchSizeGauge.WithLabelValues(pipelineID).Set(size)
}

// IncrementBatchCount increments the batch counter
func (m *PipelineMetrics) IncrementBatchCount(pipelineID string, status string) {
	m.throughputCounter.WithLabelValues(pipelineID, status).Inc()
}

// SetGraphNodes updates neural graph node count
func (m *PipelineMetrics) SetGraphNodes(graphID string, nodeType string, count float64) {
	m.graphNodeCount.WithLabelValues(graphID, nodeType).Set(count)
}

// SetLayerDimension updates layer dimension gauge
func (m *PipelineMetrics) SetLayerDimension(layerID string, dimensionType string, value float64) {
	m.layerDimensionsGauge.WithLabelValues(layerID, dimensionType).Set(value)
}

// Plugin-specific metric methods

// SetPluginStatus updates plugin status gauge
func (m *PipelineMetrics) SetPluginStatus(pluginID string, status string, value float64) {
	m.pluginStatusGauge.WithLabelValues(pluginID, status).Set(value)
}

// RecordPluginError increments plugin error counter
func (m *PipelineMetrics) RecordPluginError(pluginID string, errorType string) {
	m.pluginErrorsCounter.WithLabelValues(pluginID, errorType).Inc()
}

// SetPluginMetric sets a custom plugin metric
func (m *PipelineMetrics) SetPluginMetric(pluginID string, metricName string, value interface{}) {
	// Convert string value to float64 if needed
	var floatValue float64
	switch v := value.(type) {
	case float64:
		floatValue = v
	case int:
		floatValue = float64(v)
	case string:
		// Try to maintain backward compatibility by setting 1.0 for non-empty strings
		if v != "" {
			floatValue = 1.0
		}
	default:
		// For unsupported types, set to 0
		floatValue = 0
	}

	m.pluginMetricsGauge.WithLabelValues(pluginID, metricName).Set(floatValue)
}
