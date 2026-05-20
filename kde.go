package kde

import "math"

// Option is a functional option for configuring KDE.
type Option func(*KDE)

// WithKernel sets the kernel function for KDE.
func WithKernel(kernel Kernel) Option {
	return func(k *KDE) {
		k.kernel = kernel
	}
}

// WithKernelType sets the kernel type for KDE.
func WithKernelType(t KernelType) Option {
	return func(k *KDE) {
		k.kernel = t.KernelFunc()
	}
}

// WithBandwidth sets the bandwidth manually.
func WithBandwidth(bw float64) Option {
	return func(k *KDE) {
		k.bandwidth = bw
		k.bandwidthMethod = Manual
	}
}

// WithBandwidthMethod sets the bandwidth selection method.
func WithBandwidthMethod(method BandwidthMethod) Option {
	return func(k *KDE) {
		k.bandwidthMethod = method
	}
}

// WithWeights sets the weights for each data point.
func WithWeights(weights []float64) Option {
	return func(k *KDE) {
		if len(weights) == len(k.data) {
			k.weights = weights
		}
	}
}

// KDE represents a kernel density estimator.
type KDE struct {
	data           []float64
	weights        []float64
	kernel         Kernel
	bandwidth      float64
	bandwidthMethod BandwidthMethod
}

// New creates a new KDE with the given data and options.
func New(data []float64, opts ...Option) (*KDE, error) {
	if len(data) == 0 {
		return nil, ErrEmptyData
	}

	k := &KDE{
		data:            data,
		kernel:          Gaussian.KernelFunc(),
		bandwidthMethod: Silverman,
	}

	for _, opt := range opts {
		opt(k)
	}

	// Compute bandwidth if not manually set
	if k.bandwidthMethod == Manual && k.bandwidth == 0 {
		// User set method to Manual but didn't provide bandwidth
		// Fall back to computing it
	}

	// Compute bandwidth if not already set
	if k.bandwidth == 0 {
		k.bandwidth = computeBandwidth(data, k.bandwidthMethod, k.bandwidth)
	}

	return k, nil
}

// Evaluate computes the density at the given evaluation points.
func (k *KDE) Evaluate(points []float64) []float64 {
	if k.bandwidth == 0 || len(points) == 0 {
		return make([]float64, len(points))
	}

	densities := make([]float64, len(points))
	n := float64(len(k.data))

	for i, p := range points {
		sum := 0.0
		for j, d := range k.data {
			u := (p - d) / k.bandwidth
			kernelVal := k.kernel(u)
			weight := 1.0
			if k.weights != nil && j < len(k.weights) {
				weight = k.weights[j]
			}
			sum += weight * kernelVal
		}
		densities[i] = sum / (n * k.bandwidth)
	}

	return densities
}

// Bandwidth returns the computed bandwidth.
func (k *KDE) Bandwidth() float64 {
	return k.bandwidth
}

// Kernel returns the kernel function.
func (k *KDE) Kernel() Kernel {
	return k.kernel
}

// Data returns the training data.
func (k *KDE) Data() []float64 {
	return k.data
}

// EvaluateGrid evaluates the density on a grid.
func (k *KDE) EvaluateGrid(min, max float64, numPoints int) ([]float64, []float64) {
	if numPoints <= 0 {
		return nil, nil
	}

	step := (max - min) / float64(numPoints-1)
	points := make([]float64, numPoints)
	for i := 0; i < numPoints; i++ {
		points[i] = min + float64(i)*step
	}

	densities := k.Evaluate(points)
	return points, densities
}

// Probability returns the probability that a random variable
// drawn from the KDE falls within the given range [a, b].
func (k *KDE) Probability(a, b float64, numSamples int) float64 {
	if numSamples <= 0 {
		return 0
	}

	// Use midpoint rule for numerical integration
	step := (b - a) / float64(numSamples)
	sum := 0.0

	for i := 0; i < numSamples; i++ {
		x := a + float64(i)*step + step/2
		points := []float64{x}
		density := k.Evaluate(points)[0]
		sum += density * step
	}

	return sum
}

// CDF computes the cumulative distribution function at point x.
func (k *KDE) CDF(x float64, numSamples int) float64 {
	// Find the minimum of the data to use as lower bound
	minVal := math.Inf(1)
	maxVal := math.Inf(-1)
	for _, v := range k.data {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	// Use a slightly wider range
	margin := (maxVal - minVal) * 0.1
	return k.Probability(minVal-margin, x, numSamples)
}
