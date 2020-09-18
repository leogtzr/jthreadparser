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

	threads, err := jthreadparser.ParseFrom(os.Stdin)
	if err != nil {
		panic(err)
	}

	for _, th := range threads {
		fmt.Println(th)
	}

}
