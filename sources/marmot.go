package main

import "net"

// Represents a client
type Marmot struct {
	conn net.Conn
}

// Represents the clients list
type Marmots []*Marmot

type MarmotI interface {
	// send a ping to the client to see if the connection is always up
	Ping() (bool, error)
}

func (m Marmot) Ping() (bool, error) {

}
