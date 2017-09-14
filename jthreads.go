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
        fmt.Printf("Name '%s'\n", th.Name)
        /*
        fmt.Printf("Priority '%s'\n", th.Priority)
        fmt.Printf("Native ID '%s'\n", th.NativeID)
        fmt.Printf("State '%s'\n", th.State)
        fmt.Println(th.StackTrace)
        */
    }

    fmt.Println(len(threads))
    for k, v := range jthreadparser.Holds(&threads) {
        fmt.Printf("[%s] => %s\n", k, v)
    }

    /*
    for k, v := range jthreadparser.AwaitingNotification(&threads) {
        fmt.Printf("[%s]\n", k)
        for _, waiting := range v {
            fmt.Println(waiting)
        }
        fmt.Println("~~~~~~~~~~~~~~~~~~~~~")
    }*/

}

