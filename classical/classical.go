package classical

import (
  "fmt"
  . "github.com/zond/godip/classical/common"
  "github.com/zond/godip/classical/orders"
  "github.com/zond/godip/classical/start"
  "github.com/zond/godip/common"
  "github.com/zond/godip/judge"
)

type phase struct {
  year   int
  season string
  typ    common.PhaseType
}

func (self *phase) Year() int {
  return self.year
}

func (self *phase) Season() string {
  return self.season
}

func (self *phase) Type() common.PhaseType {
  return self.typ
}

func (self *phase) Prev() common.Phase {
  if self.typ == Retreat {
    return &phase{
      year:   self.year,
      season: self.season,
      typ:    Movement,
    }
  } else if self.typ == Movement {
    if self.season == Spring {
      if self.year == 1901 {
        return nil
      }
      return &phase{
        year:   self.year - 1,
        season: Winter,
        typ:    Build,
      }
    } else {
      return &phase{
        year:   self.year,
        season: Spring,
        typ:    Retreat,
      }
    }
  } else {
    return &phase{
      year:   self.year,
      season: Fall,
      typ:    Retreat,
    }
  }
  return nil
}

func (self *phase) Next() common.Phase {
  if self.typ == Movement {
    return &phase{
      year:   self.year,
      season: self.season,
      typ:    Retreat,
    }
  } else if self.typ == Retreat {
    if self.season == Spring {
      return &phase{
        year:   self.year,
        season: Fall,
        typ:    Movement,
      }
    } else {
      return &phase{
        year:   self.year,
        season: Winter,
        typ:    Build,
      }
    }
  } else {
    return &phase{
      year:   self.year + 1,
      season: Spring,
      typ:    Movement,
    }
  }
  return nil
}

func Blank(phase common.Phase) *judge.Judge {
  return judge.New(start.Graph(), phase, BackupRule, DefaultOrderGenerator)
}

func Start() *judge.Judge {
  return judge.New(start.Graph(), &phase{1901, Spring, Movement}, BackupRule, DefaultOrderGenerator).
    SetUnits(start.Units()).
    SetSupplyCenters(start.SupplyCenters())
}

func DefaultOrderGenerator(prov common.Province) common.Order {
  return orders.Hold(prov)
}

/*
BackupRule will make sets of only Move orders succeed, while orders with at least one Convoy all fail.
Any other alternative will cause a panic.
*/
func BackupRule(resolver common.Resolver, prov common.Province, deps map[common.Province]bool) error {
  only_moves := true
  convoys := false
  for prov, _ := range deps {
    if order := resolver.Order(prov); order != nil {
      if order.Type() != Move {
        only_moves = false
      }
      if order.Type() == Convoy {
        convoys = true
      }
    }
  }

  if only_moves {
    return nil
  }
  if convoys {
    return ErrConvoyParadox
  }
  panic(fmt.Errorf("Unknown circular dependency between %v", deps))
}
