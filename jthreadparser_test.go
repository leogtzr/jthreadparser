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
	th := extractThreadInfoFromLine(threadInfo)
	if th.Name != "Attach Listener" {
		t.Error("Expected 'Attach Listener, got ", th.Name)
	}
	if th.ID != "0x00002aaab74c5000" {
		t.Error("Expected '0x00002aaab74c5000', got ", th.ID)
	}
	if th.NativeID != "0x2ea5" {
		t.Error("Expected '0x2ea5', got ", th.NativeID)
	}
}

func TestExtractThreadState(t *testing.T) {
	const expectedState = "TIMED_WAITING"
	state := extractThreadState(threadStateLine)
	if state != expectedState {
		t.Errorf("Expected '%s', got '%s'", expectedState, state)
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
	expectedThreadStringInfo := "Thread Id: '0x00007f90d0106000' (daemon), Name: 'Attach Listener', State: 'RUNNABLE'"

	th := ThreadInfo{ID: "0x00007f90d0106000", Name: "Attach Listener", State: "RUNNABLE", Daemon: true}
	if th.String() != expectedThreadStringInfo {
		t.Errorf("got=[%s], expected=[%s]", th.String(), expectedThreadStringInfo)
	}

	th.Daemon = false
	expectedThreadStringInfo = "Thread Id: '0x00007f90d0106000', Name: 'Attach Listener', State: 'RUNNABLE'"
	if th.String() != expectedThreadStringInfo {
		t.Errorf("got=[%s], expected=[%s]", th.String(), expectedThreadStringInfo)
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
	threads, err := ParseFrom(strings.NewReader(threadInfoWithLocks))
	if err != nil {
		t.Error("Error parsing thread information.")
	}
	for _, thread := range threads {
		holds := HoldsForThread(&thread)
		if len(holds) != expectedNumberOfLocksInThreadWithLockInfo {
			t.Errorf("expected=[%d], got=[%d]", expectedNumberOfLocksInThreadWithLockInfo, len(holds))
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
	threads, err := ParseFromFile("tdump.sample")
	if err != nil {
		t.Error("Unable to parse thread dump sample file")
	}

	mostUsedMethods := MostUsedMethods(&threads)
	javaMethodName := "java.lang.Object.wait(Native Method)"

	expectedNumberOfThreadsWithMethodName := 82

	if mostUsedMethods[javaMethodName] != expectedNumberOfThreadsWithMethodName {
		t.Errorf("Should have identified %d threads, got=%d", expectedNumberOfThreadsWithMethodName, mostUsedMethods[javaMethodName])
	}
}

func TestIdenticalStrackTrace(t *testing.T) {
	expectedNumberOfThreadsWithIdenticalStackTrace := 20
	stacktrace := `at sun.misc.Unsafe.park(Native Method)
- parking to wait for  <0x0000000750785368> (a java.util.concurrent.locks.AbstractQueuedSynchronizer$ConditionObject)
at java.util.concurrent.locks.LockSupport.park(LockSupport.java:156)
at java.util.concurrent.locks.AbstractQueuedSynchronizer$ConditionObject.await(AbstractQueuedSynchronizer.java:1987)
at java.util.concurrent.LinkedBlockingQueue.take(LinkedBlockingQueue.java:399)
at java.util.concurrent.ThreadPoolExecutor.getTask(ThreadPoolExecutor.java:957)
at java.util.concurrent.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:917)
at java.lang.Thread.run(Thread.java:682)`

	threads, err := ParseFromFile("tdump.sample")
	if err != nil {
		t.Error("Unable to parse thread dump sample file")
	}

	identicalStackTrace := IdenticalStackTrace(&threads)
	if identicalStackTrace[stacktrace] != expectedNumberOfThreadsWithIdenticalStackTrace {
		t.Errorf("Should have identified %d (got=%d) threads with the following stacktrace:\n[%s]",
			expectedNumberOfThreadsWithIdenticalStackTrace, identicalStackTrace[stacktrace], stacktrace)
	}

}
