package orders

import (
  "fmt"
  cla "github.com/zond/godip/classical/common"
  dip "github.com/zond/godip/common"
)

func Move(source, dest dip.Province) *move {
  return &move{
    targets: []dip.Province{source, dest},
  }
}

type move struct {
  targets []dip.Province
}

func (self *move) Type() dip.OrderType {
  return cla.Move
}

func (self *move) Targets() []dip.Province {
  return self.targets
}

func (self *move) Adjudicate(resolver dip.Resolver) (result bool, err error) {
  // if head to head: defend strength of h2h < attack strength
  // else: hold strength of target < attack strength
  return true, nil
}

func (self *move) Validate(validator dip.Validator) error {
  if validator.Phase().Type() != cla.Movement {
    return cla.ErrInvalidPhase
  }
  if len(self.targets) != 2 {
    return cla.ErrTargetLength
  }
  if !validator.Graph().Has(self.targets[0]) {
    return cla.ErrInvalidSource
  }
  if !validator.Graph().Has(self.targets[1]) {
    return cla.ErrInvalidDestination
  }
  if unit, ok := validator.Unit(self.targets[0]); !ok {
    return cla.ErrMissingUnit
  } else {
    if unit.Type == cla.Army {
      if !validator.Graph().Flags(self.targets[1])[cla.Land] {
        return cla.ErrIllegalDestination
      }
    } else if unit.Type == cla.Fleet {
      if !validator.Graph().Flags(self.targets[1])[cla.Sea] {
        return cla.ErrIllegalDestination
      }
    } else {
      panic(fmt.Errorf("Unknown unit type %v", unit.Type))
    }
  }
  found, path := validator.Graph().Path(self.targets[0], self.targets[1], nil)
  if !found {
    return cla.ErrMissingPath
  }
  if len(path) > 2 {
    if unit, _ := validator.Unit(self.targets[0]); unit.Type == cla.Army {
      if found, _ = validator.Graph().Path(self.targets[0], self.targets[1], func(name dip.Province, flags map[dip.Flag]bool, sc *dip.Nationality) bool {
        return name == self.targets[0] || name == self.targets[1] || !flags[cla.Land]
      }); !found {
        return cla.ErrMissingSeaPath
      }
      if found, _ = validator.Graph().Path(self.targets[0], self.targets[1], func(name dip.Province, flags map[dip.Flag]bool, sc *dip.Nationality) bool {
        unit, ok := validator.Unit(name)
        return ok && unit.Type == cla.Fleet && (name == self.targets[0] || name == self.targets[1] || flags[cla.Sea])
      }); !found {
        return cla.ErrMissingConvoyPath
      }
    } else {
      return cla.ErrIllegalDistance
    }
  }
  return nil
}

func (self *move) Execute(state dip.State) {
  state.Move(self.targets[0], self.targets[1])
}
