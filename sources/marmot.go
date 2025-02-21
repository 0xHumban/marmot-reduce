package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
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

func NewMarmot(connA net.Conn) *Marmot {
	return &Marmot{
		conn:     connA,
		start:    make(chan bool, 1),
		end:      make(chan bool, 1),
		data:     "",
		response: "",
	}
}

// Represents the clients list
type Marmots []*Marmot

type MarmotI interface {
	// send a ping to the client to see if the connection is always up
	Ping() (bool, error)
}

func (ms Marmots) performAction(fctToExecute func(*Marmot)) {
	// wait group to await goroutine
	var wg sync.WaitGroup

	for _, m := range ms {
		if m != nil {
			wg.Add(1)
			go func(m *Marmot) {
				defer wg.Done()
				fctToExecute(m)
			}(m)
		}
	}
	// All clients are connected we can start calculations
	for _, marmot := range ms {
		if marmot != nil {
			marmot.start <- true
		}
	}

	// wait for all goroutines to end
	wg.Wait()

}

// send ping to all current clients
func (ms Marmots) Pings() {

	ms.performAction((*Marmot).Ping)
	// if client goroutine has 'end' = false
	// it means there is an error and we remove it from the list
	for i, m := range ms {
		if m != nil && !<-m.end {
			ms[i] = nil
			printDebug("@" + m.conn.RemoteAddr().String() + " has been removed of the clients list")
		}
	}
}

// returns the number of clients connected
func (ms Marmots) clientsLen() int {
	res := 0
	for i := range len(ms) {
		if ms[i] != nil {
			res++
		}
	}
	return res
}

// sends batch of letters, and asked to clients to count occurence of a letter
func (ms Marmots) CountingLetters(letter rune, batchSize int) {
	// Send ping to check if clients always connected
	ms.Pings()

	clientsNumber := ms.clientsLen()
	dataset := generateRandomArray(clientsNumber, batchSize)
	i := 0
	for _, m := range ms {
		if m != nil {
			m.data = fmt.Sprintf("%d%c%s\n", 1, letter, dataset[i])
			i++
		}
	}
	ms.performAction((*Marmot).CountLetter)

}

// show current clients connected
func (ms Marmots) ShowConnected() {
	ms.Pings()
	fmt.Println("\nCurrent clients connected:")
	for i, m := range ms {
		if m != nil {
			fmt.Printf("%d. @%s\n", (i + 1), m.conn.RemoteAddr())
		}
	}
}

// close all current conncections with clients
func (ms Marmots) CloseConnections() {
	for i, m := range ms {
		if m != nil {
			m.Close()
			ms[i] = nil
			printDebug("@" + m.conn.RemoteAddr().String() + " has been closed and removed of the client list")
		}
	}

}

// close the connection with the client
// it sends 'exit' for properly closed
func (m *Marmot) Close() {
	defer m.conn.Close()
	m.data = "exit\n"
	if !m.writeData(true) {
		printDebug("error sending 'exit'")
	}
}

func (m *Marmot) CountLetter() {
	// wait for start / timeout
	<-m.start
	// sending data
	res := m.writeData(false)
	if res {
		res = m.readResponse()
	} else {
		printDebug("error sending 'Ping/Pong'")
		m.end <- false
		return
	}
	if !res {
		printDebug("error receiving 'Ping/Pong'")
		m.end <- false
		return
	}
	m.end <- true

}

func (m *Marmot) Ping() {
	// wait for start / timeout
	<-m.start
	// send 'ping' to client
	m.data = "0Ping\n"
	res := m.writeData(false)
	if res {
		res = m.readResponse()
	} else {
		printDebug("error sending 'Ping/Pong'")
		m.end <- false
		return
	}

	if !res {
		printDebug("error receiving 'Ping/Pong'")
		m.end <- false
		return
	}

	m.end <- true

}

// read client response
// print in DEBUG mode and store the message in response
func (m *Marmot) readResponse() bool {

	message, err := bufio.NewReader(m.conn).ReadString('\n')
	if err != nil {
		printError(fmt.Sprintf("Reading client (@"+m.conn.RemoteAddr().String()+") message '%s'", err))
		return false
	}
	printDebug("Message received from @" + m.conn.RemoteAddr().String() + ": " + message)
	m.response = message
	return true
}

// write the Data store in Marmot to the client
// print in DEBUG mode
func (m *Marmot) writeData(show bool) bool {

	// message, err := bufio.NewReader(m.conn).ReadString('\n')
	_, err := m.conn.Write([]byte(m.data))
	if err != nil {
		printError(fmt.Sprintf("Sending message to client '(@"+m.conn.RemoteAddr().String()+" %s'", err))
		return false
	}
	printDebugCondition("Message sent: '"+m.data+"'", show)
	return true
}
