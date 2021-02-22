package jthreadparser

const (
	// LockedState ...
	LockedState SynchronizerState = 0
	// ParkingToWaitForState ...
	ParkingToWaitForState SynchronizerState = 1
	// WaitingOnState ...
	WaitingOnState SynchronizerState = 2
	// WaitingToLockState ...
	WaitingToLockState SynchronizerState = 3
)

const (
	threadInformationBegins   = "\""
	threadHeaderInfoRgx       = `^\"(.*).*tid=(\w*) nid=(\w*)\s\w*`
	threadNameRgx             = `^"(.*)".*\stid=(\w*) nid=(\w*)\s\w*`
	threadNameWithPriorityRgx = `^"(.*)".*\sprio=(\d+).*\stid=(\w*) nid=(\w*)\s\w*`
	stateRgx                  = `\s*java.lang.Thread.State: (.*)`
	runnableStateRgx          = `runnable\s{1,2}$`
	// SMRInfoRgx ...
	SMRInfoRgx                  = `^Threads class SMR info:$`
	waitingOnStateRgx           = `waiting on condition\s{1,2}$`
	synchronizerRgx             = `\s*\- (waiting on|parking to wait for|locked|waiting to lock)\s*<(.*)>\s*\(a\s(.*)\)`
	stackTraceRgx               = `^\s+(at|\-\s).*\)$`
	stackTraceRgxMethodName     = `at\s+(.*)$`
	threadNameRgxGroupIndex     = 1
	threadPriorityRgxGroupIndex = 2
	threadIDRgxGroupIndex       = 3
	threadNativeIDRgxGroupIndex = 4
)
