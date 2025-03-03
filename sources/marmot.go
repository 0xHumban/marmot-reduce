package main

import (
	"bufio"
	"context"
	"fmt"
	"math"
	"net"
	"strconv"
	"sync"
	"time"
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

func (ms Marmots) Pings() {
	printDebug("Start Pings")
	ms.performAction((*Marmot).Ping)
	// if client goroutine has 'end' = false
	// it means there is an error and we remove it from the list
	for i, m := range ms {
		if m != nil && !<-m.end {
			ms[i] = nil
			printDebug("@" + m.conn.RemoteAddr().String() + " has been removed of the clients list")
		}
	}
	printDebug("End Pings")
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
			m.data = fmt.Sprintf("%d%c%s\n", 1, letter, dataset[i])
			i++
		}
	}
	ms.performAction((*Marmot).CountLetter)
	ms.WaitEnd()
	printDebug("End counting letters")

}

// Create a range from 2 to sqrt(potentialNumber)
// Divide this range by clients number
// send range and wait for result
func (ms Marmots) PrimeNumberCalculation(potentialPrime int) {
	printDebug("Start prime number calculation")
	// Send ping to check if clients always connected
	ms.Pings()
	clientsNumber := ms.clientsLen()
	if clientsNumber == 0 {
		printError("No client connected, retry after connecting clients")
		return
	}
	start := 2
	i := 1
	sqrtNumber := math.Sqrt(float64(potentialPrime))
	subRangeLength := int(sqrtNumber / float64(clientsNumber))
	// for little number: set minimal range to 3
	if subRangeLength < 4 {
		subRangeLength = 3
	}
	for _, m := range ms {
		if m != nil {
			m.data = fmt.Sprintf("%d%d@%d@%d\n", 2, potentialPrime, start, (subRangeLength * i))
			start += subRangeLength
		}
		i++
	}
	ms.performAction((*Marmot).PrimeNumber)
	res := false
	for _, m := range ms {
		if m != nil && <-m.end {
			if m.response != "-1" {
				fmt.Printf("The number '%d' is not prime, first factor found: '%s' on %s\n", potentialPrime, m.response, m.conn.RemoteAddr())
				res = true
			}
		}
	}
	if !res {
		fmt.Printf("The number '%d' is prime\n", potentialPrime)
	}
	printDebug("End prime number calculation")
}

// Create a range from 0 to given number
// Divide this range by clients number
// send range and wait for result and calculate a pi estimation
func (ms Marmots) PiCalculation(numSamples int) float64 {
	printDebug("Start PI calculation")
	// Send ping to check if clients always connected
	ms.Pings()
	clientsNumber := ms.clientsLen()
	if clientsNumber == 0 {
		printError("No client connected, retry after connecting clients")
		return -1
	}
	samplesPerWorker := numSamples / clientsNumber
	for _, m := range ms {
		if m != nil {
			m.data = fmt.Sprintf("%d%d\n", 3, samplesPerWorker)
		}
	}
	ms.performAction((*Marmot).PiCalculation)
	insideTotal := 0
	for _, m := range ms {
		if m != nil && <-m.end {
			numberinside, err := strconv.Atoi(m.response)
			if err != nil {
				printError("during response conversion to number")
			} else {
				insideTotal += numberinside
			}
		}
	}
	pi := float64(insideTotal) / float64(numSamples) * 4
	printDebug("End PI calculation")
	printDebug(fmt.Sprintf("The PI estimation is: ~%.20f , with %d clients, a total sample of %d\n", pi, clientsNumber, numSamples))
	return pi
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

// wait for the end of the current action started
// if end true the client correctly sent the answer
// if end false then an error occurs during the communication (can be timeout or other)
// Use: use it after using `performAction` and you're not reading the `<-m.end` after
// if you miss this, it will bloqued the program
func (ms Marmots) WaitEnd() {
	for _, m := range ms {
		if m != nil {
			<-m.end
		}
	}

}

func (m *Marmot) PrimeNumber() {
	// wait for start / timeout
	<-m.start
	// sending data
	res := m.writeData(true)
	if res {
		res = m.readResponse()
	} else {
		printDebug("error sending 'Prime number calculation'")
		m.end <- false
		return
	}
	if !res {
		printDebug("error receiving 'Prime number calculation result'")
		m.end <- false
		return
	}
	m.end <- true

}

func (m *Marmot) PiCalculation() {
	m.SendAndReceiveData("Pi calculation", true)
}

func (m *Marmot) CountLetter() {
	m.SendAndReceiveData("Count letter", false)
}

func (m *Marmot) Ping() {

	// send 'ping' to client
	m.data = "0Ping\n"
	m.SendAndReceiveData("Ping/Pong", false)
}

// Can be used with a wrapper
// Send the data inside the marmot
// And wait for a marmot response
func (m *Marmot) SendAndReceiveData(functionName string, showMessageSent bool) {

	// wait for start / timeout
	<-m.start
	res := m.writeData(showMessageSent)
	if res {
		res = m.readResponse()
	} else {
		printDebug(fmt.Sprintf("error sending '%s'", functionName))
		m.end <- false
		return
	}

	if !res {
		printDebug(fmt.Sprintf("error receiving '%s'", functionName))
		m.end <- false
		return
	}

	m.end <- true
}

// read client response
// print in DEBUG mode and store the message in response
func (m *Marmot) readResponse() bool {
	fctToExecute := func(ctx context.Context, resultChan chan bool) {
		reader := bufio.NewReader(m.conn)
		type result struct {
			message string
			err     error
		}
		innerChan := make(chan result, 1)

		go func() {
			message, err := reader.ReadString('\n')
			innerChan <- result{message, err}
		}()

		// wait for potential timeout
		select {
		case res := <-innerChan:
			if res.err != nil {
				printError(fmt.Sprintf("Reading client (@"+m.conn.RemoteAddr().String()+") message '%s'", res.err))
				resultChan <- false
			} else {
				printDebug("Message received from @" + m.conn.RemoteAddr().String() + ": " + res.message[:len(res.message)-1])
				m.response = res.message[:len(res.message)-1]
				resultChan <- true
			}
		case <-ctx.Done():
			// timeout
			resultChan <- false
		}
	}

	errorMessage := fmt.Sprintf("Timeout while receiving message from client '@%s'", m.conn.RemoteAddr())
	return m.executeFunctionWithTimeout(TimeoutServerRequestSeconds*time.Second, fctToExecute, errorMessage)
}

// write the Data store in Marmot to the client
// print in DEBUG mode
func (m *Marmot) writeData(show bool) bool {
	fctToExecute := func(ctx context.Context, resultChan chan bool) {
		type result struct {
			n   int
			err error
		}
		innerChan := make(chan result, 1)

		go func() {
			n, err := m.conn.Write([]byte(m.data))
			innerChan <- result{n, err}
		}()

		select {
		case res := <-innerChan:
			if res.err != nil {
				printError(fmt.Sprintf("Sending message to client '(@"+m.conn.RemoteAddr().String()+" %s'", res.err))
				resultChan <- false
			} else {
				printDebugCondition("Message sent: '"+m.data[:len(m.data)-1]+"'", show)
				resultChan <- true
			}
		case <-ctx.Done():
			resultChan <- false
		}
	}

	errorMessage := "Timeout while sending message to client '@%s'"
	return m.executeFunctionWithTimeout(TimeoutServerRequestSeconds*time.Second, fctToExecute, errorMessage)
}

// execute a function with a timeout
// use it as wrapper
func (m *Marmot) executeFunctionWithTimeout(timeout time.Duration, fctToExecute func(ctx context.Context, resultChan chan bool), errorMessage string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resultChan := make(chan bool, 1)
	go fctToExecute(ctx, resultChan)
	// timeout implementation:
	select {
	case res := <-resultChan:
		return res
	case <-ctx.Done():
		printError(errorMessage)
		return false
	}

}
