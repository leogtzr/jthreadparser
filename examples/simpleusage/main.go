package main

import (
	"fmt"

	"github.com/leogtzr/jthreadparser"
)

func main() {

	threadDumpFile := "../../threaddumpsamples/9.0.4.0-test.txt"

	threads, err := jthreadparser.ParseFromFile(threadDumpFile)
	if err != nil {
		panic(err)
	}

	for _, th := range threads {
		fmt.Printf("'%s' (%s) (%s) [%s]\n", th.Name, th.ID, th.NativeID, th.State)
	}

}
