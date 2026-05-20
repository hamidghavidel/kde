package kde

import "math"

// BandwidthMethod represents a method for computing bandwidth.
type BandwidthMethod string

const (
	// Silverman uses Silverman's rule of thumb for bandwidth selection.
	Silverman BandwidthMethod = "silverman"
	// Scott uses Scott's rule of thumb for bandwidth selection.
	Scott BandwidthMethod = "scott"
	// Manual means the bandwidth is user-provided.
	Manual BandwidthMethod = "manual"
)

// computeBandwidth calculates the bandwidth using the specified method.
func computeBandwidth(data []float64, method BandwidthMethod, manualBw float64) float64 {
	n := float64(len(data))
	if n == 0 {
		return 0
	}

	switch method {
	case Manual:
		return manualBw
	case Silverman:
		// Silverman's rule: h = 0.9 * min(sd, IQR/1.34) * n^(-1/5)
		sd := standardDeviation(data)
		iqr := interquartileRange(data)
		bandwidth := 0.9 * math.Min(sd, iqr/1.34) * math.Pow(n, -0.2)
		return bandwidth
	case Scott:
		// Scott's rule: h = 1.059 * sd * n^(-1/5)
		sd := standardDeviation(data)
		bandwidth := 1.059 * sd * math.Pow(n, -0.2)
		return bandwidth
	default:
		// Default to Silverman
		sd := standardDeviation(data)
		iqr := interquartileRange(data)
		return 0.9 * math.Min(sd, iqr/1.34) * math.Pow(n, -0.2)
	}
}

// standardDeviation computes the standard deviation of the data.
func standardDeviation(data []float64) float64 {
	n := float64(len(data))
	if n == 0 {
		return 0
	}

	mean := 0.0
	for _, v := range data {
		mean += v
	}
	mean /= n

	variance := 0.0
	for _, v := range data {
		diff := v - mean
		variance += diff * diff
	}
	variance /= n - 1 // Sample variance

	return math.Sqrt(variance)
}

// interquartileRange computes the IQR (Q3 - Q1) of the data.
func interquartileRange(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}

	// Sort a copy of the data
	sorted := make([]float64, len(data))
	copy(sorted, data)
	quickSort(sorted, 0, len(sorted)-1)

	n := len(sorted)
	q1Index := n / 4
	q3Index := 3 * n / 4

	return sorted[q3Index] - sorted[q1Index]
}

// quickSort sorts the data in place using quicksort.
func quickSort(data []float64, low, high int) {
	if low < high {
		pivotIndex := partition(data, low, high)
		quickSort(data, low, pivotIndex-1)
		quickSort(data, pivotIndex+1, high)
	}
}

// partition partitions the data around a pivot.
func partition(data []float64, low, high int) int {
	pivot := data[high]
	i := low - 1
	for j := low; j < high; j++ {
		if data[j] <= pivot {
			i++
			data[i], data[j] = data[j], data[i]
		}
	}
	data[i+1], data[high] = data[high], data[i+1]
	return i + 1
}
