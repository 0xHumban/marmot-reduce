package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

const CoutingLetterBatchSize = 10000000

func (m *Marmot) CountLetter() {
	m.SendAndReceiveData("Count letter", false)
}

// sends batch of letters, and asked to clients to count occurence of a letter
func (ms Marmots) CountingLetters(letter rune, batchSize int) {
	printDebug("Start counting letters")
	// Send ping to check if clients always connected
	ms.Pings()

	clientsNumber := ms.clientsLen()
	if clientsNumber == 0 {
		printError("No client connected, retry after connecting clients")
		return
	}
	dataset := generateRandomArray(clientsNumber, batchSize)
	i := 0
	for _, m := range ms {
		if m != nil {
			m.data = createMessage("2", String, []byte(fmt.Sprintf("%c%s", letter, dataset[i])))
			i++
		}
	}
	ms.performAction((*Marmot).CountLetter)
	ms.WaitEnd()
	printDebug("End counting letters")

}

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
