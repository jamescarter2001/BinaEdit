package main

import (
	"carter.io/binaedit/bina"
	"carter.io/binaedit/gedit"
)

func main() {
	binaFile := bina.Read("w1f01_obj_area01.gedit")
	bina.Print(binaFile)
	gedit.Read(binaFile.Nodes[0])
}
