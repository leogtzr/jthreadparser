package main

import (
	"fmt"
	"log"
	"os"

	"github.com/leogtzr/jthreadparser"
)

func main() {
	threadDump, err := jthreadparser.ParseFromFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	for _, th := range threadDump.Threads {
		fmt.Println(th.Name)
		fmt.Println(th.Priority)
		fmt.Println(th.NativeID)
		fmt.Println(th.State)
		fmt.Println(th.StackTrace)
	}
}
