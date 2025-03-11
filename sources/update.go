package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// package used to execute self update

// wil send to all clients the latest update file
func (ms Marmots) SendUpdateFile() {
	ms.Pings()
	printDebug("Start Send update file")
	// get the local file and store it inside the marmot data
	data, err := generateUpdateFile()
	if err != nil {
		printError(fmt.Sprintf("during send update file, generating update file: %s", err))
		return
	}
	for _, m := range ms {
		if m != nil {
			m.data = data
		}
	}

	ms.performAction((*Marmot).SendUpdateFile)
	ms.WaitEnd()
	printDebug("End Send update file")
}

// The data inside the marmot have to be initialized before using this function
func (m *Marmot) SendUpdateFile() {
	// check if the data is correct
	if !m.isUpdateFile(true) {
		printError("send update file to marmot, data inside the marmot is not initialized")
		return
	}
	// send current data to client
	m.SendAndReceiveData("UpdateFile", true)
}

// returns if the data inside the Marmot is an Update File
// checkData: bool: if we have to check inside the data or the response attribut
func (m Marmot) isUpdateFile(checkData bool) bool {
	if checkData {
		return m.data != nil && m.data.ID == "-1" && m.data.Type == BinaryFile
	} else {

		return m.response != nil && m.response.ID == "-1" && m.response.Type == BinaryFile
	}
}

// will store the file from date inside the Marmot, run it and kill the old one
func (m *Marmot) SelfUpdateClient() (bool, error) {
	// check if the data is correct
	if !m.isUpdateFile(false) {
		printError("client side update client, data inside the marmot is not initialized / valid")
		return false, fmt.Errorf("client side update client, data inside the marmot is not initialized / valid")
	}
	fileData, err := decodeFile(m.response.Data)
	if err != nil {
		printError(fmt.Sprintf("client side update client, decodin File: %s", err))
		return false, err
	}

	// check versions
	printDebug(fmt.Sprintf("Current version: %d vs Version received: %d", ClientVersion, fileData.Version))
	if fileData.Version <= ClientVersion {
		printDebug("The client is already using the latest client version")
		return false, nil
	}

	// write file
	filename := fmt.Sprintf("%s%d", UpdateFilePath, fileData.Version)
	printDebug(fmt.Sprintf("New client file: '%s'", filename))

	err = os.WriteFile(filename, fileData.Data, 0755)
	if err != nil {
		printError(fmt.Sprintf("client side update client, writing file: %s", err))
		return false, err
	}

	return executeFile(filename)

}

// execute the file and exit the current client
func executeFile(filename string) (bool, error) {
	cmd := exec.Command("./" + filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err := cmd.Start()

	if err != nil {
		printError(fmt.Sprintf("client side update client, executing new file (%s): %s", filename, err))
		return false, err
	}
	return true, nil
}
