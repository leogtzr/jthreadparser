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

type SynchronizerState int
