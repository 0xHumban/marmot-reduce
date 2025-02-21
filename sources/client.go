package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

const LocalServerIP = "192.168.223.130:8080"
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

	} else if
	// count 'e' in response
	response[0] == '1' {
		printDebug("Start simulating calculus\n")
		letterToCount := response[1]
		calculus := fmt.Sprintf("%d", simulateClientCalculus(response[1:], rune(letterToCount)))
		printDebug(fmt.Sprintf("End simulating calculus -- Result'e': %s\n", calculus))
		_, _ = conn.Write([]byte(fmt.Sprintf("%s\n", calculus)))

	}
}

// simulate client calculus
// Take a word and returns occurrences of 'a'
func simulateClientCalculus(word string, letterToCount rune) int {
	res := 0
	for _, letter := range word {
		if letter == letterToCount {
			res++
		}
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
