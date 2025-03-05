package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
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
