package main

import (
	"fmt"
	"net"
)

const ServerPort = ":8080"
const ClientNumber = 3

// open a port to allow client to connect
// In:
// - port: port to open
// - handleFct: function pointer to handle different connexions
func openConnection(port string, marmots Marmots) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("ERROR during listening:", err)
	}

	defer ln.Close()

	printDebug("Server waiting for connections")
	// DATASET FOR CLIENT
	// marmots := make([]Marmot, clientNumber)
	// dataset := generateRandomArray(ClientNumber, 1000000)
	for i := 0; i < ClientNumber; i++ {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("ERROR accepting connection: ", err)
			continue
		} else {
			printDebug("New client connected: @" + conn.RemoteAddr().String())
		}
		marmots[i] = NewMarmot(conn)
	}

}
