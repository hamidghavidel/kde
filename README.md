# KDE - Kernel Density Estimation in Go

A pure Go implementation of Kernel Density Estimation (KDE).

## Table of Contents

- [What is Kernel Density Estimation?](#what-is-kernel-density-estimation)
- [How It Works](#how-it-works)
- [The Mathematics](#the-mathematics)
- [Kernels](#kernels)
- [Bandwidth Selection](#bandwidth-selection)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [API Reference](#api-reference)
- [Examples](#examples)

## What is Kernel Density Estimation?

Kernel Density Estimation (KDE) is a non-parametric technique for estimating the probability density function of a dataset. Unlike parametric methods (e.g., assuming data follows a normal distribution), KDE makes no such assumptions—it lets the data speak for itself.

KDE is widely used in:
- Statistical analysis and data science
- Anomaly detection
- Data visualization
- Machine learning (as a density estimator)
- Scientific research

### Why use KDE?

```go
// Parametric: assumes normal distribution
mean, std := calculateMeanAndStd(data)
density := normalPDF(x, mean, std)

// KDE: no distribution assumption
k, _ := kde.New(data)
density := k.Evaluate([]float64{x})[0]
```

## How It Works

Given a dataset {x₁, x₂, ..., xₙ}, KDE estimates the probability density at a point x using:

```
        1   n
f(x) = ----  Σ  K((x - xᵢ) / h)
        n·h  i=1
```

Where:
- **n** is the number of data points
- **h** is the bandwidth (smoothing parameter)
- **K** is the kernel function
- **xᵢ** are the individual data points

The kernel function K determines the shape of the "bump" placed at each data point. The bandwidth h controls how wide each bump is—small h means sharp peaks, large h means smooth curves.

### Intuition

Imagine placing a small mound (kernel) at each data point. The final density estimate is the average height of all these mounds at any given point. The bandwidth controls how spread out each mound is.

## The Mathematics

### The Kernel Function

A kernel function K(u) must satisfy:
1. **Non-negative**: K(u) ≥ 0
2. **Symmetric**: K(u) = K(-u)
3. **Integrates to 1**: ∫K(u)du = 1

The kernel function is evaluated at the scaled distance u = (x - xᵢ) / h from each data point.

### Choosing a Kernel

The choice of kernel affects the smoothness of the estimate. The Gaussian kernel produces the smoothest curves. The Epanechnikov kernel is optimal in terms of mean integrated squared error. All commonly used kernels produce similar results with sufficient data.

### The Bandwidth (h)

The bandwidth is the **most critical parameter** in KDE:

| Bandwidth | Effect |
|-----------|--------|
| Too small (h → 0) | Overfitting—spiky, noisy estimate that follows every fluctuation |
| Too large (h → ∞) | Underfitting—oversmoothed estimate that hides structure |
| Just right | Balances bias and variance to reveal true density shape |

The bandwidth controls the bias-variance tradeoff:
- **Low bandwidth**: Low bias, high variance (captures detail but noisy)
- **High bandwidth**: High bias, low variance (smooth but misses features)

### Mean Integrated Squared Error (MISE)

The quality of a KDE is often measured by MISE:

```
MISE(h) = E[∫(f̂(x) - f(x))² dx]
```

The optimal bandwidth minimizes MISE. Both Silverman's and Scott's rules are approximations to this optimal bandwidth under the assumption that the true density is normal.

## Kernels

### Gaussian Kernel

```
K(u) = (1 / √(2π)) · e^(-u²/2)
```

The Gaussian kernel produces the smoothest estimates and is the default choice.

**Properties:**
- Infinitely differentiable (very smooth)
- Has infinite support (non-zero for all u, though very small for |u| > 3)
- Computationally straightforward

```go
// Using Gaussian kernel (default)
k, _ := kde.New(data)

// Explicitly set
k, _ := kde.New(data, kde.WithKernelName("gaussian"))
```

### Epanechnikov Kernel

```
K(u) = 0.75 · (1 - u²)  for |u| ≤ 1
K(u) = 0                  for |u| > 1
```

The Epanechnikov kernel is optimal in terms of MISE efficiency (asymptotically optimal).

**Properties:**
- Compactly supported (zero outside [-1, 1])
- Slightly faster computation
- Optimal efficiency (95% of Gaussian at infinite sample size)

```go
k, _ := kde.New(data, kde.WithKernelName("epanechnikov"))
```

### Uniform (Box) Kernel

```
K(u) = 0.5  for |u| ≤ 1
K(u) = 0     for |u| > 1
```

The Uniform kernel gives a "moving histogram" effect.

**Properties:**
- Compactly supported
- Simplest kernel
- Rarely used in practice (less efficient)

```go
k, _ := kde.New(data, kde.WithKernelName("uniform"))
```

### Kernel Comparison

```
Gaussian:       ___________
              /           \
             /             \
            /               \
           /                 \
          /                   \
----------------------------------------

Epanechnikov:   _________
               /         \
              /           \
             /             \
----------------------------------------

Uniform:       ________
              |        |
              |        |
----------------------------------------
```

## Bandwidth Selection

### Silverman's Rule of Thumb

The most common method. Assumes the underlying distribution is normal.

```
h = 0.9 · min(σ, IQR/1.34) · n^(-1/5)
```

Where:
- σ is the standard deviation
- IQR is the interquartile range (Q3 - Q1)
- n is the sample size

This rule is robust to non-normality because it uses min(σ, IQR/1.34), preferring the IQR when data has heavy tails.

```go
// Default: uses Silverman's rule
k, _ := kde.New(data)
```

### Scott's Rule of Thumb

Similar to Silverman's but uses a constant of 1.059:

```
h = 1.059 · σ · n^(-1/5)
```

Scott's rule is slightly more conservative (larger bandwidth) and assumes normality more strictly.

```go
k, _ := kde.New(data, kde.WithBandwidthMethod(kde.Scott))
```

### Manual Bandwidth

For fine-tuning or when you have domain knowledge:

```go
// Smaller than computed: more detail
k, _ := kde.New(data, kde.WithBandwidth(0.5))

// Larger than computed: smoother
k, _ := kde.New(data, kde.WithBandwidth(2.0))
```

### Bandwidth Rules Comparison

| Sample Size | Silverman | Scott | Recommendation |
|-------------|-----------|-------|----------------|
| 50 | ~0.8σ | ~0.9σ | Silverman more robust |
| 200 | ~0.5σ | ~0.6σ | Similar |
| 1000 | ~0.3σ | ~0.35σ | Similar |
| Any | Uses min(σ, IQR) | Uses σ only | Silverman for non-normal data |

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

## Examples

### Comparing Kernels

See how different kernels affect the density estimate:

```go
data := []float64{1.0, 2.0, 2.5, 3.0, 3.5, 4.0, 5.0, 5.5, 6.0, 7.0}

kGaussian, _ := kde.New(data, kde.WithKernelName("gaussian"))
kEpanechnikov, _ := kde.New(data, kde.WithKernelName("epanechnikov"))
kUniform, _ := kde.New(data, kde.WithKernelName("uniform"))

x, yGaussian := kGaussian.EvaluateGrid(0, 8, 100)
_, yEpanechnikov := kEpanechnikov.EvaluateGrid(0, 8, 100)
_, yUniform := kUniform.EvaluateGrid(0, 8, 100)
```

### Bandwidth Effect

Compare different bandwidth settings:

```go
data := []float64{/* your data */}

kNarrow, _ := kde.New(data, kde.WithBandwidth(0.1))
kDefault, _ := kde.New(data) // Silverman's rule
kWide, _ := kde.New(data, kde.WithBandwidth(3.0))
```

### Probability and CDF

Compute probabilities and cumulative distribution values:

```go
k, _ := kde.New(data)

// Probability of falling within [2, 5]
prob := k.Probability(2.0, 5.0, 1000)

// CDF value at x=3
cdf := k.CDF(3.0, 1000)

// Check that total probability is ~1
totalProb := k.Probability(-math.MaxFloat64, math.MaxFloat64, 10000)
```

### Multi-Modal Data

KDE naturally captures multi-modal distributions:

```go
// Bimodal data (two peaks)
data := []float64{
    1.0, 1.2, 1.4, 1.1, 1.3,  // cluster 1
    5.0, 5.2, 5.4, 5.1, 5.3,  // cluster 2
}

k, _ := kde.New(data)
x, densities := k.EvaluateGrid(0, 7, 200)
```

### Weights (Advanced)

When data points have different importance:

```go
data := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
weights := []float64{1.0, 2.0, 1.0, 2.0, 1.0}

k, _ := kde.New(data, kde.WithWeights(weights))

densities := k.Evaluate([]float64{2.5, 3.5})
```

## License

MIT
