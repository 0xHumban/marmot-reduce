package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const LocalServerIP = "192.168.1.25:8080"
const ServerIP = "127.0.0.1:8080"
const RetryDelais = 5

// handle connection client side
// client waiting for server instructions
// 'exit': connection closed
// '1': count 'e' in response
// returns if the connection has been asked by server
func handleConnectionClientSide(conn net.Conn) bool {
	defer conn.Close()
	response := ""
	for response != "exit" {
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("ERROR reading server response", err)
			return false
		}
		response = response[:len(response)-1]

		// fmt.Printf("Server response: '%s'\n", response)
		if response == "exit" {
			printDebug("EXIT request received")
			return true
		}

		treatServerResponse(conn, response)

	}
	return true
}

// treats the server response
// choose whats the next step, which function the client have to execute
func treatServerResponse(conn net.Conn, response string) {
	// Ping request
	if response[0] == '0' {
		printDebug("Ping pong request received")
		message := fmt.Sprintf("'Pong' from @%s\n", conn.LocalAddr().String())
		_, _ = conn.Write([]byte(message))
		printDebug("Ping pong response sent")

	} else if
	// count 'e' in response
	response[0] == '1' {
		printDebug("Start counting letter occurrences\n")
		letterToCount := response[1]
		calculus := fmt.Sprintf("%d", countLetterOccurrence(response[1:], rune(letterToCount)))
		printDebug(fmt.Sprintf("End couting letter occurrences -- Result for '%c': %s\n", rune(letterToCount), calculus))
		_, _ = conn.Write([]byte(fmt.Sprintf("%s\n", calculus)))

	} else if
	// calculate if a number is prime in a given range
	response[0] == '2' {
		printDebug("Start prime number calculation\n")
		parts := strings.Split(response[1:], "@")

		if len(parts) != 3 {
			printError("Invalid format")
			_, _ = conn.Write([]byte(fmt.Sprintf("%s\n", "Invalid format")))
			return
		}

		potentialPrime, err1 := strconv.Atoi(parts[0])
		start, err2 := strconv.Atoi(parts[1])
		end, err3 := strconv.Atoi(parts[2])

		if err1 != nil || err2 != nil || err3 != nil {
			printError("during conversion")
			_, _ = conn.Write([]byte(fmt.Sprintf("%s\n", "Conversion error")))
			return
		}
		calculus := fmt.Sprintf("%d", calculatePrimeNumber(potentialPrime, start, end))
		printDebug(fmt.Sprintf("End prime number calculation -- Result: %s\n", calculus))
		_, _ = conn.Write([]byte(fmt.Sprintf("%s\n", calculus)))

	}
}

// Take a word and returns occurrences of the given letter
func countLetterOccurrence(word string, letterToCount rune) int {
	res := 0
	for _, letter := range word {
		if letter == letterToCount {
			res++
		}
	}
	return res
}

// calculate if a number is prime in a given range
// returns '-1' if there is no number product of the potentialPrime, else returns the first one found
//
//	    IN: potentialPrime: The number to calculate if is prime
//		IN: start:          Range start
//		IN: end:            Range end
func calculatePrimeNumber(potentialPrime, start, end int) int {
	res := -1
	i := start
	for i < end && res == -1 {
		if potentialPrime%i == 0 {
			res = i
		}
		i++
	}

	return res
}

func connectToServer(ip string) {
	connectionClosedProperly := false

	for !connectionClosedProperly {
		conn, err := net.Dial("tcp", ip)
		if err != nil {
			fmt.Println("ERROR connecting to server", err)
		} else {
			// DEBUG
			printDebug("Local address: " + conn.LocalAddr().String())
			printDebug("Remote address: " + conn.RemoteAddr().String())

			connectionClosedProperly = handleConnectionClientSide(conn)
		}
		if !connectionClosedProperly {
			time.Sleep(RetryDelais * time.Second)
		}
	}
}
