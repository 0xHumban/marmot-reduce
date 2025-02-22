package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func handlePrimeNumberCalculationMenu(marmots Marmots) {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("======= Prime number Calculation ======= ")
		fmt.Println("It will divide by client number the range from 2 to sqrt(YourNumber), to know if the given number is prime")
		fmt.Println("This algorithm used the naive version with a loop in the range")
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
				marmots.PrimeNumberCalculation(number)

				return
			}
		}

	}
}
