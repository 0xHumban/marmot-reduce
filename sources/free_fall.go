package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
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

const FreeFallDefaultTick = 0.001
const FreeFallDefaultTotalSteps = 10000
const FreeFallDefaultcy0 = 100000
const FreeFallDefaultg = 9.81

// / represents simulation parameters
type FreeFall struct {
	Y0         float64
	V0         float64
	G          float64
	TotalSteps float64
	Tick       float64
	Start      float64 // used for client
	End        float64 // used for client
	Results    []FreeFallResult
}

type FreeFallResult struct {
	Time float64 `json:"time"`
	Y    float64 `json:"y"`
}

func formatFreeFallResult(r FreeFallResult) string {
	return fmt.Sprintf("%.6f\t%.6f", r.Time, r.Y)
}

func (m FreeFall) encode() ([]byte, error) {
	return encode(m)
}

func decodeFreeFall(data []byte) (*FreeFall, error) {
	return decode[FreeFall](data)
}

// create a new Free fall struct
func NewFreeFall(y0, v0, g, totalSteps, tick, start, end float64) *FreeFall {
	return &FreeFall{y0, v0, g, totalSteps, tick, start, end, make([]FreeFallResult, 0)}
}

// create new freefall configuration for earth
func NewFreeFallEarth(y0, totalSteps, tick float64) *FreeFall {
	return NewFreeFall(y0, 0, 9.81, totalSteps, 0, totalSteps, tick)
}

func NewDefaultFreeFall() *FreeFall {
	return NewFreeFall(FreeFallDefaultcy0, 0, FreeFallDefaultg, FreeFallDefaultTotalSteps, FreeFallDefaultTick, 0, FreeFallDefaultTotalSteps)
}

// compute a position from his time frame
func (f FreeFall) ComputeFreeFallPosition(t float64) float64 {
	return f.Y0 + f.V0*t - 0.5*f.G*math.Pow(t, 2)
}

// compute positions in range of the given time
func (f *FreeFall) ComputeFreeFallPositionRange(start, end float64) {
	for t := start; t < end; t += f.Tick {
		y := f.ComputeFreeFallPosition(t)
		f.Results = append(f.Results, FreeFallResult{Time: t, Y: y})
	}
}

// compute all positions of free fall simulation
func (f *FreeFall) ComputeFreeFallAllPositionsStartToEnd() {
	printDebug("Start generating free fall position")
	start := time.Now()
	f.ComputeFreeFallPositionRange(f.Start, f.End)
	elapsed := time.Since(start)
	printDebug(fmt.Sprintf("Time to generate free fall positions: %v", elapsed))
}

// generate a plot from
func (f FreeFall) generatePlot() {
	// Convert results to plotter.XY format
	var data []plotter.XY
	for _, r := range f.Results {
		data = append(data, plotter.XY{X: r.Time, Y: r.Y})
	}

	sortDataByTime(data)
	GeneratePlot(data, "Free fall simulation", "Time (s)", "Height (m)", "free_fall_simulation.png")
}

func (f FreeFall) generateResultsFile() {
	headers := []string{"Time (s)", "Height (m)"}
	err := saveResultsToFile(f.Results, "freefall.dat", headers, formatFreeFallResult)
	if err != nil {
		printError(fmt.Sprintf("generating Results File: %s", err))
	}
}

func (f *FreeFall) generateFreeFallData(generatePlot, generateResultsFile bool) {
	f.ComputeFreeFallAllPositionsStartToEnd()
	if generatePlot {
		f.generatePlot()
	}
	if generateResultsFile {
		f.generateResultsFile()
	}
}

func (ff *FreeFall) PrintProperties() {
	fmt.Printf("Initial Height (y0): %.2f m\n", ff.Y0)
	fmt.Printf("Initial Velocity (v0): %.2f m/s\n", ff.V0)
	fmt.Printf("Gravity (g): %.2f m/s²\n", ff.G)
	fmt.Printf("Total Steps: %.2f\n", ff.TotalSteps)
	fmt.Printf("Tick (time step): %.5f\n", ff.Tick)
}

func handleSimulateFreeFallMenu(marmots Marmots) {
	scanner := bufio.NewScanner(os.Stdin)
	ff := NewDefaultFreeFall()

	for {
		fmt.Println("======= FreeFall Simulation ======= ")
		fmt.Println("Current configuration:")
		ff.PrintProperties()
		fmt.Println("1. Set Initial Height (y0)")
		fmt.Println("2. Set Initial Velocity (v0)")
		fmt.Println("3. Set Total Steps")
		fmt.Println("4. Set Time Step (Tick)")
		fmt.Println("5. Calculate FreeFall")
		fmt.Println("6. Print Current Properties")
		fmt.Println("7. Exit")
		fmt.Print("Enter your choice: ")

		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			fmt.Print("Enter the initial height (y0) in meters: ")
			scanner.Scan()
			input := scanner.Text()
			y0, err := strconv.ParseFloat(input, 64)
			if err != nil {
				fmt.Println("Invalid input. Please try again.")
			} else {
				ff.Y0 = y0
			}

		case "2":
			fmt.Print("Enter the initial velocity (v0) in meters per second: ")
			scanner.Scan()
			input := scanner.Text()
			v0, err := strconv.ParseFloat(input, 64)
			if err != nil {
				fmt.Println("Invalid input. Please try again.")
			} else {
				ff.V0 = v0
			}

		case "3":
			fmt.Print("Enter the total number of steps: ")
			scanner.Scan()
			input := scanner.Text()
			steps, err := strconv.ParseFloat(input, 64)
			if err != nil {
				fmt.Println("Invalid input. Please try again.")
			} else {
				ff.TotalSteps = steps
			}

		case "4":
			fmt.Print("Enter the time step (Tick): ")
			scanner.Scan()
			input := scanner.Text()
			tick, err := strconv.ParseFloat(input, 64)
			if err != nil {
				fmt.Println("Invalid input. Please try again.")
			} else {
				ff.Tick = tick
			}

		case "5":

			// Ask user if they want to generate a plot and/or save to file
			fmt.Println("\nDo you want to generate a plot and/or save the results?")
			fmt.Println("1. Generate plot")
			fmt.Println("2. Save results to file")
			fmt.Println("3. Both")
			fmt.Println("4. Neither")
			fmt.Print("Enter your choice: ")

			scanner.Scan()
			subChoice := strings.TrimSpace(scanner.Text())
			generateplot := false
			generateResultsFile := false
			switch subChoice {
			case "1":
				generateplot = true
			case "2":
				generateResultsFile = true
			case "3":
				generateplot = true
				generateResultsFile = true
			case "4":
				fmt.Println("No additional actions taken.")
			default:
				fmt.Println("Invalid choice. No actions taken.")
			}
			// starts repartition + calculations + agregation
			startTime := time.Now()
			ff_res := marmots.FreeFallCalculation(*ff)
			elapsed := time.Since(startTime)
			printDebug(fmt.Sprintf("Time to perform free fall calculation %v", elapsed))

			if ff_res == nil {
				printError("during free fall calculation (freeFall struct is nil)")
				return
			}

			if generateplot {
				ff_res.generatePlot()
			}
			if generateResultsFile {
				ff_res.generateResultsFile()
			}

		case "6":
			ff.PrintProperties()

		case "7":
			fmt.Println("Exiting program.")
			return

		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
