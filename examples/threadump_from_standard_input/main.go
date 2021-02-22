package main

import (
	"fmt"
	"os"

	"github.com/leogtzr/jthreadparser"
)

/*
	reading a thread dump from stdin.
*/
func main() {

	threadDump := jthreadparser.ParseFrom(os.Stdin)

	for _, th := range threadDump.Threads {
		fmt.Println(th)
	}

}
