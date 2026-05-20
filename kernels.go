package kde

import "math"

// Kernel represents a kernel function for KDE.
type Kernel func(u float64) float64

// Gaussian kernel: K(u) = (1/sqrt(2*pi)) * exp(-0.5 * u^2)
func Gaussian(u float64) float64 {
	return math.Exp(-0.5*u*u) / math.Sqrt(2*math.Pi)
}

// Epanechnikov kernel: K(u) = 0.75 * (1 - u^2) for |u| <= 1, 0 otherwise
func Epanechnikov(u float64) float64 {
	if u > 1 || u < -1 {
		return 0
	}
	return 0.75 * (1 - u*u)
}

// Uniform kernel: K(u) = 0.5 for |u| <= 1, 0 otherwise
func Uniform(u float64) float64 {
	if u > 1 || u < -1 {
		return 0
	}
	return 0.5
}

// KernelFunc returns a Kernel function by name.
func KernelFunc(name string) (Kernel, error) {
	switch name {
	case "gaussian":
		return Gaussian, nil
	case "epanechnikov":
		return Epanechnikov, nil
	case "uniform":
		return Uniform, nil
	default:
		return nil, ErrUnknownKernel
	}
}

// AvailableKernels returns a list of available kernel names.
func AvailableKernels() []string {
	return []string{"gaussian", "epanechnikov", "uniform"}
}
