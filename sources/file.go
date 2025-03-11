package main

import (
	"fmt"
	"os"
)

// used to send file across the network
// with a version number and the file himself
type File struct {
	Version int
	Data    []byte
}

func createFile(version int, data []byte) *File {
	return &File{version, data}
}

// will take client version and file in local and generate a struct for the file
func generateUpdateFile() (*Message, error) {
	printDebug("Update filename: " + UpdateFilename)
	data, err := os.ReadFile(UpdateFilename)

	if err != nil {
		printError(fmt.Sprintf("During update file generation: %s", err))
		return nil, err
	}
	file, err := createFile(ClientVersion, data).encode()
	if err != nil {
		printError(fmt.Sprintf("During update file generation: %s", err))
		return nil, err
	}
	return createMessage("-1", BinaryFile, file), nil
}

func (f File) encode() ([]byte, error) {
	return encode(f)
}

func decodeFile(data []byte) (*File, error) {
	return decode[File](data)
}
