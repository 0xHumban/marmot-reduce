package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

const RedColor = "\033[31m"
const YellowColor = "\033[33m"
const ResetColor = "\033[0m"

func generateRandomArray(arraylength, stringLength int) []string {
	res := make([]string, arraylength)
	for i := range res {
		res[i] = generateRandomString(stringLength)
	}

	return res
}

func generateRandomString(length int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func printDebug(text string) {
	fmt.Println(YellowColor + "\nDEBUG: " + text + ResetColor)
}

func printError(text string) {
	fmt.Println(RedColor + "\nERROR: " + text + ResetColor)
}

func showMenu() {
	fmt.Println("\n===== Menu ===== ")
	fmt.Println("1. Show connected marmot")
	fmt.Println("2. Send ping to clients")
	fmt.Println("3. Close connections")
	fmt.Print("Choose an option:")
}

func handleMenu(marmots Marmots) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		showMenu()
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			fmt.Println("TODO")
		case "2":
			marmots.Pings()
		case "3":
			fmt.Println("TODO")
		default:
			printError("Invalid option, please try again.")
		}
	}
}
