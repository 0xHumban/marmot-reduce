package main

import (
	"fmt"
	"math"
	"time"

	"gonum.org/v1/plot/plotter"
)

/// Package used to simulate free fall simulation
/*
	y = y0 + v0.t - (1/2)g.t²

	y0 = initial height (m)
	v0 = initial velocity (m/s)
	g = 9.81 m / s² (gravity)
	t = time (s)
*/

// / represents simulation parameters
type FreeFall struct {
	y0         float64
	v0         float64
	g          float64
	totalSteps float64
	tick       float64
	results    []FreeFallResult
}

type FreeFallResult struct {
	Time float64 `json:"time"`
	Y    float64 `json:"y"`
}

func formatFreeFallResult(r FreeFallResult) string {
	return fmt.Sprintf("%.6f\t%.6f", r.Time, r.Y)
}

// create a new Free fall struct
func NewFreeFall(y0, v0, g, totalSteps, tick float64) *FreeFall {
	return &FreeFall{y0, v0, g, totalSteps, tick, make([]FreeFallResult, 0)}
}

// create new freefall configuration for earth
func NewFreeFallEarth(y0, totalSteps, tick float64) *FreeFall {
	return NewFreeFall(y0, 0, 9.81, totalSteps, tick)
}

// compute a position from his time frame
func (f FreeFall) ComputeFreeFallPosition(t float64) float64 {
	return f.y0 + f.v0*t - 0.5*f.g*math.Pow(t, 2)
}

// compute positions in range of the given time
func (f *FreeFall) ComputeFreeFallPositionRange(start, end float64) {
	for t := start; t < end; t += f.tick {
		y := f.ComputeFreeFallPosition(t)
		f.results = append(f.results, FreeFallResult{Time: t, Y: y})
	}
}

// compute all positions of free fall simulation
func (f *FreeFall) ComputeFreeFallAllPositions() {
	printDebug("Start generating free fall data")
	start := time.Now()
	f.ComputeFreeFallPositionRange(0, f.totalSteps)
	elapsed := time.Since(start)
	printDebug(fmt.Sprintf("Time to generate free fall points: %v", elapsed))
}

// generate a plot from
func (f FreeFall) generatePlot() {
	// Convert results to plotter.XY format
	var data []plotter.XY
	for _, r := range f.results {
		data = append(data, plotter.XY{X: r.Time, Y: r.Y})
	}

	sortDataByTime(data)
	GeneratePlot(data, "Free fall simulation", "Time (s)", "Height (m)", "free_fall_simulation.png")
}

func (f FreeFall) generateResultsFile() {
	headers := []string{"Time (s)", "Height (m)"}
	err := saveResultsToFile(f.results, "freefall.dat", headers, formatFreeFallResult)
	if err != nil {
		printError(fmt.Sprintf("generating Results File: %s", err))
	}
}

func (f FreeFall) generateFreeFallData(generatePlot, generateResultsFile bool) {
	f.ComputeFreeFallAllPositions()
	if generatePlot {
		f.generatePlot()
	}
	if generateResultsFile {
		f.generateResultsFile()
	}
}
