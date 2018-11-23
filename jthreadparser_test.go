package jthreadparser

import (
    "testing"
)

const (
    threadInfo = `"Attach Listener" daemon prio=10 tid=0x00002aaab74c5000 nid=0x2ea5 waiting on condition [0x0000000000000000]`
    threadStateLine = `   java.lang.Thread.State: TIMED_WAITING (parking)`
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
