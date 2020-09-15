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

	awaiting := jthreadparser.AwaitingNotification(&threads)
	for lock, threads := range awaiting {
		// fmt.Println
		fmt.Printf("Lock: [%q], threads: %q\n\n", lock, threads)
	}

	//for _, th := range threads {

	// if th.State == "BLOCKED" {
	// 	fmt.Println(th)
	// }

	// if th.Name == "http-nio-8080-exec-1" {
	// 	for _, hold := range jthreadparser.HoldsForThread(&th) {
	// 		fmt.Println(hold)
	// 	}
	// }
	//}
}
