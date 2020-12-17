package main

import (
	"fmt"
)

type Vessel interface {
	Storage() *Storage
	DeltaV() int64 // in m/s
}

// Storage holds values for determining the cargo capacity of a Vessel
type Storage struct {
	Volume int64 // in m3
	Tonnage int64 // in metric tons
	// Tonnage, as a measure of cargo-carrying capacity is actually a measure of volume
	// Here, we're going to use it as weight because the in-industry term for this,
	// displacement, doesn't make sense for spaceships
}

// Station is a stationary Vessel that produces and stores cargo
type Station struct {
}

// Ship is a mobile Vessel that holds and transports cargo
type Ship struct {
}
