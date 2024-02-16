package cerror

import "errors"

var (
    ErrPassengerWeightIsZero = errors.New("Passenger weight is zero")
    ErrObjectDependencyMissing = errors.New("Object dependency missing")
    ErrObjectNotFound = errors.New("Object not found")
)
