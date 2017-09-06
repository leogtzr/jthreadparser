package main

import (
    "fmt"
    "log"
    "os"
    "./jthreadparser"
)

func main() {

    if len(os.Args) == 1 {
        return
    }

    threads, err := jthreadparser.Parse(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }

    for _, th := range threads {
        fmt.Println(th)
    }

}