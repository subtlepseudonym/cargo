package main

import "sync"

type Vessel interface {
	Storage() Storage
	Position() (int64, int64) // TODO: 2D for now
	DeltaV() int64            // in m/s
}

// Storage holds values for determining the cargo capacity of a Vessel
type Storage struct {
	Volume  int64 // in m3
	Tonnage int64 // in metric tons
	// Tonnage, as a measure of cargo-carrying capacity is actually a measure of volume
	// Here, we're going to use it as weight because the in-industry term for this,
	// displacement, doesn't make sense for spaceships
	Contents []Cargo
}

// Station is a stationary Vessel that produces, consumes, and stores cargo
type Station struct {
	Hold     Storage
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
	Volume   float64 // TODO: might want to use int64 for simplicity / performance
	Quantity int64
}

type BlastFurnace struct {
	// TODO: come up with volatiles / impurities for all input materials
	Hematite  int64 // Fe2 O3
	Magnetite int64 // Fe3 O4

	Coke     int64 // 95% pure carbon
	Charcoal int64 // 75% pure carbon

	Oxygen int64 // TODO: how is this imported / refined?

	CarbonDioxide int64
	Iron          int64

	mu *sync.Mutex
}

// Refine starts the iron refining process. This refines as much iron as possible
// from the input materials, throwing away excess intermediate materials
func (b *BlastFurnace) Refine() {
	// take snapshot of buffers for thread safety
	b.mu.Lock()
	hematite := b.Hematite
	magnetite := b.Magnetite
	coke := b.Coke
	charcoal := b.Charcoal
	oxygen := b.Oxygen
	b.mu.Unlock()

	cokeCarbon := int64(float64(coke) * 0.95)
	charcoalCarbon := int64(float64(charcoal) * 0.75)
	carbon := cokeCarbon + charcoalCarbon

	// 2C + O2 -> 2CO
	carbonMonoxideFactor := Min(carbon/2, oxygen/2)
	carbonMonoxide := carbonMonoxideFactor * 2

	usedOxygen := carbonMonoxideFactor * 2

	usedCarbon := Min(cokeCarbon, carbonMonoxideFactor*2)
	usedCoke := int64(float64(usedCarbon) / 0.95)
	// FIXME: may lose some due to rounding errors (is this a problem?)

	remainingCarbon := carbon - usedCarbon
	usedCharcoal := int64(float64(Min(charcoalCarbon, remainingCarbon)) / 0.75)
	// FIXME: may lose some due to rounding errors (is this a problem?)

	// 3(Fe2O3) + CO -> 2(Fe3O4) + C02
	magnetiteFactor := Min(hematite/3, carbonMonoxide)
	carbonMonoxide = carbonMonoxide - magnetiteFactor // consumed
	producedMagnetite := 2 * magnetiteFactor          // produced
	carbonDioxide := magnetiteFactor                  // produced

	usedHematite := magnetiteFactor * 3

	// Fe3O4 + CO -> 3(FeO) + CO2
	ironOxideFactor := Min(producedMagnetite+magnetite, carbonMonoxide)
	carbonMonoxide = carbonMonoxide - ironOxideFactor // consumed
	ironOxide := ironOxideFactor * 3                  // produced
	carbonDioxide = carbonDioxide + ironOxideFactor   // produced

	usedMagnetite := ironOxideFactor

	// FeO + CO -> Fe + CO2
	ironFactor := Min(ironOxide, carbonMonoxide)
	iron := ironFactor
	carbonDioxide = carbonDioxide + ironFactor

	// update buffers
	b.mu.Lock()
	b.Hematite = b.Hematite - usedHematite
	// TODO: add something about magnetite reclamation to blast furnace description
	b.Magnetite = b.Magnetite + producedMagnetite - usedMagnetite
	b.Coke = b.Coke - usedCoke
	b.Charcoal = b.Charcoal - usedCharcoal
	b.Oxygen = b.Oxygen - usedOxygen

	b.Iron = b.Iron + iron
	b.CarbonDioxide = b.CarbonDioxide + carbonDioxide
	b.mu.Unlock()
}

func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}
