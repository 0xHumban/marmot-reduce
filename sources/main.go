package main

import "os"

func main() {
	argWithoutProg := os.Args[1:]
	if len(argWithoutProg) > 0 {
		openConnection(ServerPort, handleConnection)
	} else {
		connectToServer(ServerIP)
	}
}
