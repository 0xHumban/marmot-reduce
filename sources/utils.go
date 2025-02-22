package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
	fmt.Println(YellowColor + "DEBUG: " + text + ResetColor)
}

func printError(text string) {
	fmt.Println(RedColor + "ERROR: " + text + ResetColor)
}

func showMenu() {
	fmt.Println("\n===== Menu ===== ")
	fmt.Println("1. Show connected marmot")
	fmt.Println("2. Send ping to clients")
	fmt.Println("3. Close connections")
	fmt.Println("4. Execute calculations")
	fmt.Println("5. Exit (will let clients trying to reconnect to server)")
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
	fmt.Println("4. Back")
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
			return
		default:
			printError("Invalid option, please try again.")
		}
	}

}
