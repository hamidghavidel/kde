package kde

import (
	"math"
	"testing"
)

func TestGaussianKernel(t *testing.T) {
	tests := []struct {
		name     string
		u        float64
		expected float64
	}{
		{"u=0", 0, 1 / math.Sqrt(2*math.Pi)},
		{"u=1", 1, math.Exp(-0.5) / math.Sqrt(2*math.Pi)},
		{"u=-1", -1, math.Exp(-0.5) / math.Sqrt(2*math.Pi)},
		{"u=2", 2, math.Exp(-2) / math.Sqrt(2*math.Pi)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gaussianKernel(tt.u)
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("Gaussian(%v) = %v, want %v", tt.u, got, tt.expected)
			}
		})
	}
}

func TestEpanechnikovKernel(t *testing.T) {
	tests := []struct {
		name     string
		u        float64
		expected float64
	}{
		{"u=0", 0, 0.75},
		{"u=0.5", 0.5, 0.75 * (1 - 0.25)},
		{"u=1", 1, 0},
		{"u=-1", -1, 0},
		{"u=1.5", 1.5, 0},
		{"u=-1.5", -1.5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := epanechnikovKernel(tt.u)
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("Epanechnikov(%v) = %v, want %v", tt.u, got, tt.expected)
			}
		})
	}
}

func TestUniformKernel(t *testing.T) {
	tests := []struct {
		name     string
		u        float64
		expected float64
	}{
		{"u=0", 0, 0.5},
		{"u=0.5", 0.5, 0.5},
		{"u=1", 1, 0.5},
		{"u=-1", -1, 0.5},
		{"u=1.5", 1.5, 0},
		{"u=-1.5", -1.5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uniformKernel(tt.u)
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("Uniform(%v) = %v, want %v", tt.u, got, tt.expected)
			}
		})
	}
}

func TestKernelFunc(t *testing.T) {
	tests := []struct {
		name    string
		kernel  string
		wantErr bool
	}{
		{"gaussian", "gaussian", false},
		{"epanechnikov", "epanechnikov", false},
		{"uniform", "uniform", false},
		{"unknown", "triangular", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := KernelFunc(tt.kernel)
			if (err != nil) != tt.wantErr {
				t.Errorf("KernelFunc(%q) error = %v, wantErr %v", tt.kernel, err, tt.wantErr)
			}
		})
	}
}

func TestAvailableKernels(t *testing.T) {
	kernels := AvailableKernels()
	if len(kernels) != 3 {
		t.Errorf("AvailableKernels() returned %d kernels, want 3", len(kernels))
	}
}

func TestStandardDeviation(t *testing.T) {
	tests := []struct {
		name     string
		data     []float64
		expected float64
	}{
		{"simple", []float64{1, 2, 3}, 1.0},
		{"single", []float64{5}, 0},
		{"empty", []float64{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := standardDeviation(tt.data)
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("standardDeviation(%v) = %v, want %v", tt.data, got, tt.expected)
			}
		})
	}
}

func TestInterquartileRange(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	iqr := interquartileRange(data)
	// For [1,2,3,4,5,6,7,8]: Q1=3, Q3=7, IQR=4
	expected := 4.0
	if math.Abs(iqr-expected) > 1e-10 {
		t.Errorf("interquartileRange(%v) = %v, want %v", data, iqr, expected)
	}
}

func TestComputeBandwidth(t *testing.T) {
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}

	// Test Silverman
	bw := computeBandwidth(data, Silverman, 0)
	if bw <= 0 {
		t.Errorf("computeBandwidth with Silverman returned non-positive bandwidth: %v", bw)
	}

	// Test Scott
	bw = computeBandwidth(data, Scott, 0)
	if bw <= 0 {
		t.Errorf("computeBandwidth with Scott returned non-positive bandwidth: %v", bw)
	}

	// Test Manual
	bw = computeBandwidth(data, Manual, 0.5)
	if bw != 0.5 {
		t.Errorf("computeBandwidth with Manual returned %v, want 0.5", bw)
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		data    []float64
		wantErr bool
	}{
		{"valid", []float64{1, 2, 3, 4, 5}, false},
		{"empty", []float64{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKDEWithOptions(t *testing.T) {
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0}

	// Test with custom bandwidth
	k, err := New(data, WithBandwidth(0.5))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if k.Bandwidth() != 0.5 {
		t.Errorf("Bandwidth() = %v, want 0.5", k.Bandwidth())
	}

	// Test with custom kernel
	k, err = New(data, WithKernelType(Epanechnikov))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if k.Kernel() == nil {
		t.Errorf("Kernel() = nil, want non-nil")
	}

	// Test with Scott's rule
	k, err = New(data, WithBandwidthMethod(Scott))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if k.Bandwidth() <= 0 {
		t.Errorf("Bandwidth() = %v, want positive value", k.Bandwidth())
	}
}

func TestKDEEvaluate(t *testing.T) {
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	k, err := New(data, WithBandwidth(0.5))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Test evaluation
	points := []float64{1.0, 2.5, 5.0}
	densities := k.Evaluate(points)

	if len(densities) != len(points) {
		t.Errorf("Evaluate() returned %d densities, want %d", len(densities), len(points))
	}

	// All densities should be non-negative for valid data
	for i, d := range densities {
		if d < 0 {
			t.Errorf("Evaluate()[%d] = %v, want non-negative", i, d)
		}
	}
}

func TestKDEEvaluateGrid(t *testing.T) {
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	k, err := New(data, WithBandwidth(0.5))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	x, y := k.EvaluateGrid(0, 6, 100)

	if len(x) != 100 || len(y) != 100 {
		t.Errorf("EvaluateGrid() returned wrong lengths: x=%d, y=%d, want 100", len(x), len(y))
	}

	// All densities should be non-negative
	for i, d := range y {
		if d < 0 {
			t.Errorf("EvaluateGrid()[%d] = %v, want non-negative", i, d)
		}
	}
}

func TestKDEProbability(t *testing.T) {
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	k, err := New(data, WithBandwidth(1.0))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Probability should be between 0 and 1
	prob := k.Probability(0, 10, 100)
	if prob < 0 || prob > 1 {
		t.Errorf("Probability() = %v, want between 0 and 1", prob)
	}
}

func TestKDECDF(t *testing.T) {
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	k, err := New(data, WithBandwidth(1.0))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// CDF should be between 0 and 1
	cdf := k.CDF(3.0, 100)
	if cdf < 0 || cdf > 1 {
		t.Errorf("CDF() = %v, want between 0 and 1", cdf)
	}
}
