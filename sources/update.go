package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
)

// package used to execute self update

// enum used to identify message type
// is it string or binary file for example
type MessageType int

const (
	String MessageType = iota
	BinaryFile
)

// Struct used to send / receive data between client and server
// ID: represents the action id
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
		return fmt.Sprintf("The current message is a binary file")
	default:
		return fmt.Sprintf("%s", string(m.Data))
	}
}

// Create and returns serialize struct
// used to avoid a idiot struct init in others function
func generateNewMessage(id string, messageType MessageType, data []byte) []byte {
	m := Message{id, messageType, data}
	return m.encode()
}

// serializes the struct
func (m Message) encode() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(m)
	if err != nil {
		printError(fmt.Sprintf("During serialisation: %s", err))
		return nil
	}
	return buffer.Bytes()
}

// deserializes the struct
func decodeMessage(buffer bytes.Buffer) *Message {
	decoder := gob.NewDecoder(&buffer)
	var message Message
	err := decoder.Decode(&message)
	if err != nil {
		printError(fmt.Sprintf("During deserialisation: %s", err))
		return nil
	}
	return &message
}
