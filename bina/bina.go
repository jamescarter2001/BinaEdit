package bina

import (
	"encoding/binary"
	"fmt"
	"os"
)

// #include "stringtable.h"
import "C"

type Header struct {
	Signature  [4]byte
	Version    [3]byte
	EndianFlag byte
	FileSize   uint32
	NodeCount  uint16
	Unknown1   uint16
}

type NodeHeader struct {
	Signature [4]byte
	Length    uint32
}

type NodeDataHeader struct {
	StringTableOffset    uint32
	StringTableLength    uint32
	OffsetTableLength    uint32
	AdditionalDataLength uint16
	Padding1             uint16
}

func ReadString(stringTable string, offset int) string {
	cstr := C.readString(C.CString(stringTable), C.int(offset))
	return C.GoString(cstr)
}

func Read(path string) {
	f, err := os.Open(path)

	if err == nil {
		header := Header{}

		err := binary.Read(f, binary.LittleEndian, &header)

		if err != nil {
			fmt.Println(err)
		}

		// headerSize := binary.Size(header)

		if string(header.Signature[:]) == "BINA" {
			// INFO
			fmt.Printf("Signature: %s\n", header.Signature)
			fmt.Printf("Version: %s\n", header.Version)
			fmt.Printf("Endian: %c\n", header.EndianFlag)
			fmt.Printf("Node Count: %d\n", header.NodeCount)

			for i := 0; i < int(header.NodeCount); i++ {
				fmt.Printf("----------------\n")

				fmt.Printf("Node %d\n\n", i+1)

				nodeHeader := NodeHeader{}

				binary.Read(f, binary.LittleEndian, &nodeHeader)
				fmt.Printf("Node Type: %s\n", nodeHeader.Signature)
				fmt.Printf("Node Length: %d\n", nodeHeader.Length)

				if string(nodeHeader.Signature[:]) == "DATA" {
					nodeDataHeader := NodeDataHeader{}

					binary.Read(f, binary.LittleEndian, &nodeDataHeader)

					additionalData := make([]byte, nodeDataHeader.AdditionalDataLength)
					dataBlock := make([]byte, nodeHeader.Length-(0x18+uint32(nodeDataHeader.AdditionalDataLength)))

					binary.Read(f, binary.LittleEndian, &additionalData)
					binary.Read(f, binary.LittleEndian, &dataBlock)

					stringTable := string(dataBlock[nodeDataHeader.StringTableOffset : nodeDataHeader.StringTableOffset+nodeDataHeader.StringTableLength])

					fmt.Printf("%s\n", C.GoString(C.readString(C.CString(stringTable), 12)))
				}

				fmt.Printf("----------------\n")
			}

		} else {
			fmt.Println("Invalid file.")
		}
	}
}
