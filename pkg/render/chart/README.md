# Bhojpur GUI - Chart Library

The `chart` library is part of rendering system. It is designed to supports `timeseries`
and __continuous__ `line charts`. It is pre-integrated with document renderers.

## Installation

To install `chart` library, run the following:

```bash
> go get -u github.com/bhojpur/gui
```

Most of the components are interchangeable.

## Sinple Usage

Everything starts with the `chart.Chart` object. The bare minimum to draw a chart would be the following:

```golang

import (
    ...
    "bytes"
    ...
    "github.com/bhojpur/gui/pkg/render/chart" //exposes "chart"
)

graph := chart.Chart{
    Series: []chart.Series{
        chart.ContinuousSeries{
            XValues: []float64{1.0, 2.0, 3.0, 4.0},
            YValues: []float64{1.0, 2.0, 3.0, 4.0},
        },
    },
}

buffer := bytes.NewBuffer([]byte{})
err := graph.Render(chart.PNG, buffer)
```

**Explanation of the above example:**

A `chart` can have many `Series`, and a `Series` is a collection of things that need to be drawn according to the X range and the Y range(s).

Here, we have a single series with x range values as float64s, rendered to a PNG.

NOTE: we can pass any type of `io.Writer` into `Render(...)`, meaning that we can render the chart to a file or a response or anything else that implements `io.Writer`.

## Charting API

Everything on the `chart.Chart` object has defaults that can be overriden. Whenever a software
developer sets a property on the chart object, it is to be assumed that value will be used
instead of the default.