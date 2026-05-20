package kde

import "math"

// KernelType represents the type of kernel function used in KDE.
type KernelType int

const (
	// Gaussian kernel: K(u) = (1/sqrt(2*pi)) * exp(-0.5 * u^2)
	Gaussian KernelType = iota
	// Epanechnikov kernel: K(u) = 0.75 * (1 - u^2) for |u| <= 1, 0 otherwise
	Epanechnikov
	// Uniform kernel: K(u) = 0.5 for |u| <= 1, 0 otherwise
	Uniform
)

// String returns the string representation of the kernel type.
func (k KernelType) String() string {
	switch k {
	case Gaussian:
		return "gaussian"
	case Epanechnikov:
		return "epanechnikov"
	case Uniform:
		return "uniform"
	default:
		return "unknown"
	}
}

// KernelFunc returns the Kernel function for the given kernel type.
func (k KernelType) KernelFunc() Kernel {
	switch k {
	case Gaussian:
		return gaussianKernel
	case Epanechnikov:
		return epanechnikovKernel
	case Uniform:
		return uniformKernel
	default:
		return gaussianKernel
	}
}

// kernel is the internal function type for kernel functions.
type kernel func(u float64) float64

// Kernel represents a kernel function for KDE.
type Kernel func(u float64) float64

// gaussianKernel is the Gaussian kernel function.
func gaussianKernel(u float64) float64 {
	return math.Exp(-0.5*u*u) / math.Sqrt(2*math.Pi)
}

// epanechnikovKernel is the Epanechnikov kernel function.
func epanechnikovKernel(u float64) float64 {
	if u > 1 || u < -1 {
		return 0
	}
	return 0.75 * (1 - u*u)
}

// uniformKernel is the Uniform kernel function.
func uniformKernel(u float64) float64 {
	if u > 1 || u < -1 {
		return 0
	}
	return 0.5
}

// KernelFunc returns a Kernel function by name.
// Deprecated: Use KernelType constants directly with WithKernel or KernelType.KernelFunc.
func KernelFunc(name string) (Kernel, error) {
	switch name {
	case "gaussian":
		return gaussianKernel, nil
	case "epanechnikov":
		return epanechnikovKernel, nil
	case "uniform":
		return uniformKernel, nil
	default:
		return nil, ErrUnknownKernel
	}
}

// AvailableKernels returns a list of available kernel types.
func AvailableKernels() []KernelType {
	return []KernelType{Gaussian, Epanechnikov, Uniform}
}
