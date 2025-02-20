package main

import "os"

func main() {
	argWithoutProg := os.Args[1:]
	marmots := make([]*Marmot, ClientNumber)
	if len(argWithoutProg) > 0 {
		go openConnection(ServerPort, marmots)
		handleMenu(marmots)
	} else {
		connectToServer(ServerIP)
	}
}
