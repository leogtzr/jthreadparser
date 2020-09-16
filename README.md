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

To get information about which threads are waiting on what (synchronizer states) you can use the SynchronizersByThread() method:

```go
threadDumpFile := "../../threaddumpsamples/13.0.2.0.txt"

threads, err := jthreadparser.ParseFromFile(threadDumpFile)
if err != nil {
    panic(err)
}

syncs := jthreadparser.SynchronizersByThread(&threads)
for thread, threadSyncs := range syncs {
    fmt.Printf("Thread [%s (%s)], synchronizers: %q\n", thread.Name, thread.ID, threadSyncs)
}
```

#### Output
```
...
Thread [http-nio-8080-exec-10 (0x00007f49484d4000)]
	{0x000000060dc25418 java.util.concurrent.locks.AbstractQueuedSynchronizer$ConditionObject ParkingToWaitForState}
Thread [Catalina-utility-1 (0x00007f4948d87000)]
	{0x000000060cf0ffc0 java.util.concurrent.locks.AbstractQueuedSynchronizer$ConditionObject ParkingToWaitForState}
Thread [http-nio-8080-BlockPoller (0x00007f4948f82000)]
	{0x000000060dc19118 sun.nio.ch.Util$2 LockedState}
	{0x000000060dc190c0 sun.nio.ch.EPollSelectorImpl LockedState}
Thread [scheduling-1 (0x00007f494899a800)]
	{0x000000060dc31e60 java.lang.Class for com.thdump.calls.CallResult LockedState}
	{0x000000060dc31ed8 java.lang.Class for com.thdump.calls.Call9 LockedState}
	{0x000000060dc31f48 java.lang.Class for com.thdump.calls.Call8 LockedState}
	{0x000000060dc31fb8 java.lang.Class for com.thdump.calls.Call7 LockedState}
	{0x000000060dc32028 java.lang.Class for com.thdump.calls.Call6 LockedState}
	{0x000000060dc32098 java.lang.Class for com.thdump.calls.Call5 LockedState}
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