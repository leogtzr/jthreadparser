package main

import (
	"fmt"

	"github.com/leogtzr/jthreadparser"
)

func main() {
	threadDumpFile := "../../threaddumpsamples/13.0.2.0.txt"

	threads, err := jthreadparser.ParseFromFile(threadDumpFile)
	if err != nil {
		panic(err)
	}

	syncs := jthreadparser.SynchronizersByThread(&threads)
	for thread, threadSyncs := range syncs {
		fmt.Printf("Thread [%s (%s)]\n", thread.Name, thread.ID)
		for _, s := range threadSyncs {
			fmt.Printf("\t%s\n", s)
		}
	}

}
