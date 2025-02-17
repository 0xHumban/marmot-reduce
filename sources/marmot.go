package main

import (
	"bufio"
	"fmt"
	"net"
)

// Represents a client
// conn: used to communicate with client
// start: waiting to send data to client
// end: used to synchronise client response
// data: dataset to send to client
// response: response of the client after calculations
type Marmot struct {
	conn     net.Conn
	start    chan bool
	end      chan bool
	data     string
	response string
}

// Represents the clients list
type Marmots []*Marmot

type MarmotI interface {
	// send a ping to the client to see if the connection is always up
	Ping() (bool, error)
}

// func (m Marmot) Ping() (bool, error) {

// }

// Wait for 'start'
// send dataset to client
// waiting for response
// Send end signal with channel
func (m *Marmot) handleMarmot() {
	defer m.conn.Close()

	// DEBUG
	fmt.Println("Local address: ", m.conn.LocalAddr())
	fmt.Println("Remote address: ", m.conn.RemoteAddr())
	fmt.Println("MARMOT NOT STARTED YET")
	<-m.start
	fmt.Println("MARMOT STARTED")

	// read client message
	message, err := bufio.NewReader(m.conn).ReadString('\n')
	if err != nil {
		fmt.Println("ERROR reading client message", err)
		return
	}

	fmt.Printf("Message received: %s", message)

	// send response to client
	m.data = "1e" + m.data + "\n"
	// response := "1eMessage received by server\n"
	response := m.data
	_, _ = m.conn.Write([]byte(response))

	// DELETE ---
	// read client message
	message, err = bufio.NewReader(m.conn).ReadString('\n')
	if err != nil {
		fmt.Println("ERROR reading client message", err)
		return
	}
	fmt.Printf("Message received by client: %s", message)

	// End connection by sending 'exit'
	fmt.Printf("Connection closed\n")
	response = "exit\n"
	_, _ = m.conn.Write([]byte(response))

	m.end <- true

}
