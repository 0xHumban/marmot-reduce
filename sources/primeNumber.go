package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Create a range from 2 to sqrt(potentialNumber)
// Divide this range by clients number
// send range and wait for result
func (ms Marmots) PrimeNumberCalculation(potentialPrime int) {
	printDebug("Start prime number calculation")
	// Send ping to check if clients always connected
	ms.Pings()
	clientsNumber := ms.clientsLen()
	if clientsNumber == 0 {
		printError("No client connected, retry after connecting clients")
		return
	}
	start := 2
	i := 1
	sqrtNumber := math.Sqrt(float64(potentialPrime))
	subRangeLength := int(sqrtNumber / float64(clientsNumber))
	// for little number: set minimal range to 3
	if subRangeLength < 4 {
		subRangeLength = 3
	}
	for _, m := range ms {
		if m != nil {
			m.data = createMessage("3", String, []byte(fmt.Sprintf("%d@%d@%d", potentialPrime, start, (subRangeLength*i))))
			start += subRangeLength
		}
		i++
	}
	ms.performAction((*Marmot).PrimeNumber)
	res := false
	for _, m := range ms {
		if m != nil && <-m.end {
			if string(m.response.Data) != "-1" {
				fmt.Printf("The number '%d' is not prime, first factor found: '%s' on %s\n", potentialPrime, m.response, m.conn.RemoteAddr())
				res = true
			}
		}
	}
	if !res {
		fmt.Printf("The number '%d' is prime\n", potentialPrime)
	}
	printDebug("End prime number calculation")
}

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
