package chart

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import "time"

// SecondsPerXYZ
const (
	SecondsPerHour = 60 * 60
	SecondsPerDay  = 60 * 60 * 24
)

// TimeMillis returns a duration as a float millis.
func TimeMillis(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}

// DiffHours returns the difference in hours between two times.
func DiffHours(t1, t2 time.Time) (hours int) {
	t1n := t1.Unix()
	t2n := t2.Unix()
	var diff int64
	if t1n > t2n {
		diff = t1n - t2n
	} else {
		diff = t2n - t1n
	}
	return int(diff / (SecondsPerHour))
}

// TimeMin returns the minimum and maximum times in a given range.
func TimeMin(times ...time.Time) (min time.Time) {
	if len(times) == 0 {
		return
	}
	min = times[0]
	for index := 1; index < len(times); index++ {
		if times[index].Before(min) {
			min = times[index]
		}

	}
	return
}

// TimeMax returns the minimum and maximum times in a given range.
func TimeMax(times ...time.Time) (max time.Time) {
	if len(times) == 0 {
		return
	}
	max = times[0]

	for index := 1; index < len(times); index++ {
		if times[index].After(max) {
			max = times[index]
		}
	}
	return
}

// TimeMinMax returns the minimum and maximum times in a given range.
func TimeMinMax(times ...time.Time) (min, max time.Time) {
	if len(times) == 0 {
		return
	}
	min = times[0]
	max = times[0]

	for index := 1; index < len(times); index++ {
		if times[index].Before(min) {
			min = times[index]
		}
		if times[index].After(max) {
			max = times[index]
		}
	}
	return
}

// TimeToFloat64 returns a float64 representation of a time.
func TimeToFloat64(t time.Time) float64 {
	return float64(t.UnixNano())
}

// TimeFromFloat64 returns a time from a float64.
func TimeFromFloat64(tf float64) time.Time {
	return time.Unix(0, int64(tf))
}

// TimeDescending sorts a given list of times ascending, or min to max.
type TimeDescending []time.Time

// Len implements sort.Sorter
func (d TimeDescending) Len() int { return len(d) }

// Swap implements sort.Sorter
func (d TimeDescending) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

// Less implements sort.Sorter
func (d TimeDescending) Less(i, j int) bool { return d[i].After(d[j]) }

// TimeAscending sorts a given list of times ascending, or min to max.
type TimeAscending []time.Time

// Len implements sort.Sorter
func (a TimeAscending) Len() int { return len(a) }

// Swap implements sort.Sorter
func (a TimeAscending) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less implements sort.Sorter
func (a TimeAscending) Less(i, j int) bool { return a[i].Before(a[j]) }

// Days generates a seq of timestamps by day, from -days to today.
func Days(days int) []time.Time {
	var values []time.Time
	for day := days; day >= 0; day-- {
		values = append(values, time.Now().AddDate(0, 0, -day))
	}
	return values
}

// Hours returns a sequence of times by the hour for a given number of hours
// after a given start.
func Hours(start time.Time, totalHours int) []time.Time {
	times := make([]time.Time, totalHours)

	last := start
	for i := 0; i < totalHours; i++ {
		times[i] = last
		last = last.Add(time.Hour)
	}

	return times
}

// HoursFilled adds zero values for the data bounded by the start and end of the xdata array.
func HoursFilled(xdata []time.Time, ydata []float64) ([]time.Time, []float64) {
	start, end := TimeMinMax(xdata...)
	totalHours := DiffHours(start, end)

	finalTimes := Hours(start, totalHours+1)
	finalValues := make([]float64, totalHours+1)

	var hoursFromStart int
	for i, xd := range xdata {
		hoursFromStart = DiffHours(start, xd)
		finalValues[hoursFromStart] = ydata[i]
	}

	return finalTimes, finalValues
}
