package cerror

import (
	"errors"

	"gorm.io/gorm"
)

var (
    ErrPassengerWeightIsZero = errors.New("Passenger weight is zero")
    ErrObjectDependencyMissing = errors.New("Object dependency missing")
    ErrObjectNotFound = gorm.ErrRecordNotFound
    ErrNoRealtimeHandlerFound = errors.New("No realtime handler found")
)
