package classical

import (
	"bytes"
	"fmt"
	cla "github.com/zond/godip/classical/common"
	"github.com/zond/godip/classical/orders"
	dip "github.com/zond/godip/common"
	"sort"
	"strings"
)

type phase struct {
	PhaseYear   int
	PhaseSeason dip.Season
	PhaseType   dip.PhaseType
}

func (self *phase) shortestDistance(s dip.State, src dip.Province, dst []dip.Province) (result int) {
	var unit dip.Unit
	var ok bool
	unit, src, ok = s.Unit(src)
	if !ok {
		panic(fmt.Errorf("No unit at %v", src))
	}
	var filter dip.PathFilter
	found := false
	for _, destination := range dst {
		if unit.Type == cla.Fleet {
			filter = func(p dip.Province, edgeFlags, nodeFlags map[dip.Flag]bool, sc *dip.Nation) bool {
				return edgeFlags[cla.Sea] && nodeFlags[cla.Sea]
			}
		} else {
			filter = func(p dip.Province, edgeFlags, nodeFlags map[dip.Flag]bool, sc *dip.Nation) bool {
				if p.Super() == destination.Super() {
					return true
				}
				u, _, ok := s.Unit(p)
				return (edgeFlags[cla.Land] && nodeFlags[cla.Land]) || (ok && !nodeFlags[cla.Land] && u.Nation == unit.Nation && u.Type == cla.Fleet)
			}
		}
		for _, coast := range s.Graph().Coasts(destination) {
			for _, srcCoast := range s.Graph().Coasts(src) {
				if srcCoast == destination {
					result = 0
					found = true
				} else {
					if path := s.Graph().Path(srcCoast, coast, filter); path != nil {
						if !found || len(path) < result {
							result = len(path)
							found = true
						}
					}
					if path := s.Graph().Path(srcCoast, coast, nil); path != nil {
						if !found || len(path) < result {
							result = len(path)
							found = true
						}
					}
				}
			}
		}
	}
	return
}

type remoteUnitSlice struct {
	provinces []dip.Province
	distances map[dip.Province]int
	units     map[dip.Province]dip.Unit
}

func (self remoteUnitSlice) String() string {
	var l []string
	for _, prov := range self.provinces {
		l = append(l, fmt.Sprintf("%v:%v", prov, self.distances[prov]))
	}
	return strings.Join(l, ", ")
}

func (self remoteUnitSlice) Len() int {
	return len(self.provinces)
}

func (self remoteUnitSlice) Swap(i, j int) {
	self.provinces[i], self.provinces[j] = self.provinces[j], self.provinces[i]
}

func (self remoteUnitSlice) Less(i, j int) bool {
	if self.distances[self.provinces[i]] == self.distances[self.provinces[j]] {
		u1 := self.units[self.provinces[i]]
		u2 := self.units[self.provinces[j]]
		if u1.Type == cla.Fleet && u2.Type == cla.Army {
			return true
		}
		if u2.Type == cla.Fleet && u1.Type == cla.Army {
			return false
		}
		return bytes.Compare([]byte(self.provinces[i]), []byte(self.provinces[j])) < 0
	}
	return self.distances[self.provinces[i]] > self.distances[self.provinces[j]]
}

func (self *phase) sortedUnits(s dip.State, n dip.Nation) []dip.Province {
	provs := remoteUnitSlice{
		distances: make(map[dip.Province]int),
		units:     make(map[dip.Province]dip.Unit),
	}
	provs.provinces, _, _ = s.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
		if u != nil && u.Nation == n {
			provs.distances[p] = self.shortestDistance(s, p, s.Graph().SCs(n))
			provs.units[p] = *u
			return true
		}
		return false
	})
	sort.Sort(provs)
	dip.Logf("Sorted units for %v is %v", n, provs)
	return provs.provinces
}

func (self *phase) DefaultOrder(p dip.Province) dip.Adjudicator {
	if self.PhaseType == cla.Movement {
		return orders.Hold(p)
	}
	return nil
}

func (self *phase) PostProcess(s dip.State) {
	if self.PhaseType == cla.Retreat {
		for prov, _ := range s.Dislodgeds() {
			s.RemoveDislodged(prov)
			s.SetResolution(prov, cla.ErrForcedDisband)
		}
		s.ClearDislodgers()
		s.ClearBounces()
		if self.PhaseSeason == cla.Fall {
			s.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
				if u != nil {
					if s.Graph().SC(p) != nil {
						s.SetSC(p.Super(), u.Nation)
					}
				}
				return false
			})
		}
	} else if self.PhaseType == cla.Adjustment {
		for _, nationality := range cla.Nations {
			_, _, balance := cla.AdjustmentStatus(s, nationality)
			if balance < 0 {
				su := self.sortedUnits(s, nationality)[:-balance]
				for _, prov := range su {
					dip.Logf("Removing %v due to forced disband", prov)
					s.RemoveUnit(prov)
					s.SetResolution(prov, cla.ErrForcedDisband)
				}
			}
		}
	} else if self.PhaseType == cla.Movement {
		for prov, unit := range s.Dislodgeds() {
			hasRetreat := false
			for _, edge := range s.Graph().Edges(prov) {
				if _, _, ok := s.Unit(edge); !ok && !s.Bounce(prov, edge) {
					if path := cla.Path(s, unit.Type, prov, edge); len(path) == 1 {
						dip.Logf("%v can retreat to %v", prov, edge)
						hasRetreat = true
						break
					}
				}
			}
			if !hasRetreat {
				s.RemoveDislodged(prov)
				dip.Logf("Removing %v since it has no retreat", prov)
			}
		}
	}
}

func (self *phase) Year() int {
	return self.PhaseYear
}

func (self *phase) Season() dip.Season {
	return self.PhaseSeason
}

func (self *phase) Type() dip.PhaseType {
	return self.PhaseType
}

func (self *phase) Prev() dip.Phase {
	if self.PhaseType == cla.Retreat {
		return &phase{
			PhaseYear:   self.PhaseYear,
			PhaseSeason: self.PhaseSeason,
			PhaseType:   cla.Movement,
		}
	} else if self.PhaseType == cla.Movement {
		if self.PhaseSeason == cla.Spring {
			if self.PhaseYear == 1901 {
				return nil
			}
			return &phase{
				PhaseYear:   self.PhaseYear - 1,
				PhaseSeason: cla.Fall,
				PhaseType:   cla.Adjustment,
			}
		} else {
			return &phase{
				PhaseYear:   self.PhaseYear,
				PhaseSeason: cla.Spring,
				PhaseType:   cla.Retreat,
			}
		}
	} else {
		return &phase{
			PhaseYear:   self.PhaseYear,
			PhaseSeason: cla.Fall,
			PhaseType:   cla.Retreat,
		}
	}
	return nil
}

func (self *phase) Next() dip.Phase {
	if self.PhaseType == cla.Movement {
		return &phase{
			PhaseYear:   self.PhaseYear,
			PhaseSeason: self.PhaseSeason,
			PhaseType:   cla.Retreat,
		}
	} else if self.PhaseType == cla.Retreat {
		if self.PhaseSeason == cla.Spring {
			return &phase{
				PhaseYear:   self.PhaseYear,
				PhaseSeason: cla.Fall,
				PhaseType:   cla.Movement,
			}
		} else {
			return &phase{
				PhaseYear:   self.PhaseYear,
				PhaseSeason: cla.Fall,
				PhaseType:   cla.Adjustment,
			}
		}
	} else {
		return &phase{
			PhaseYear:   self.PhaseYear + 1,
			PhaseSeason: cla.Spring,
			PhaseType:   cla.Movement,
		}
	}
	return nil
}
