package main

import (
	"fmt"

	"github.com/leogtzr/jthreadparser"
)

func main() {
	threadDumpFile := "../../threaddumpsamples/15.0.txt"

	threadDump, err := jthreadparser.ParseFromFile(threadDumpFile)
	if err != nil {
		panic(err)
	}

	for _, id := range threadDump.SMR {
		fmt.Println(id)
	}
}
