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

### Threads with same stacktrace

```go
threads, err := jthreadparser.ParseFromFile("thread_dump.txt")
...
indenticalStackTrace := jthreadparser.IdenticalStackTrace(&threads)

for stackTrace, threadCount := range indenticalStackTrace {
    fmt.Printf("%d threads having this stacktrace: \n", threadCount)
    fmt.Println(stackTrace)
}
```

Output:
```
...
20 threads having this stacktrace: 
at sun.misc.Unsafe.park(Native Method)
- parking to wait for  <0x000000074efc2310> (a java.util.concurrent.locks.AbstractQueuedSynchronizer$ConditionObject)
at java.util.concurrent.locks.LockSupport.park(LockSupport.java:156)
at java.util.concurrent.locks.AbstractQueuedSynchronizer$ConditionObject.await(AbstractQueuedSynchronizer.java:1987)
at java.util.concurrent.LinkedBlockingQueue.take(LinkedBlockingQueue.java:399)
at java.util.concurrent.ThreadPoolExecutor.getTask(ThreadPoolExecutor.java:957)
at java.util.concurrent.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:917)
at java.lang.Thread.run(Thread.java:682)
1 threads having this stacktrace: 
at java.lang.Object.wait(Native Method)
- waiting on <0x000000075c8e5d40> (a atg.service.datacollection.DataCollectorQueue)
at java.lang.Object.wait(Object.java:485)
at atg.service.queue.EventQueue.getElement(EventQueue.java:236)
- locked <0x000000075c8e5d40> (a atg.service.datacollection.DataCollectorQueue)
at atg.service.queue.EventQueue.dispatchQueueElements(EventQueue.java:285)
at atg.service.queue.EventQueue$Handler.run(EventQueue.java:91)
1 threads having this stacktrace: 
at java.lang.Object.wait(Native Method)
...
```