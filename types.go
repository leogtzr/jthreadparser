package jthreadparser

// ThreadInfo ...
type ThreadInfo struct {
	Name, ID, NativeID, Priority, State, StackTrace string
	Daemon                                          bool
}

// Locked ...
type Locked struct {
	LockID, LockecObjectName string
}

// SynchronizerState ...
type SynchronizerState int

// Synchronizer ...
type Synchronizer struct {
	ID, ObjectName string
	State          SynchronizerState
}

// SMR section, holds a list of all the non-JVM thread IDs.
type SMR struct {
	IDs []string
}

// ThreadDump ...
type ThreadDump struct {
	Threads []ThreadInfo
	SMR     []SMR
}
