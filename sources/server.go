package main

import (
	"bufio"
	"fmt"
	"net"
)

const ServerPort = ":8080"

// open a port to allow client to connect
// In:
// - port: port to open
// - handleFct: function pointer to handle different connexions
func openConnection(port string, handleFct func(conn net.Conn)) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("ERROR during listening:", err)
	}

	defer ln.Close()

	fmt.Println("Server waiting for connections")
	// DATASET FOR CLIENT
	clientNumber := 3
	marmots := make([]Marmot, clientNumber)
	dataset := generateRandomArray(clientNumber, 100000000)
	for i := 0; i < 3; i++ {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("ERROR accepting connection: ", err)
			continue
		}
		marmots[i] = Marmot{conn, make(chan bool), make(chan bool), dataset[i], ""}
		go marmots[i].handleMarmot()
		// go handleFct(conn)
	}

	// All clients are connected we can start calculations
	for _, marmot := range marmots {
		marmot.start <- true
	}
	for _, marmot := range marmots {
		<-marmot.end
	}
}

// handle the client connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// DEBUG
	fmt.Println("Local address: ", conn.LocalAddr())
	fmt.Println("Remote address: ", conn.RemoteAddr())

	// read client message
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("ERROR reading client message", err)
		return
	}

	fmt.Printf("Message received: %s", message)

	// send response to client
	response := "1eMessage received by server\n"
	_, _ = conn.Write([]byte(response))

	// DELETE ---
	// read client message
	message, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("ERROR reading client message", err)
		return
	}
	fmt.Printf("Message received by client: %s", message)

	// End connection by sending 'exit'
	fmt.Printf("Connection closed\n")
	response = "exit\n"
	_, _ = conn.Write([]byte(response))

}
