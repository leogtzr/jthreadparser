# General usage

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

### Most Used Methods

You can check how many threads include a specific Method using the MostUsedMethods() function:

```go
threads, err := jthreadparser.ParseFromFile("thread_dump.txt")
...

mostUsedMethods := jthreadparser.MostUsedMethods(&threads)
for javaMethodName, threadCount := range mostUsedMethods {
    fmt.Printf("%d thread(s) having '%s'\n", threadCount, javaMethodName)
}
```

Output:
```
...
241 thread(s) having 'java.util.concurrent.ThreadPoolExecutor.getTask(ThreadPoolExecutor.java:1074)'
233 thread(s) having 'java.util.concurrent.LinkedBlockingQueue.take(LinkedBlockingQueue.java:442)'
59 thread(s) having 'java.util.concurrent.ThreadPoolExecutor.runWorker(ThreadPoolExecutor.java:1149)'
...
```

