package jthreadparser

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const (
	threadInfo      = `"Attach Listener" daemon prio=10 tid=0x00002aaab74c5000 nid=0x2ea5 waiting on condition [0x0000000000000000]`
	threadStateLine = `   java.lang.Thread.State: TIMED_WAITING (parking)`

	threadInformation = `"HDScanner" prio=10 tid=0x00002aaaebf6e800 nid=0x4367 runnable [0x00000000451cf000]
   java.lang.Thread.State: RUNNABLE
	at java.io.UnixFileSystem.getBooleanAttributes0(Native Method)
	at java.io.UnixFileSystem.getBooleanAttributes(UnixFileSystem.java:228)
	at java.io.File.isHidden(File.java:804)

   Locked ownable synchronizers:
	- <0x00000007500b7250> (a java.util.concurrent.locks.ReentrantLock$NonfairSync)
	`

	daemonThreadInformation = `"Attach Listener" #6085 daemon prio=9 os_prio=0 tid=0x00007f90d0106000 nid=0x18a1 runnable [0x0000000000000000]
	java.lang.Thread.State: RUNNABLE
 
	Locked ownable synchronizers:
	 - None`

	threadInfoWithLocks = `"default task-23" #349 prio=5 os_prio=0 tid=0x00007f8fe400c800 nid=0x72fa waiting for monitor entry [0x00007f8f7228e000]
	java.lang.Thread.State: BLOCKED (on object monitor)
	 at java.security.Provider.getService(Provider.java:1039)
	 - locked <0x0000000682e5f948> (a sun.security.provider.Sun)
	 at sun.security.jca.ProviderList.getService(ProviderList.java:332)
	 at sun.security.jca.GetInstance.getInstance(GetInstance.java:157)
	 at java.security.Security.getImpl(Security.java:695)
	 at java.security.MessageDigest.getInstance(MessageDigest.java:167)
	 at sun.security.rsa.RSASignature.<init>(RSASignature.java:79)
	 at sun.security.rsa.RSASignature$SHA512withRSA.<init>(RSASignature.java:305)
	 at sun.reflect.NativeConstructorAccessorImpl.newInstance0(Native Method)
	 at sun.reflect.NativeConstructorAccessorImpl.newInstance(NativeConstructorAccessorImpl.java:62)
	 at sun.reflect.DelegatingConstructorAccessorImpl.newInstance(DelegatingConstructorAccessorImpl.java:45)
	 at java.lang.reflect.Constructor.newInstance(Constructor.java:423)
	 at java.security.Provider$Service.newInstance(Provider.java:1595)
	 at java.security.Signature$Delegate.newInstance(Signature.java:1020)
	 at java.security.Signature$Delegate.chooseProvider(Signature.java:1114)
	 - locked <0x00000007bc531138> (a java.lang.Object)
	 at java.security.Signature$Delegate.engineInitSign(Signature.java:1188)
	 at java.security.Signature.initSign(Signature.java:553)
	 at sun.security.ssl.HandshakeMessage$ECDH_ServerKeyExchange.<init>(HandshakeMessage.java:1031)
	 at sun.security.ssl.ServerHandshaker.clientHello(ServerHandshaker.java:971)
	 at sun.security.ssl.ServerHandshaker.processMessage(ServerHandshaker.java:228)
	 at sun.security.ssl.Handshaker.processLoop(Handshaker.java:1052)
	 at sun.security.ssl.Handshaker$1.run(Handshaker.java:992)
	 at sun.security.ssl.Handshaker$1.run(Handshaker.java:989)
	 at java.security.AccessController.doPrivileged(Native Method)
	 at sun.security.ssl.Handshaker$DelegatedTask.run(Handshaker.java:1467)
	 - locked <0x00000007bbbac500> (a sun.security.ssl.SSLEngineImpl)
	 at io.undertow.protocols.ssl.SslConduit$5.run(SslConduit.java:1021)
	 at java.util.concurrent.ThreadPoolExecutor.runWorker(ThreadPoolExecutor.java:1149)
	 at java.util.concurrent.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:624)
	 at java.lang.Thread.run(Thread.java:748)
 
	Locked ownable synchronizers:
	 - <0x00000006a43d5c08> (a java.util.concurrent.ThreadPoolExecutor$Worker)
 `
	threadDumpSample = `2017-06-02 19:02:52
Full thread dump Java HotSpot(TM) 64-Bit Server VM (20.141-b32 mixed mode):

"Attach Listener" daemon prio=10 tid=0x00002aaab74c5000 nid=0x2ea5 waiting on condition [0x0000000000000000]
	java.lang.Thread.State: RUNNABLE

	Locked ownable synchronizers:
	- None

"RMI TCP Connection(idle)" daemon prio=10 tid=0x00002aaac69b1800 nid=0x2dec waiting on condition [0x0000000051388000]
	java.lang.Thread.State: TIMED_WAITING (parking)
	at sun.misc.Unsafe.park(Native Method)
	- parking to wait for  <0x0000000740c99708> (a java.util.concurrent.SynchronousQueue$TransferStack)
	at java.util.concurrent.locks.LockSupport.parkNanos(LockSupport.java:196)
	at java.util.concurrent.SynchronousQueue$TransferStack.awaitFulfill(SynchronousQueue.java:424)
	at java.lang.Thread.run(Thread.java:682)

	Locked ownable synchronizers:
	- None

"RMI TCP Connection(idle)" daemon prio=10 tid=0x00002aaad5029000 nid=0x2bf8 waiting on condition [0x0000000051287000]
	java.lang.Thread.State: TIMED_WAITING (parking)
	at sun.misc.Unsafe.park(Native Method)
	- parking to wait for  <0x0000000740c99708> (a java.util.concurrent.SynchronousQueue$TransferStack)
	at java.util.concurrent.locks.LockSupport.parkNanos(LockSupport.java:196)
	at java.util.concurrent.SynchronousQueue$TransferStack.awaitFulfill(SynchronousQueue.java:424)
	at java.lang.Thread.run(Thread.java:682)

	Locked ownable synchronizers:
	- None`
)

func TestThreadInfo(t *testing.T) {

	type testCase struct {
		threadInfoHeadLine string
		want               ThreadInfo
	}

	tests := []testCase{
		testCase{
			threadInfoHeadLine: `"http-nio-8080-BlockPoller" #18 daemon prio=5 os_prio=0 cpu=13.17ms elapsed=116.26s tid=0x00007f4948f82000 nid=0xa35f runnable  [0x00007f48a7dfe000]`,
			want: ThreadInfo{
				Daemon:   true,
				Name:     `http-nio-8080-BlockPoller`,
				NativeID: "0xa35f",
				Priority: "5",
				ID:       "0x00007f4948f82000",
			},
		},
		testCase{
			threadInfoHeadLine: `"GC Thread#6" os_prio=0 cpu=12.53ms elapsed=106.94s tid=0x00007f1914009000 nid=0xb090 runnable`,
			want: ThreadInfo{
				Daemon:   false,
				Name:     `GC Thread#6`,
				NativeID: "0xb090",
				Priority: "",
				ID:       `0x00007f1914009000`,
			},
		},
		testCase{
			threadInfoHeadLine: `"GC Thread#0" os_prio=0 cpu=25.37ms elapsed=107.25s tid=0x00007f195c06a800 nid=0xb079 runnable`,
			want: ThreadInfo{
				Daemon:   false,
				Name:     `GC Thread#0`,
				NativeID: `0xb079`,
				Priority: "",
				ID:       `0x00007f195c06a800`,
			},
		},
		testCase{
			threadInfoHeadLine: `"scheduling-1" #31 prio=5 os_prio=0 cpu=23.59ms elapsed=116.23s tid=0x00007f494899a800 nid=0xa36c waiting on condition  [0x00007f48a70f0000]`,
			want: ThreadInfo{
				Daemon:   false,
				Name:     `scheduling-1`,
				NativeID: `0xa36c`,
				ID:       `0x00007f494899a800`,
				Priority: "5",
			},
		},
	}

	for _, tc := range tests {
		got := extractThreadInfoFromLine(tc.threadInfoHeadLine)
		if got.ID != tc.want.ID {
			t.Errorf("got=[%s], want=[%s]", got.ID, tc.want.ID)
		}
		if got.Daemon != tc.want.Daemon {
			t.Errorf("got=[%t], want=[%t]", got.Daemon, tc.want.Daemon)
		}
		if got.Name != tc.want.Name {
			t.Errorf("got=[%s], want=[%s]", got.Name, tc.want.Name)
		}
		if got.NativeID != tc.want.NativeID {
			t.Errorf("got=[%s], want=[%s]", got.NativeID, tc.want.NativeID)
		}
		if got.Priority != tc.want.Priority {
			t.Errorf("got=[%s], want=[%s]", got.Priority, tc.want.Priority)
		}
		if got.ID != tc.want.ID {
			t.Errorf("got=[%s], want=[%s]", got.ID, tc.want.ID)
		}
	}
}

func TestExtractThreadState(t *testing.T) {

	type testCase struct {
		threadStateLine, want string
	}

	tests := []testCase{
		testCase{
			threadStateLine: `   java.lang.Thread.State: BLOCKED (on object monitor)`,
			want:            "BLOCKED",
		},
		testCase{
			threadStateLine: `java.lang.Thread.State: TIMED_WAITING (on object monitor)`,
			want:            `TIMED_WAITING`,
		},
	}

	for _, tc := range tests {
		got := extractThreadState(tc.threadStateLine)
		if got != tc.want {
			t.Errorf("got=[%s], want=[%s]", got, tc.want)
		}
	}
}

func TestShouldIdentifyDaemonThread(t *testing.T) {
	threads, err := ParseFrom(strings.NewReader(daemonThreadInformation))
	if err != nil {
		t.Error("Error parsing daemon thread")
	}
	if len(threads) != 1 {
		t.Error("Error parsing single daemon thread dump")
	}
	if !threads[0].Daemon {
		t.Error("Thread should be daemon")
	}
}

func TestShouldTagCorrectlyDaemonThread(t *testing.T) {

	type testCase struct {
		threadInfo ThreadInfo
		want       string
	}

	tests := []testCase{
		testCase{
			threadInfo: ThreadInfo{ID: "0x00007f90d0106000", Name: "Attach Listener", State: "RUNNABLE", Daemon: true},
			want:       `Thread Id: '0x00007f90d0106000' (daemon), Name: 'Attach Listener', State: 'RUNNABLE'`,
		},

		testCase{
			threadInfo: ThreadInfo{ID: "0x00007f90d0106000", Name: "Attach Listener", State: "RUNNABLE", Daemon: false},
			want:       `Thread Id: '0x00007f90d0106000', Name: 'Attach Listener', State: 'RUNNABLE'`,
		},
	}

	for _, tc := range tests {
		got := tc.threadInfo.String()
		if got != tc.want {
			t.Errorf("got=[%s], expected=[%s]", tc.threadInfo.String(), tc.want)
		}
	}
}

func TestHolds(t *testing.T) {
	expectedNumberOfThreads := 1
	expectedNumberOfLocksInThread := 3

	expectedLocks := []Locked{
		Locked{LockID: "0x0000000682e5f948", LockecObjectName: "sun.security.provider.Sun"},
		Locked{LockID: "0x00000007bc531138", LockecObjectName: "java.lang.Object"},
		Locked{LockID: "0x00000007bbbac500", LockecObjectName: "sun.security.ssl.SSLEngineImpl"},
	}

	threads, err := ParseFrom(strings.NewReader(threadInfoWithLocks))
	if err != nil {
		t.Error("Error parsing thread information.")
	}
	holds := Holds(&threads)
	if len(holds) != expectedNumberOfThreads {
		t.Errorf("got=[%d], expected=[%d] for holds map length.", len(holds), expectedNumberOfThreads)
	}
	for _, locks := range holds {
		if len(locks) != expectedNumberOfLocksInThread {
			t.Errorf("It should have identified %d, got=%d", expectedNumberOfThreads, len(locks))
		}
		for i, lock := range locks {
			if expectedLocks[i] != lock {
				t.Errorf("got=[%s], expected=[%s]", lock, expectedLocks[i])
			}
		}
	}

}

func TestHoldsForThread(t *testing.T) {
	expectedNumberOfLocksInThreadWithLockInfo := 3
	expectedLocks := []Locked{
		Locked{LockID: "0x0000000682e5f948", LockecObjectName: "sun.security.provider.Sun"},
		Locked{LockID: "0x00000007bc531138", LockecObjectName: "java.lang.Object"},
		Locked{LockID: "0x00000007bbbac500", LockecObjectName: "sun.security.ssl.SSLEngineImpl"},
	}

	threads, err := ParseFrom(strings.NewReader(threadInfoWithLocks))
	if err != nil {
		t.Error("Error parsing thread information.")
	}
	for _, thread := range threads {
		holds := HoldsForThread(&thread)
		if len(holds) != expectedNumberOfLocksInThreadWithLockInfo {
			t.Errorf("expected=[%d], got=[%d]", expectedNumberOfLocksInThreadWithLockInfo, len(holds))
		}
		for i, lock := range holds {
			if expectedLocks[i] != lock {
				t.Errorf("got=[%s], expected=[%s]", lock, expectedLocks[i])
			}
		}

	}

	// Test now with a thread having a stacktrace without locking/hold information
	// It should return an empty slice.
	threads, err = ParseFrom(strings.NewReader(threadInformation))
	if err != nil {
		t.Error("Error parsing thread information.")
	}

	for _, thread := range threads {
		holds := HoldsForThread(&thread)
		if len(holds) != 0 {
			t.Errorf("Should be empty")
		}
	}

}

func TestParseFromFile(t *testing.T) {
	file, err := ioutil.TempFile("", "xxxx.")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(file.Name())

	err = ioutil.WriteFile(file.Name(), []byte(threadDumpSample), 0644)
	if err != nil {
		t.Error(err)
	}

	expectedNumberOfThreadsInSample := 3
	threads, err := ParseFromFile(file.Name())

	if len(threads) != expectedNumberOfThreadsInSample {
		t.Errorf("got=[%d], expected=[%d]", len(threads), expectedNumberOfThreadsInSample)
	}

	threads, err = ParseFromFile("does_not_exist...")
	if threads != nil {
		t.Errorf("got=[%q], expected a nil reference", threads)
	}
}

func TestTopMethodsInThreads(t *testing.T) {

	type testCase struct {
		javaMethodName string
		want           int
	}

	threads, err := ParseFromFile("tdump.sample")
	if err != nil {
		t.Error("Unable to parse thread dump sample file")
	}

	tests := []testCase{
		testCase{
			javaMethodName: "java.lang.Object.wait(Native Method)",
			want:           82,
		},
	}

	for _, tc := range tests {
		got := MostUsedMethods(&threads)
		if got[tc.javaMethodName] != tc.want {
			t.Errorf("Should have identified %d threads, got=%d", tc.want, got[tc.javaMethodName])
		}
	}

}

func TestIdenticalStrackTrace(t *testing.T) {

	type testCase struct {
		stacktrace string
		want       int
	}

	tests := []testCase{
		testCase{
			stacktrace: `at sun.misc.Unsafe.park(Native Method)
- parking to wait for  <0x0000000750785368> (a java.util.concurrent.locks.AbstractQueuedSynchronizer$ConditionObject)
at java.util.concurrent.locks.LockSupport.park(LockSupport.java:156)
at java.util.concurrent.locks.AbstractQueuedSynchronizer$ConditionObject.await(AbstractQueuedSynchronizer.java:1987)
at java.util.concurrent.LinkedBlockingQueue.take(LinkedBlockingQueue.java:399)
at java.util.concurrent.ThreadPoolExecutor.getTask(ThreadPoolExecutor.java:957)
at java.util.concurrent.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:917)
at java.lang.Thread.run(Thread.java:682)`,
			want: 20,
		},
	}

	threads, err := ParseFromFile("tdump.sample")
	if err != nil {
		t.Error("Unable to parse thread dump sample file")
	}

	for _, tc := range tests {
		got := IdenticalStackTrace(&threads)
		if got[tc.stacktrace] != tc.want {
			t.Errorf("Should have identified %d (got=%d) threads with the following stacktrace:\n%s\nEND",
				tc.want, got[tc.stacktrace], tc.stacktrace)
		}
	}

}

func TestVerifyNumberOfThreadsInSamples(t *testing.T) {
	type testCase struct {
		sampleFileName string
		want           int
	}

	tests := []testCase{
		{"threaddumpsamples/10.0.2.0.txt", 51},
		{"threaddumpsamples/10.0.2.2.txt", 51},
		{"threaddumpsamples/10.0.2.3.txt", 51},
		{"threaddumpsamples/11.0.2.0.txt", 43},
		{"threaddumpsamples/11.0.2.1.txt", 46},
		{"threaddumpsamples/11.0.2.2.txt", 50},
		{"threaddumpsamples/11.0.2.3.txt", 50},
		{"threaddumpsamples/11.0.8.0-amazon.txt", 99},
		{"threaddumpsamples/11.0.8.1-amazon.txt", 99},
		{"threaddumpsamples/11.0.8.2-amazon.txt", 48},
		{"threaddumpsamples/11.0.8.3-amazon.txt", 48},
		{"threaddumpsamples/12.0.2.0.txt", 49},
		{"threaddumpsamples/12.0.2.1.txt", 48},
		{"threaddumpsamples/12.0.2.2.txt", 48},
		{"threaddumpsamples/12.0.2.3.txt", 48},
		{"threaddumpsamples/13.0.2.0.txt", 42},
		{"threaddumpsamples/13.0.2.1.txt", 42},
		{"threaddumpsamples/13.0.2.2.txt", 50},
		{"threaddumpsamples/13.0.2.3.txt", 48},
		{"threaddumpsamples/14.0.1.0.txt", 51},
		{"threaddumpsamples/14.0.1.1.txt", 51},
		{"threaddumpsamples/14.0.1.2.txt", 51},
		{"threaddumpsamples/14.0.1.3.txt", 50},
		{"threaddumpsamples/15.0.txt", 43},
		{"threaddumpsamples/15.1.txt", 51},
		{"threaddumpsamples/15.2.txt", 51},
		{"threaddumpsamples/15.3.txt", 49},
		{"threaddumpsamples/1.8-amazon.0.txt", 37},
		{"threaddumpsamples/1.8-amazon.1.txt", 37},
		{"threaddumpsamples/1.8-amazon.2.txt", 37},
		{"threaddumpsamples/1.8-amazon.3.txt", 37},
		{"threaddumpsamples/9.0.4.0.txt", 51},
		{"threaddumpsamples/9.0.4.1.txt", 51},
		{"threaddumpsamples/9.0.4.2.txt", 51},
		{"threaddumpsamples/9.0.4.3.txt", 51},
		{"threaddumpsamples/tdump2.sample", 58},
		{"threaddumpsamples/tdump.jdk11.idea.txt", 56},
		{"threaddumpsamples/tdump.sample", 203},
		{"threaddumpsamples/x.txt", 0},
	}

	for _, tc := range tests {
		threads, err := ParseFromFile(tc.sampleFileName)
		if err != nil {
			t.Error(err)
		}
		if len(threads) != tc.want {
			t.Errorf("got=[%d], want=[%d]", len(threads), tc.want)
		}
	}

}

func TestNoThreadDumpFile(t *testing.T) {

	type testCase struct {
		sampleFileName string
		want           int
	}

	tests := []testCase{
		{"threaddumpsamples/x.txt", 0},
	}

	for _, tc := range tests {
		threads, err := ParseFromFile(tc.sampleFileName)
		if err != nil {
			t.Error(err)
		}

		if len(threads) != tc.want {
			t.Errorf("got=[%d], want=[%d]", len(threads), tc.want)
		}
	}

}

func TestHasRunnableState(t *testing.T) {

	type testCase struct {
		threadHeaderLine string
		want             bool
	}

	tests := []testCase{
		testCase{
			threadHeaderLine: `GC Thread#3" os_prio=0 cpu=21.13ms elapsed=146.91s tid=0x00007f1914004000 nid=0xb08d runnable  `,
			want:             true,
		},
		testCase{
			threadHeaderLine: `"Finalizer" #3 daemon prio=8 os_prio=0 cpu=0.77ms elapsed=741.83s tid=0x00007f6f5c29e000 nid=0x6b90 in Object.wait()  [0x00007f6f13ffe000]`,
			want:             false,
		},
	}

	for _, tc := range tests {
		got := hasRunnableState(tc.threadHeaderLine)
		if got != tc.want {
			t.Errorf("got=[%t], want=[%t]", got, tc.want)
		}
	}

}

func TestWaitingOnConditionState(t *testing.T) {
	type testCase struct {
		threadHeaderLine string
		want             bool
	}

	tests := []testCase{
		testCase{
			threadHeaderLine: `"VM Periodic Task Thread" os_prio=0 tid=0x00007f74bc22d800 nid=0xd9c4 waiting on condition `,
			want:             true,
		},
		testCase{
			threadHeaderLine: `"Finalizer" #3 daemon prio=8 os_prio=0 cpu=0.77ms elapsed=741.83s tid=0x00007f6f5c29e000 nid=0x6b90 in Object.wait()  [0x00007f6f13ffe000]`,
			want:             false,
		},
	}

	for _, tc := range tests {
		got := hasWaitingOnState(tc.threadHeaderLine)
		if got != tc.want {
			t.Errorf("got=[%t], want=[%t]", got, tc.want)
		}
	}
}

func TestHasThreadHeaderInformation(t *testing.T) {
	type testCase struct {
		threadHeader string
		want         bool
	}

	tests := []testCase{
		testCase{
			threadHeader: `"VM Thread" os_prio=0 cpu=235.82ms elapsed=147.51s tid=0x00007f904819cef0 nid=0xbfd2 runnable`,
			want:         true,
		},
		testCase{
			threadHeader: `"http-nio-8080-exec-8" #27 daemon prio=5 os_prio=0 cpu=1.67ms elapsed=144.34s tid=0x00007f9049128ec0 nid=0xc03f waiting for monitor entry  [0x00007f8faf0f0000]`,
			want:         true,
		},
		testCase{
			threadHeader: `"OPPORTUNITY"`,
			want:         false,
		},
		testCase{
			threadHeader: `Hello!`,
			want:         false,
		},
	}

	for _, tc := range tests {
		got := hasThreadHeaderInformation(tc.threadHeader)
		if got != tc.want {
			t.Errorf("got=[%t], want=[%t]", got, tc.want)
		}
	}

}

func TestCheckThreadInfo(t *testing.T) {

	threads, err := ParseFromFile("threaddumpsamples/9.0.4.0-test.txt")
	if err != nil {
		t.Error(err)
	}

	const expectedNumberOfThreadsInSampleFile = 5
	if len(threads) != expectedNumberOfThreadsInSampleFile {
		t.Errorf("got=[%d] threads, expected=[%d]", len(threads), expectedNumberOfThreadsInSampleFile)
	}

	type testCase struct {
		threadName string
		isDaemon   bool
		threadID   string
		nativeID   string
		state      string
	}

	tests := []testCase{
		testCase{
			threadName: `Attach Listener`,
			isDaemon:   true,
			threadID:   `0x00007f321c001000`,
			nativeID:   `0x5ac6`,
			state:      `RUNNABLE`,
		},
		testCase{
			threadName: `DestroyJavaVM`,
			isDaemon:   false,
			threadID:   `0x00007f32b4012000`,
			nativeID:   `0x5934`,
			state:      `RUNNABLE`,
		},
		testCase{
			threadName: `scheduling-1`,
			isDaemon:   false,
			threadID:   `0x00007f32b556c000`,
			nativeID:   `0x596c`,
			state:      `TIMED_WAITING`,
		},
		testCase{
			threadName: `http-nio-8080-Acceptor`,
			isDaemon:   true,
			threadID:   `0x00007f32b53d7000`,
			nativeID:   `0x596b`,
			state:      `RUNNABLE`,
		},
		testCase{
			threadName: `http-nio-8080-ClientPoller`,
			isDaemon:   true,
			threadID:   `0x00007f32b53f1000`,
			nativeID:   `0x596a`,
			state:      `RUNNABLE`,
		},
	}

	for i := 0; i < expectedNumberOfThreadsInSampleFile; i++ {
		got := threads[i]
		if got.Name != tests[i].threadName {
			t.Errorf("got=[%s], expected=[%s]", got.Name, tests[i].threadName)
		}

		if got.Daemon != tests[i].isDaemon {
			t.Errorf("got=[%t], want=[%t]", got.Daemon, tests[i].isDaemon)
		}

		if got.ID != tests[i].threadID {
			t.Errorf("got=[%s], want=[%s]", got.ID, tests[i].threadID)
		}

		if got.NativeID != tests[i].nativeID {
			t.Errorf("got=[%s], want=[%s]", got.NativeID, tests[i].nativeID)
		}

		if got.State != tests[i].state {
			t.Errorf("got=[%s], want=[%s]", got.State, tests[i].state)
		}
	}

}

func TestCheckThreadInfo2(t *testing.T) {

	threads, err := ParseFromFile("threaddumpsamples/14.0.1.1-together.txt")
	if err != nil {
		t.Error(err)
	}

	const expectedNumberOfThreadsInSampleFile = 9
	if len(threads) != expectedNumberOfThreadsInSampleFile {
		t.Errorf("got=[%d] threads, expected=[%d]", len(threads), expectedNumberOfThreadsInSampleFile)
	}

	type testCase struct {
		threadName string
		isDaemon   bool
		threadID   string
		nativeID   string
		state      string
		stackTrace string
	}

	test := []testCase{
		testCase{
			threadName: `Reference Handler`,
			isDaemon:   true,
			threadID:   `0x00007f195c29c000`,
			nativeID:   `0xb07f`,
			state:      `RUNNABLE`,
			stackTrace: `at java.lang.ref.Reference.waitForReferencePendingList(java.base@14.0.1/Native Method)
at java.lang.ref.Reference.processPendingReferences(java.base@14.0.1/Reference.java:241)
at java.lang.ref.Reference$ReferenceHandler.run(java.base@14.0.1/Reference.java:213)
`,
		},

		testCase{
			threadName: `VM Thread`,
			isDaemon:   false,
			threadID:   `0x00007f195c299000`,
			nativeID:   `0xb07e`,
			state:      `runnable`,
			stackTrace: ``,
		},

		testCase{
			threadName: `GC Thread#7`,
			isDaemon:   false,
			threadID:   `0x00007f191400a800`,
			nativeID:   `0xb091`,
			state:      `runnable`,
			stackTrace: ``,
		},

		testCase{
			threadName: `G1 Main Marker`,
			isDaemon:   false,
			threadID:   `0x00007f195c08c000`,
			nativeID:   `0xb07a`,
			state:      `runnable`,
			stackTrace: ``,
		},

		testCase{
			threadName: `G1 Conc#0`,
			isDaemon:   false,
			threadID:   `0x00007f195c08d800`,
			nativeID:   `0xb07b`,
			state:      `runnable`,
			stackTrace: ``,
		},

		testCase{
			threadName: `G1 Conc#1`,
			isDaemon:   false,
			threadID:   `0x00007f1924001000`,
			nativeID:   `0xb097`,
			state:      `runnable`,
			stackTrace: ``,
		},

		testCase{
			threadName: `G1 Refine#0`,
			isDaemon:   false,
			threadID:   `0x00007f195c20c000`,
			nativeID:   `0xb07c`,
			state:      `runnable`,
			stackTrace: ``,
		},

		testCase{
			threadName: `G1 Young RemSet Sampling`,
			isDaemon:   false,
			threadID:   `0x00007f195c20d800`,
			nativeID:   `0xb07d`,
			state:      `runnable`,
			stackTrace: ``,
		},

		testCase{
			threadName: `VM Periodic Task Thread`,
			isDaemon:   false,
			threadID:   `0x00007f195c30f000`,
			nativeID:   `0xb087`,
			state:      `waiting on condition`,
			stackTrace: ``,
		},
	}

	for i := 0; i < expectedNumberOfThreadsInSampleFile; i++ {
		got := threads[i]
		if got.Name != test[i].threadName {
			t.Errorf("got=[%s], want=[%s]", got.Name, test[i].threadName)
		}

		if got.Daemon != test[i].isDaemon {
			t.Errorf("got=[%t], want=[%t]", got.Daemon, test[i].isDaemon)
		}

		if got.ID != test[i].threadID {
			t.Errorf("got=[%s], want=[%s]", got.ID, test[i].threadID)
		}

		if got.NativeID != test[i].nativeID {
			t.Errorf("got=[%s], want=[%s]", got.NativeID, test[i].nativeID)
		}

		if got.State != test[i].state {
			t.Errorf("got=[%s], want=[%s]", got.State, test[i].state)
		}

		if got.StackTrace != test[i].stackTrace {
			t.Errorf("got=[\n{%s}\n], want=[\n{%s}\n]", got.StackTrace, test[i].stackTrace)
		}
	}

}

func TestAwaitingNotification(t *testing.T) {
	//
}
