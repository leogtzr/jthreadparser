package jthreadparser

const (
	LockedState           SynchronizerState = 0
	ParkingToWaitForState SynchronizerState = 1
	WaitingToLockState    SynchronizerState = 2
)
