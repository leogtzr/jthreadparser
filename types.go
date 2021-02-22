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

// ThreadDump ...
type ThreadDump struct {
	Threads []ThreadInfo
	SMR     []string
}
