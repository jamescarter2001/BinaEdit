package gedit

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"carter.io/binaedit/bina"
)

type GEditHeader struct {
	Padding1                uint64
	Padding2                uint64
	ObjectOffsetTableOffset uint64
	ObjectCount             uint64
	ObjectCount2            uint64
	Padding3                uint64
}

type GEditDataHeader struct {
	ObjectOffsetTable []uint64
}

type ForcesObjectReference struct {
	ID      uint16
	GroupID uint16
}

type ObjectEntry struct {
	Padding1              uint64
	ObjectTypeOffset      uint64
	ObjectNameOffset      uint64
	ID                    ForcesObjectReference
	ParentID              ForcesObjectReference
	Position              [3]float32
	Rotation              [3]float32
	ChildPositionOffset   [3]float32
	ChildRotationOffset   [3]float32
	TagsOffsetTableOffset uint64
	TagCount              uint64
	TagCount2             uint64
	Padding3              uint64
	ParametersOffset      uint64
}

func Read(node bina.Node) {
	header := GEditHeader{}
	r := bytes.NewReader(node.Data.DataBlock)

	err := binary.Read(r, binary.LittleEndian, &header)

	if err == nil {
		objectOffsetTable := make([]uint64, header.ObjectCount)
		binary.Read(r, binary.LittleEndian, &objectOffsetTable)

		fmt.Println("DONE")
	} else {
		fmt.Println("Not a gedit file.")
	}
}
