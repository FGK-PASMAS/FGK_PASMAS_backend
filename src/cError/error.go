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
    ErrIncludeNotSupported = errors.New("Include not supported")
    ErrFilterNotSupported = errors.New("Filter not supported")

    ErrRealtimeEventCouldNotBeCreated = errors.New("Realtime event could not be created")
)
