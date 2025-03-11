package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func handlePiEstimationMenu(marmots Marmots) {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("======= PI Estimation ======= ")
		fmt.Println("It will divide by client number the range, from 0 to the number of points given")
		fmt.Println("This algorithm is the Monte Carlo Algorithm")
		fmt.Println("Enter the number you want ")
		fmt.Println("(Enter '-1' to leave)")

		scanner.Scan()

		choice := strings.TrimSpace(scanner.Text())
		if choice == "-1" {
			return
		} else {
			number, err := strconv.Atoi(choice)
			if err != nil {
				printError("Invalid option, please try again.")
			} else {
				startTime := time.Now()
				piEstimate := marmots.PiCalculation(number)
				endTime := time.Now()
				// Calculate the duration
				duration := endTime.Sub(startTime)

				fmt.Printf("Estimation of Pi: %.20f\n", piEstimate)
				fmt.Printf("Time taken: %v\n", duration)

				return
			}
		}

	}
}

func (m *Marmot) PiCalculation() {
	m.SendAndReceiveData("Pi calculation", true)
}

// Create a range from 0 to given number
// Divide this range by clients number
// send range and wait for result and calculate a pi estimation
func (ms Marmots) PiCalculation(numSamples int) float64 {
	printDebug("Start PI calculation")
	// Send ping to check if clients always connected
	ms.Pings()
	clientsNumber := ms.clientsLen()
	if clientsNumber == 0 {
		printError("No client connected, retry after connecting clients")
		return -1
	}
	samplesPerWorker := numSamples / clientsNumber
	for _, m := range ms {
		if m != nil {
			m.data = createMessage("4", String, []byte(fmt.Sprintf("%d", samplesPerWorker)))
		}
	}
	ms.performAction((*Marmot).PiCalculation)
	insideTotal := 0
	for _, m := range ms {
		if m != nil && <-m.end {
			numberinside, err := strconv.Atoi(string(m.response.Data))
			if err != nil {
				printError("during response conversion to number")
			} else {
				insideTotal += numberinside
			}
		}
	}
	pi := float64(insideTotal) / float64(numSamples) * 4
	printDebug("End PI calculation")
	printDebug(fmt.Sprintf("The PI estimation is: ~%.20f , with %d clients, a total sample of %d\n", pi, clientsNumber, numSamples))
	return pi
}
