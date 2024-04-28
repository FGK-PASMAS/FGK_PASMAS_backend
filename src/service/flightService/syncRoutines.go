package flightService

import "sync"

// flightCreation is used to lock the creation of a flight to ensure data consistency
var flightCreation sync.Mutex

// flightUpdate is used to lock the update of a flight to ensure data consistency
var flightUpdate sync.Mutex


func LockFlightCreation() {
    flightCreation.Lock()
}

func UnlockFlightCreation() {
    flightCreation.Unlock()
}

func LockFlightUpdate() {
    flightUpdate.Lock()
}

func UnlockFlightUpdate() {
    flightUpdate.Unlock()
}

func LockAll() {
    flightCreation.Lock()
    flightUpdate.Lock()
}

func UnlockAll() {
    flightCreation.Unlock()
    flightUpdate.Unlock()
}
