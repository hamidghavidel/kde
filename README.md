# KDE - Kernel Density Estimation in Go

A pure Go implementation of Kernel Density Estimation (KDE).

## Installation

```bash
go get github.com/hamidghavidel/kde
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/hamidghavidel/kde"
)

func main() {
    // Sample data
    data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}

    // Create KDE with default settings (Gaussian kernel, Silverman's bandwidth)
    k, err := kde.New(data)
    if err != nil {
        panic(err)
    }

    // Evaluate density at specific points
    points := []float64{2.5, 5.0, 7.5}
    densities := k.Evaluate(points)
    fmt.Println("Densities:", densities)

    // Or evaluate on a grid
    x, y := k.EvaluateGrid(0, 11, 100)
    fmt.Printf("Grid: %d points from %.2f to %.2f\n", len(x), x[0], x[len(x)-1])
}
```

## Features

- **Multiple kernels**: Gaussian, Epanechnikov, Uniform
- **Automatic bandwidth selection**: Silverman's rule, Scott's rule, or manual
- **Pure Go**: No C dependencies, easy to deploy
- **Simple API**: Functional options for configuration

## Kernels

Available kernels:
- `gaussian` (default)
- `epanechnikov`
- `uniform`

```go
// Use a specific kernel
k, _ := kde.New(data, kde.WithKernelName("epanechnikov"))
```

## Bandwidth Selection

```go
// Automatic (Silverman's rule - default)
k, _ := kde.New(data)

// Automatic with Scott's rule
k, _ := kde.New(data, kde.WithBandwidthMethod(kde.Scott))

// Manual
k, _ := kde.New(data, kde.WithBandwidth(0.5))
```

## API Reference

### New(data []float64, opts ...Option) (*KDE, error)

Create a new KDE estimator.

Options:
- `WithKernel(k Kernel)` - Set kernel function
- `WithKernelName(name string)` - Set kernel by name
- `WithBandwidth(bw float64)` - Set bandwidth manually
- `WithBandwidthMethod(m BandwidthMethod)` - Set bandwidth selection method

### (k *KDE) Evaluate(points []float64) []float64

Evaluate density at the given points.

### (k *KDE) EvaluateGrid(min, max float64, numPoints int) ([]float64, []float64)

Evaluate density on a grid from min to max with numPoints.

### (k *KDE) Probability(a, b float64, numSamples int) float64

Compute probability that a value falls within [a, b].

### (k *KDE) CDF(x float64, numSamples int) float64

Compute cumulative distribution function at x.

## License

MIT
