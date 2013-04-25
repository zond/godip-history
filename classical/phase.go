package classical

import (
  "fmt"
  cla "github.com/zond/godip/classical/common"
  dip "github.com/zond/godip/common"
  "sort"
)

type phase struct {
  year   int
  season dip.Season
  typ    dip.PhaseType
}

func (self *phase) shortestDistance(s dip.State, src dip.Province, dst []dip.Province) (result int) {
  var unit dip.Unit
  var ok bool
  unit, src, ok = s.Unit(src)
  if !ok {
    panic(fmt.Errorf("No unit at %v", src))
  }
  var filter dip.PathFilter
  if unit.Type == cla.Fleet {
    filter = func(p dip.Province, edgeFlags, nodeFlags map[dip.Flag]bool, sc *dip.Nation) bool {
      return edgeFlags[cla.Sea] && nodeFlags[cla.Sea]
    }
  } else {
    filter = func(p dip.Province, edgeFlags, nodeFlags map[dip.Flag]bool, sc *dip.Nation) bool {
      u, _, ok := s.Unit(p)
      return (edgeFlags[cla.Land] && nodeFlags[cla.Land]) || (ok && u.Nation == unit.Nation && u.Type == cla.Fleet)
    }
  }
  found := false
  for _, destination := range dst {
    for _, coast := range s.Graph().Coasts(destination) {
      if path := s.Graph().Path(src, coast, filter); path != nil {
        if !found || len(path) < result {
          result = len(path)
          found = true
        }
      }
    }
  }
  return
}

type remoteUnitSlice struct {
  provinces []dip.Province
  distances map[dip.Province]int
}

func (self remoteUnitSlice) Len() int {
  return len(self.provinces)
}

func (self remoteUnitSlice) Swap(i, j int) {
  self.provinces[i], self.provinces[j] = self.provinces[j], self.provinces[i]
}

func (self remoteUnitSlice) Less(i, j int) bool {
  return self.distances[self.provinces[i]] > self.distances[self.provinces[j]]
}

func (self *phase) sortedUnits(s dip.State, n dip.Nation) []dip.Province {
  provs := remoteUnitSlice{
    distances: make(map[dip.Province]int),
  }
  provs.provinces, _, _ = s.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
    if u != nil && u.Nation == n {
      provs.distances[p] = self.shortestDistance(s, p, s.Graph().SCs(n))
      return true
    }
    return false
  })
  sort.Sort(provs)
  return provs.provinces
}

func (self *phase) PostProcess(s dip.State) {
  if self.typ == cla.Retreat {
    s.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
      if _, _, ok := s.Dislodged(p); ok {
        s.RemoveDislodged(p)
        s.SetError(p, cla.ErrForcedDisband)
      }
      return false
    })
    s.ClearDislodgers()
    if self.season == cla.Fall {
      s.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
        if u != nil {
          if s.Graph().SC(p) != nil {
            s.SetSC(p, u.Nation)
          }
        }
        return false
      })
    }
  } else if self.typ == cla.Adjustment {
    for _, nationality := range cla.Nations {
      _, _, balance := cla.AdjustmentStatus(s, nationality)
      if balance < 0 {
        for _, prov := range self.sortedUnits(s, nationality)[:-balance] {
          s.RemoveUnit(prov)
          s.SetError(prov, cla.ErrForcedDisband)
        }
      }
    }
  }
}

func (self *phase) Year() int {
  return self.year
}

func (self *phase) Season() dip.Season {
  return self.season
}

func (self *phase) Type() dip.PhaseType {
  return self.typ
}

func (self *phase) Prev() dip.Phase {
  if self.typ == cla.Retreat {
    return &phase{
      year:   self.year,
      season: self.season,
      typ:    cla.Movement,
    }
  } else if self.typ == cla.Movement {
    if self.season == cla.Spring {
      if self.year == 1901 {
        return nil
      }
      return &phase{
        year:   self.year - 1,
        season: cla.Fall,
        typ:    cla.Adjustment,
      }
    } else {
      return &phase{
        year:   self.year,
        season: cla.Spring,
        typ:    cla.Retreat,
      }
    }
  } else {
    return &phase{
      year:   self.year,
      season: cla.Fall,
      typ:    cla.Retreat,
    }
  }
  return nil
}

func (self *phase) Next() dip.Phase {
  if self.typ == cla.Movement {
    return &phase{
      year:   self.year,
      season: self.season,
      typ:    cla.Retreat,
    }
  } else if self.typ == cla.Retreat {
    if self.season == cla.Spring {
      return &phase{
        year:   self.year,
        season: cla.Fall,
        typ:    cla.Movement,
      }
    } else {
      return &phase{
        year:   self.year,
        season: cla.Fall,
        typ:    cla.Adjustment,
      }
    }
  } else {
    return &phase{
      year:   self.year + 1,
      season: cla.Spring,
      typ:    cla.Movement,
    }
  }
  return nil
}
