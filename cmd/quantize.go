package cmd

import (
	"errors"
	"fmt"

	"threshAI/internal/core/memory"

	"github.com/spf13/cobra"
)

var (
	quantizeCmd = &cobra.Command{
		Use:   "quantize",
		Short: "Quantize model weights",
		Long: `Quantize model weights using specified precision.
Supported quantization modes:
- Q4_K_M: 4-bit medium quantization (4.3GB)
- Q5_K_S: 5-bit small quantization (5.1GB)`,
		RunE: runQuantize,
	}

	quantizeMode string
	memoryMap    = map[string]uint64{
		"Q4_K_M": 4617089843, // 4.3GB
		"Q5_K_S": 5476083302, // 5.1GB
	}
)

func init() {
	rootCmd.AddCommand(quantizeCmd)
	quantizeCmd.Flags().StringVarP(&quantizeMode, "quantize", "q", "", "Quantization mode (Q4_K_M, Q5_K_S)")
}

func runQuantize(cmd *cobra.Command, args []string) error {
	if quantizeMode == "" {
		return errors.New("quantization mode must be specified")
	}

	memSize, ok := memoryMap[quantizeMode]
	if !ok {
		return fmt.Errorf("unsupported quantization mode: %s", quantizeMode)
	}

	// Initialize GPU allocator
	allocator := memory.NewGPUMemoryAllocator()
	defer allocator.FreeAll()

	// Allocate memory for quantized weights
	ptr, err := allocator.Alloc4BitTensorCache(int(memSize))
	if err != nil {
		return fmt.Errorf("failed to allocate GPU memory: %v", err)
	}
	_ = ptr // Prevent unused variable error

	// Enable mixed precision
	if err := allocator.EnableMixedPrecision(); err != nil {
		return fmt.Errorf("failed to enable mixed precision: %v", err)
	}

	fmt.Printf("Successfully initialized %s quantization with %d bytes allocated\n",
		quantizeMode, memSize)

	return nil
}
