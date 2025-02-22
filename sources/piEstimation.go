package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
				marmots.PiCalculation(number)

				return
			}
		}

	}
}
