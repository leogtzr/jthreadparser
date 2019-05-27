```go
threads, err := jthreadparser.ParseFromFile(os.Args[1])

if err != nil {
    log.Fatal(err)
}

for _, th := range threads {
    fmt.Println(th.Name)
    fmt.Println(th.Priority)
    fmt.Println(th.NativeID)
    fmt.Println(th.State)
    fmt.Println(th.StackTrace)
}
```

### Synchronizers

To get information about which threads are waiting on what you can use the AwaitingNotification() method:

```go
threadsWaiting := jthreadparser.AwaitingNotification(&threads)
for k, v := range threadsWaiting {
    fmt.Printf("%d threads waiting for notification on lock %s\n", len(v), k.LockID)
    for _, threadWaiting := range v {
        fmt.Println("\t", threadWaiting.Name)
    }
    fmt.Println()
}
```