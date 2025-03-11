package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// enum used to identify message type
// is it string or binary file for example
type MessageType int

const (
	String MessageType = iota
	BinaryFile
)

// Struct used to send / receive data between client and server
// ID: represents the action id
// -1: Update client with new file
// 0: Ping
// 1: Close connection (exit)
// 2: Counting letter
// 3: Calculate if a number is prime
// 4: Calculate pi estimation
// Type: the of the message, string or binary file
// Data: data used to process the message
type Message struct {
	ID   string
	Type MessageType
	Data []byte
}

func (m Message) String() string {
	switch m.Type {
	case BinaryFile:
		return fmt.Sprintf("id: %s The current message is a binary file", m.ID)
	default:
		return fmt.Sprintf("id: %s Data: '%s'", m.ID, string(m.Data))
	}
}

// create new message struct
func createMessage(id string, messageType MessageType, data []byte) *Message {
	return &Message{id, messageType, data}
}

// Create and returns serialize struct
// used to avoid a idiot struct init in others function
func generateNewMessage(id string, messageType MessageType, data []byte) ([]byte, error) {
	m := Message{id, messageType, data}
	return m.encode()
}

func (m Message) encode() ([]byte, error) {
	return encode(m)
}

// deserializes the struct
func decodeMessage(data []byte) (*Message, error) {
	return decode[Message](data)
}

// serializes generic struct
func encode[T any](m T) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(m)
	if err != nil {
		printError(fmt.Sprintf("During serialisation: %s", err))
		return nil, err
	}
	return buffer.Bytes(), nil

}

// generic method used to decode struct
func decode[T any](data []byte) (*T, error) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	var obj T
	err := decoder.Decode(&obj)
	if err != nil {
		printError(fmt.Sprintf("During deserialisation: %s", err))
		return nil, err
	}
	return &obj, err
}

// / returns if the message is a 'exit' message
func (m *Message) isExit() bool {
	return m != nil && m.ID == "1"
}
