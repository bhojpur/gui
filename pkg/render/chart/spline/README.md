# Bhojpur GUI - Spline Library

The `spline` library generates a cubic spline for given points.

## Simple Usage

### Create a cubic spline

```golang
s := spline.NewCubicSpline([]float64{0, 1, 2, 3}, []float64{0, 0.5, 2, 1.5})
```

### Get an interpolated value

```golang
s.At(3.5)
```

### Get an array of interpolated values

```
s.Range(0, 3, 0.25)
```

### Supported boundaries

First derivation boundary: spline.NewClampedCubicSpline

Second derivation boundary: spline.NewNaturalCubicSpline

## Installation

Just
`go get github.com/bhojpur/gui`