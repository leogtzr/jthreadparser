package jthreadparser

import (
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

func TestExtractThreadDump(t *testing.T) {

}
