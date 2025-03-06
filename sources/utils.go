package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const RedColor = "\033[31m"
const YellowColor = "\033[33m"
const ResetColor = "\033[0m"

func printDebugCondition(text string, show bool) {
	if show {
		printDebug(text)
	}
}

func printDebug(text string) {
	now := time.Now()
	millis := fmt.Sprintf("%d", now.UnixMilli())
	fmt.Println(YellowColor + millis + "| DEBUG: " + text + ResetColor)
}

func printError(text string) {
	now := time.Now()
	millis := fmt.Sprintf("%d", now.UnixMilli())
	fmt.Println(RedColor + millis + "| ERROR: " + text + ResetColor)
}

func showMenu() {
	fmt.Println("\n===== Menu ===== ")
	fmt.Println("1. Show connected marmot")
	fmt.Println("2. Send ping to clients")
	fmt.Println("3. Close connections")
	fmt.Println("4. Execute calculations")
	fmt.Println("5. Update clients software")
	fmt.Println("6. Exit (will let clients trying to reconnect to server)")
	fmt.Print("Choose an option:\n")
}

func handleMenu(marmots Marmots) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		showMenu()
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			marmots.ShowConnected()
		case "2":
			marmots.Pings()
		case "3":
			marmots.CloseConnections()
		case "4":
			handleCalculationMenu(marmots)
		case "5":
			handleClientUpdateMenu(marmots)
		case "6":
			return
		default:
			printError("Invalid option, please try again.")
		}
	}
}

func showClientUpdateMenu() {
	// TODO: add env variable to store the latest client generate
	fmt.Println("\n===== Update Client Menu ===== ")
	fmt.Println("It will send to clients, the latest version of the client software")
	fmt.Printf("The current is: %d\n", ClientVersion)
	fmt.Println("1. YES")
	fmt.Println("2. NO (return)")
}

func handleClientUpdateMenu(marmots Marmots) {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		showClientUpdateMenu()
		scanner.Scan()
		choice := strings.ToLower(strings.TrimSpace(scanner.Text()))

		switch choice {
		case "1":
			marmots.SendUpdateFile()
		case "y":
			marmots.SendUpdateFile()
		case "2":
			return
		case "n":
			return
		default:
			printError("Invalid option, please try again.")
		}
	}
}

func showCalculationMenu() {
	fmt.Println("\n===== Calculation Menu ===== ")
	fmt.Println("1. Counting letter")
	fmt.Println("2. Calculate if a number is prime")
	fmt.Println("3. Calculate Pi estimation")
	fmt.Println("4. Simulate free fall")
	fmt.Println("5. Back")
	fmt.Print("Choose an option:\n")

}

func handleCalculationMenu(marmots Marmots) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		showCalculationMenu()
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			handleCountingLetterMenu(marmots)
		case "2":
			handlePrimeNumberCalculationMenu(marmots)
		case "3":
			handlePiEstimationMenu(marmots)
		case "4":
			handleSimulateFreeFallMenu(marmots)
		case "5":
			return
		default:
			printError("Invalid option, please try again.")
		}
	}

}

// Generic function to generate a scatter plot
func GeneratePlot(data []plotter.XY, title, xLabel, yLabel, filename string) {
	printDebug("Start generating plot")
	start := time.Now()
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel

	// Convert input data to plotter format
	points := make(plotter.XYs, len(data))
	for i, d := range data {
		points[i].X = d.X
		points[i].Y = d.Y
	}

	// Create scatter plot
	scatter, err := plotter.NewScatter(points)
	if err != nil {
		printError(fmt.Sprintf("creation of scatter plot: %s", err))
		return
	}
	scatter.GlyphStyle.Radius = vg.Points(3)

	// Add data to plot and save
	p.Add(scatter)
	if err := p.Save(8*vg.Inch, 5*vg.Inch, filename); err != nil {
		printError(fmt.Sprintf("saving the plot: %s", err))
		return
	}
	elapsed := time.Since(start)
	printDebug(fmt.Sprintf("Graphique generated: '%s' in %v", filename, elapsed))
}

// sort data by time, to avoid weird plot
func sortDataByTime(data []plotter.XY) {
	sort.Slice(data, func(i, j int) bool {
		return data[i].X < data[j].X
	})
}

// save results for gnuplot format
func saveResultsToFile[T any](results []T, filename string, headers []string, format func(T) string) error {
	printDebug("Start saving results to file")
	start := time.Now()
	file, err := os.Create(filename)
	if err != nil {
		printError(fmt.Sprintf("creating the file: %s", err))
		return fmt.Errorf("error during creation file : %v", err)
	}
	defer file.Close()

	// add headers
	if len(headers) > 0 {
		for _, header := range headers {
			fmt.Fprintf(file, "# %s\n", header)
		}
	}

	// write data
	for _, result := range results {
		// it uses the format function given
		dataLine := format(result)
		fmt.Fprintln(file, dataLine)
	}

	elapsed := time.Since(start)
	printDebug(fmt.Sprintf("File created: '%s' in %v", filename, elapsed))
	return nil
}
