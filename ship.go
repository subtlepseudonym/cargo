package main

type Vessel interface {
	Storage() Storage
	DeltaV() int64 // in m/s
}

// Storage holds values for determining the cargo capacity of a Vessel
type Storage struct {
	Volume  int64 // in m3
	Tonnage int64 // in metric tons
	// Tonnage, as a measure of cargo-carrying capacity is actually a measure of volume
	// Here, we're going to use it as weight because the in-industry term for this,
	// displacement, doesn't make sense for spaceships
}

// Station is a stationary Vessel that produces, consumes, and stores cargo
type Station struct {
	Hold Storage
	Produces []Cargo // list of distinct cargo types produced
	Consumes []Cargo // list of distinct cargo types consumed
}

func (s *Station) Storage() Storage {
	return s.Hold
}

func (s *Station) DeltaV() int64 {
	return 0 // stationary
}

// Ship is a mobile Vessel that holds and transports cargo
type Ship struct {
	Hold Storage
}

func (s *Ship) Storage() Storage {
	return s.Hold
}

// Cargo is a type of good
type Cargo struct {
	Volume float64 // TODO: might want to use int64 for simplicity / performance
}
