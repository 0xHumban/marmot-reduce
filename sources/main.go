package main

import (
	"fmt"
	"os"
)

func main() {
	printDebug(fmt.Sprintf("Current software version: %d", ClientVersion))
	argWithoutProg := os.Args[1:]
	marmots := make([]*Marmot, ClientNumber)
	if len(argWithoutProg) > 0 {
		go openConnection(ServerPort, marmots)
		handleMenu(marmots)
	} else {
		connectToServer(ServerIP)
	}
}
