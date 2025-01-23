package generator

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"threshAI/internal/core/memory"
	"threshAI/internal/telemetry"
)

const (
	maxBatchSize    = 32
	minBatchSize    = 4
	batchTimeout    = 50 * time.Millisecond
	warmupSteps     = 100
	profilingWindow = 1000
)

type OptimizedPipeline struct {
	allocator      *memory.OptimizedAllocator
	batchSize      int32
	throughput     atomic.Int64
	latencySum     atomic.Int64
	batchesCounter atomic.Int32
	warmupComplete atomic.Bool
	metrics        *telemetry.PipelineMetrics
	pipelineID     string
}

func NewOptimizedPipeline(id string) *OptimizedPipeline {
	p := &OptimizedPipeline{
		allocator:  memory.NewOptimizedAllocator(),
		batchSize:  minBatchSize,
		pipelineID: id,
		metrics:    telemetry.GetMetrics(),
	}
	p.allocator.EnableAdaptiveCheckpointing()
	p.metrics.SetBatchSize(id, float64(minBatchSize))
	return p
}

// Generate implements batched generation with adaptive sizing and graduated pressure handling
func (p *OptimizedPipeline) Generate(ctx context.Context, prompt string) (string, error) {
	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, batchTimeout)
	defer cancel()

	currentBatch := atomic.LoadInt32(&p.batchSize)
	batchCounter := p.batchesCounter.Add(1)

	// Profile and adjust batch size with warmup
	if batchCounter > warmupSteps {
		p.warmupComplete.Store(true)
		if batchCounter%profilingWindow == 0 {
			p.adjustBatchSize()
		}
	}

	// Calculate allocation size with padding for potential growth
	baseSize := len(prompt)
	growthPadding := int(float64(baseSize) * 0.1) // 10% padding
	allocSize := (baseSize + growthPadding) * int(currentBatch)

	// Try allocation with pressure monitoring
	tensor, err := p.allocator.AllocTensorWithCache(allocSize)
	if err != nil {
		p.metrics.RecordAllocationError()

		// Implement graduated fallback strategy
		if currentBatch > minBatchSize*2 {
			// Try halving the batch size first
			newBatch := currentBatch / 2
			atomic.StoreInt32(&p.batchSize, newBatch)
			p.metrics.SetBatchSize(p.pipelineID, float64(newBatch))

			// Retry with smaller batch
			allocSize = (baseSize + growthPadding) * int(newBatch)
			tensor, err = p.allocator.AllocTensorWithCache(allocSize)
			if err != nil {
				// If still failing, fall back to minimum
				atomic.StoreInt32(&p.batchSize, minBatchSize)
				p.metrics.SetBatchSize(p.pipelineID, float64(minBatchSize))
				return "", err
			}
		} else {
			// Already at or near minimum, simply fail
			atomic.StoreInt32(&p.batchSize, minBatchSize)
			p.metrics.SetBatchSize(p.pipelineID, float64(minBatchSize))
			return "", err
		}
	}

	// Record successful allocation
	p.metrics.ObserveAllocation(p.pipelineID, "prompt_tensor", float64(allocSize))

	// Execute generation with timeout handling
	start := time.Now()
	resultChan := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		defer func() {
			// Free tensor when generation is complete
			p.allocator.FreeTensor(tensor)
		}()

		// Placeholder for actual generation using the tensor
		// In real implementation, this would use the tensor for computation
		result := fmt.Sprintf("Generated (tensor:%v): %s", tensor, prompt)
		resultChan <- result
	}()

	// Wait for completion or timeout
	select {
	case result := <-resultChan:
		// Update performance metrics
		latency := time.Since(start)
		p.latencySum.Add(latency.Microseconds())
		p.throughput.Add(1)

		// Record telemetry
		p.metrics.ObserveBatchLatency(p.pipelineID, latency)
		p.metrics.IncrementBatchCount(p.pipelineID, "success")

		return result, nil
	case err := <-errChan:
		p.metrics.IncrementBatchCount(p.pipelineID, "error")
		return "", err
	case <-timeoutCtx.Done():
		p.metrics.IncrementBatchCount(p.pipelineID, "timeout")
		return "", timeoutCtx.Err()
	}
}

// adjustBatchSize implements dynamic batch sizing based on performance metrics
func (p *OptimizedPipeline) adjustBatchSize() {
	if !p.warmupComplete.Load() {
		return
	}

	currentBatch := atomic.LoadInt32(&p.batchSize)
	avgLatency := float64(p.latencySum.Load()) / float64(p.throughput.Load())
	timeoutUs := float64(batchTimeout.Microseconds())

	// Get memory metrics from allocator
	allocStats := p.allocator.GetMetrics()
	avgPressure := (allocStats["small_pool_pressure"].(float64) +
		allocStats["medium_pool_pressure"].(float64) +
		allocStats["large_pool_pressure"].(float64)) / 3.0

	highWatermark := allocStats["high_watermark"].(float64)
	criticalWatermark := allocStats["critical_watermark"].(float64)

	// Calculate target batch size based on both latency and memory pressure
	var targetBatch int32

	// Fast path - good performance and low memory pressure
	if avgLatency < timeoutUs*0.5 && avgPressure < highWatermark*0.6 && currentBatch < maxBatchSize {
		// Room for aggressive growth
		growthStep := int32(4)
		if avgPressure > highWatermark*0.4 {
			growthStep = 2 // More conservative growth under moderate pressure
		}
		targetBatch = currentBatch + growthStep
	} else if avgLatency > timeoutUs*1.5 || avgPressure > criticalWatermark*0.9 {
		// Problems detected - reduce batch size
		reductionFactor := int32(1)
		if avgLatency > timeoutUs*2 || avgPressure > criticalWatermark {
			// Severe issues - reduce more aggressively
			reductionFactor = currentBatch / 4
		} else {
			reductionFactor = currentBatch / 8
		}
		if reductionFactor < 1 {
			reductionFactor = 1
		}
		targetBatch = currentBatch - reductionFactor
	} else {
		// Stable performance - fine-tune
		if avgLatency < timeoutUs*0.8 && avgPressure < highWatermark*0.8 {
			targetBatch = currentBatch + 1
		} else if avgLatency > timeoutUs*1.2 || avgPressure > highWatermark {
			targetBatch = currentBatch - 1
		} else {
			targetBatch = currentBatch // Maintain current size
		}
	}

	// Apply bounds
	if targetBatch > maxBatchSize {
		targetBatch = maxBatchSize
	}
	if targetBatch < minBatchSize {
		targetBatch = minBatchSize
	}

	// Update batch size if changed
	if targetBatch != currentBatch {
		atomic.StoreInt32(&p.batchSize, targetBatch)
		p.metrics.SetBatchSize(p.pipelineID, float64(targetBatch))
	}

	// Update metrics
	p.metrics.SetLayerDimension(p.pipelineID, "batch_size", float64(targetBatch))
	p.metrics.SetLayerDimension(p.pipelineID, "avg_latency_us", avgLatency)
	p.metrics.SetLayerDimension(p.pipelineID, "memory_pressure", avgPressure)

	// Reset performance counters
	p.latencySum.Store(0)
	p.throughput.Store(0)
}

// GetMetrics returns current pipeline performance metrics
func (p *OptimizedPipeline) GetMetrics() map[string]interface{} {
	currentBatch := atomic.LoadInt32(&p.batchSize)
	totalBatches := p.batchesCounter.Load()
	avgLatency := float64(p.latencySum.Load()) / float64(p.throughput.Load())

	// Update prometheus metrics
	p.metrics.SetBatchSize(p.pipelineID, float64(currentBatch))
	p.metrics.SetLayerDimension(p.pipelineID, "total_batches", float64(totalBatches))

	return map[string]interface{}{
		"batch_size":      currentBatch,
		"total_batches":   totalBatches,
		"warmup_complete": p.warmupComplete.Load(),
		"avg_latency_us":  avgLatency,
		"pipeline_id":     p.pipelineID,
	}
}
