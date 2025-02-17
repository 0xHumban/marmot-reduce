package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const ServerIP = "127.0.0.1:8080"

// handle connection client side
// client waiting for server instructions
func handleConnectionClientSide(conn net.Conn) {
	defer conn.Close()
	response := ""
	for response != "exit\n" {
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("ERROR reading server response", err)
			return
		}
		if response == "exit\n" {
			fmt.Println("EXITT ASKED")
			return
		}

		fmt.Printf("Server response: '%s'", response)
		fmt.Printf("Start simulating calculus\n")
		calculus := simulateClientCalculus(response)
		fmt.Printf("End simulating calculus\n'a': %i\n", calculus)
		conn.Write([]byte(fmt.Sprintf("%d\n", calculus)))

	}

}

// simulate client calculus
// Take a word and returns occurrences of 'a'
func simulateClientCalculus(word string) int {
	res := 0
	for _, letter := range word {
		if letter == 'a' {
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

	// response, err := bufio.NewReader(conn).ReadString('\n')
	// if err != nil {
	// 	fmt.Println("ERROR reading server response", err)
	// 	return
	// }

	// fmt.Printf("Server response: %s", response)

	// for {
	// }
}

func main() {
	connectToServer(ServerIP)
}
