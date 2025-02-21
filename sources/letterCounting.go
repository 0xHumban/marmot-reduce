package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

const CoutingLetterBatchSize = 10000000

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

func handleCountingLetterMenu(marmots Marmots) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("======= Counting letter occurrence ======= ")
		fmt.Println("It will send to clients batch of random letters and to count the occurence of the letter asked")
		fmt.Printf("The current batch size is :%d\n", CoutingLetterBatchSize)
		fmt.Println("Enter the letter you want to count occurences")
		fmt.Println("(Enter '-1' to leave)")

		scanner.Scan()

		choice := strings.TrimSpace(scanner.Text())
		if choice == "-1" {
			return
		} else if len(choice) == 1 {
			marmots.CountingLetters(rune(choice[0]), CoutingLetterBatchSize)
			return
		} else {
			printError("Invalid option, please try again.")
		}

	}
}
