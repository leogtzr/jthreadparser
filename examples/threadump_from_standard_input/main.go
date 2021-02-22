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

	threads := jthreadparser.ParseFrom(os.Stdin)

	for _, th := range threads {
		fmt.Println(th)
	}

}
