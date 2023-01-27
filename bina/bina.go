package bina

import (
	"encoding/binary"
	"fmt"
	"os"
)

// #include "stringtable.h"
import "C"

// raw --

type BINAHeader struct {
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

// --

type NodeData struct {
	AdditionalData []byte
	DataBlock      []byte
	StringTable    string
	OffsetTable    string
}

type Node struct {
	Header     NodeHeader
	DataHeader NodeDataHeader
	Data       NodeData
}

type BINA struct {
	Header BINAHeader
	Nodes  []Node
}

func ReadString(stringTable string, offset int) string {
	cstr := C.readString(C.CString(stringTable), C.int(offset))
	return C.GoString(cstr)
}

func Print(file BINA) {
	// INFO
	fmt.Printf("Signature: %s\n", file.Header.Signature)
	fmt.Printf("Version: %s\n", file.Header.Version)
	fmt.Printf("Endian: %c\n", file.Header.EndianFlag)
	fmt.Printf("Node Count: %d\n", file.Header.NodeCount)

	for i, n := range file.Nodes {
		fmt.Printf("----------------\n")
		fmt.Printf("Node %d\n\n", i+1)
		fmt.Printf("Node Type: %s\n", n.Header.Signature)
		fmt.Printf("Node Length: %d\n", n.Header.Length)
		fmt.Printf("Offset Count: %d\n", len(n.Data.OffsetTable))
		fmt.Printf("----------------\n")
	}
}

func Read(path string) BINA {
	f, err := os.Open(path)

	if err == nil {
		header := BINAHeader{}

		err := binary.Read(f, binary.LittleEndian, &header)

		if err != nil {
			fmt.Println(err)
		}

		if string(header.Signature[:]) == "BINA" {
			var nodeList []Node

			for i := 0; i < int(header.NodeCount); i++ {
				nodeHeader := NodeHeader{}

				binary.Read(f, binary.LittleEndian, &nodeHeader)

				if string(nodeHeader.Signature[:]) == "DATA" {
					nodeDataHeader := NodeDataHeader{}

					binary.Read(f, binary.LittleEndian, &nodeDataHeader)

					additionalData := make([]byte, nodeDataHeader.AdditionalDataLength)
					data := make([]byte, nodeHeader.Length-(0x18+uint32(nodeDataHeader.AdditionalDataLength)))

					binary.Read(f, binary.LittleEndian, &additionalData)
					binary.Read(f, binary.LittleEndian, &data)

					offsetTableOffset := nodeDataHeader.StringTableOffset + nodeDataHeader.StringTableLength

					dataBlock := data[:nodeDataHeader.StringTableOffset]
					stringTable := string(data[nodeDataHeader.StringTableOffset:offsetTableOffset])
					offsetTable := string(data[offsetTableOffset : offsetTableOffset+nodeDataHeader.OffsetTableLength])

					nodeData := NodeData{additionalData, dataBlock, stringTable, offsetTable}
					node := Node{nodeHeader, nodeDataHeader, nodeData}

					// fmt.Printf("%s\n", ReadString(stringTable, 0xC))

					nodeList = append(nodeList, node)
				}

				fmt.Printf("----------------\n")
			}

			binaFile := BINA{header, nodeList}
			return binaFile

		} else {
			fmt.Println("Invalid file.")
			return BINA{}
		}
	}
	return BINA{}
}
