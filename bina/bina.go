package bina

import (
	"encoding/binary"
	"fmt"
	"os"
)

type Header struct {
	Signature  [4]byte
	Version    [3]byte
	EndianFlag byte
	FileSize   uint32
	NodeCount  uint16
	Unknown1   uint16
}

type Node struct {
}

func Read(path string) {
	f, err := os.Open(path)

	if err == nil {
		header := Header{}

		err := binary.Read(f, binary.LittleEndian, &header)

		if err != nil {
			fmt.Println(err)
		}

		if string(header.Signature[:]) == "BINA" {
			fmt.Printf("Signature: %s\n", header.Signature)
			fmt.Printf("Version: %s\n", header.Version)
			fmt.Printf("Endian: %c\n", header.EndianFlag)
			fmt.Printf("Node Count: %d\n", header.NodeCount)
		} else {
			fmt.Println("Invalid file.")
		}
	}
}
