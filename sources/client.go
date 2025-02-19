package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const ServerIP = "192.168.223.130:8080"
const LocalServerIP = "127.0.0.1:8080"

// handle connection client side
// client waiting for server instructions
// 'exit': connection closed
// '1': count 'e' in response
func handleConnectionClientSide(conn net.Conn) {
	defer conn.Close()
	response := ""
	for response != "exit" {
		response, err := bufio.NewReader(conn).ReadString('\n')
		response = response[:len(response)-1]
		if err != nil {
			fmt.Println("ERROR reading server response", err)
			return
		}

		// fmt.Printf("Server response: '%s'\n", response)
		if response == "exit" {
			fmt.Println("EXIT ASKED")
			return
		}

		// count 'e' in response
		if response[0] == '1' {
			fmt.Printf("Start simulating calculus\n")
			letterToCount := response[1]
			calculus := fmt.Sprintf("%d", simulateClientCalculus(response[1:], rune(letterToCount)))
			fmt.Printf("End simulating calculus -- Result'e': %s\n", calculus)
			conn.Write([]byte(fmt.Sprintf("%s\n", calculus)))

		}

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
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		fmt.Println("ERROR connecting to server", err)
		os.Exit(1)
	}
	// DEBUG
	fmt.Println("Local address: ", conn.LocalAddr())
	fmt.Println("Remote address: ", conn.RemoteAddr())

	message := fmt.Sprintf("The client address is:@%s\n", conn.LocalAddr().String())

	conn.Write([]byte(message))

	handleConnectionClientSide(conn)

}
