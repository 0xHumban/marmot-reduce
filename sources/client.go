package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var ClientVersion = 3

var UpdateFilePath = "build/client-"
var UpdateFilename = fmt.Sprintf("%s%d", UpdateFilePath, ClientVersion)

const LocalServerIP = "192.168.1.25:8080"
const ServerIP = "127.0.0.1:8080"
const RetryDelais = 5

// handle connection client side
// client waiting for server instructions
// 'exit': connection closed
// '1': count 'e' in response
// returns if the connection has been asked by server
func (m *Marmot) handleConnectionClientSide() bool {
	defer m.conn.Close()
	for !m.response.isExit() {
		res := m.readResponse()
		if !res {
			printDebug("ERROR reading server response")
			return false
		}
		if m.response.isExit() {
			printDebug("EXIT request received")
			return true
		}

		m.treatServerResponse()

	}
	return true
}

// treats the server response
// choose whats the next step, which function the client have to execute
func (m *Marmot) treatServerResponse() {
	switch m.response.Type {
	case BinaryFile:
		m.treatBinaryFileServerResponse()
	default:
		m.treatStringServerResponse()
	}
}

func (m *Marmot) treatBinaryFileServerResponse() {
	// self update client request
	if m.response.ID == "-1" {
		printDebug("Self update client request received")
		res, err := m.SelfUpdateClient()
		if err != nil {
			m.data = createMessage("-1", String, []byte(fmt.Sprintf("error during self updating client: %s", err)))
		} else {
			m.data = createMessage("-1", String, []byte("Marmot has been updated"))
		}
		_ = m.writeData(true)
		if res {
			printDebug("File executed")
			// close the current client
			os.Exit(0)
		}
	}
}

func (m *Marmot) treatStringServerResponse() {
	// Ping request
	if m.response.ID == "0" {
		printDebug("Ping pong request received")
		m.data = createMessage("0", String, []byte(fmt.Sprintf("'Pong' from @%s", m.conn.LocalAddr().String())))
		_ = m.writeData(true)
		printDebug("Ping pong response sent")

	} else if
	// count 'e' in response
	m.response.ID == "2" {
		printDebug("Start counting letter occurrences\n")
		letterToCount := string(m.response.Data)[0]
		calculus := fmt.Sprintf("%d", countLetterOccurrence(string(m.response.Data)[1:], rune(letterToCount)))
		printDebug(fmt.Sprintf("End couting letter occurrences -- Result for '%c': %s\n", rune(letterToCount), calculus))
		m.data = createMessage("2", String, []byte(calculus))
		m.writeData(true)

	} else if
	// calculate if a number is prime in a given range
	m.response.ID == "3" {
		printDebug("Start prime number calculation\n")
		parts := strings.Split(string(m.response.Data), "@")
		println(string(m.response.Data))
		if len(parts) != 3 {
			printError("Invalid format")
			m.data = createMessage("3", String, []byte(fmt.Sprintf("%s\n", "Invalid format")))
			m.writeData(true)
			return
		}

		potentialPrime, err1 := strconv.Atoi(parts[0])
		start, err2 := strconv.Atoi(parts[1])
		end, err3 := strconv.Atoi(parts[2])

		if err1 != nil || err2 != nil || err3 != nil {
			printError("during conversion")
			m.data = createMessage("3", String, []byte(fmt.Sprintf("%s\n", "Conversion error")))
			m.writeData(true)
			return
		}
		calculus := fmt.Sprintf("%d", calculatePrimeNumber(potentialPrime, start, end))
		printDebug(fmt.Sprintf("End prime number calculation -- Result: %s\n", calculus))
		m.data = createMessage("3", String, []byte(calculus))
		m.writeData(true)

	} else if
	// calculate pi estiumation
	m.response.ID == "4" {
		printDebug("Start pi estimation\n")
		samples, err := strconv.Atoi(string(m.response.Data))
		if err != nil {
			printError("during conversion")
			m.data = createMessage("4", String, []byte(fmt.Sprintf("%s\n", "Conversion error")))
			m.writeData(true)
			return
		}
		calculus := fmt.Sprintf("%d", calculePiChunk(samples))
		printDebug(fmt.Sprintf("End pi estimation -- Result: %s\n", calculus))
		m.data = createMessage("4", String, []byte(calculus))
		m.writeData(true)

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

// calculate PI chunk
// returns of points those that fall inside the quarter circle
// IN: n: the number of points to generate
func calculePiChunk(n int) int {
	inside := 0
	for i := 0; i < n; i++ {
		x, y := rand.Float64(), rand.Float64()
		if x*x+y*y <= 1 {
			inside++
		}
	}

	return inside
}

// saves the file stored in data on local system
// start it and kill the old one
// TODO: add verification

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
			marmot := NewMarmot(conn)
			connectionClosedProperly = marmot.handleConnectionClientSide()
		}
		if !connectionClosedProperly {
			time.Sleep(RetryDelais * time.Second)
		}
	}
}
