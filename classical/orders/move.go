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

func (self *move) calcAttackSupport(r dip.Resolver, src, dst dip.Province) int {
  _, supports, _ := r.Find(func(p dip.Province, o dip.Order, u dip.Unit) bool {
    if o.Type() == cla.Support && len(o.Targets()) == 3 && o.Targets()[1] == src && o.Targets()[2] == dst {
      if success, _ := r.Resolve(p); success {
        return true
      }
    }
    return false
  })
  return len(supports)
}

func (self *move) calcHoldSupport(r dip.Resolver, p dip.Province) int {
  _, supports, _ := r.Find(func(p dip.Province, o dip.Order, u dip.Unit) bool {
    if o.Type() == cla.Support && len(o.Targets()) == 2 && o.Targets()[1] == p {
      if success, _ := r.Resolve(p); success {
        return true
      }
    }
    return false
  })
  return len(supports)
}

func (self *move) Adjudicate(r dip.Resolver) (result bool, err error) {
  // my power
  attackStrength := self.calcAttackSupport(r, self.targets[0], self.targets[1]) + 1

  // competing moves for the same destination
  _, competingOrders, _ := r.Find(func(p dip.Province, o dip.Order, u dip.Unit) bool {
    return p != self.targets[0] && o.Type() == cla.Move && o.Targets()[1] == self.targets[1]
  })
  for _, competingOrder := range competingOrders {
    if self.calcAttackSupport(r, competingOrder.Targets()[0], competingOrder.Targets()[1])+1 >= attackStrength {
      return false, cla.ErrBounce{competingOrder.Targets()[0]}
    }
  }

  convoyed := false
  if unit, _ := r.Unit(self.targets[0]); unit.Type == cla.Army {
    if _, steps := r.Graph().Path(self.targets[0], self.targets[1], nil); len(steps) > 1 {
      convoyed = true
    }
  }

  if convoyed {
    if found, _ := r.Graph().Path(self.targets[0], self.targets[1], func(name dip.Province, flags map[dip.Flag]bool, sc *dip.Nationality) bool {
      if unit, ok := r.Unit(name); ok && unit.Type == cla.Fleet {
        if order, ok := r.Order(name); ok && order.Type() == cla.Convoy && order.Targets()[1] == self.targets[0] && order.Targets()[2] == self.targets[1] {
          if success, _ := r.Resolve(name); success {
            return true
          }
        }
      }
      return false
    }); !found {
      return false, cla.ErrMissingConvoy
    }
  }

  if atDest, ok := r.Order(self.targets[1]); ok {
    if !convoyed && atDest.Type() == cla.Move && atDest.Targets()[1] == self.targets[0] { // head to head
      if self.calcAttackSupport(r, atDest.Targets()[0], atDest.Targets()[1])+1 >= attackStrength {
        return false, cla.ErrBounce{self.targets[1]}
      }
    } else if atDest.Type() == cla.Move {
      if success, _ := r.Resolve(self.targets[1]); !success && 1 >= attackStrength { // attack against something that moves away
        return false, cla.ErrBounce{self.targets[1]}
      }
    } else if self.calcHoldSupport(r, self.targets[1])+1 >= attackStrength { // simple attack
      return false, cla.ErrBounce{self.targets[1]}
    }
  }
  return true, nil
}

func (self *move) Validate(v dip.Validator) error {
  if v.Phase().Type() != cla.Movement {
    return cla.ErrInvalidPhase
  }
  if !v.Graph().Has(self.targets[0]) {
    return cla.ErrInvalidSource
  }
  if !v.Graph().Has(self.targets[1]) {
    return cla.ErrInvalidDestination
  }
  if unit, ok := v.Unit(self.targets[0]); !ok {
    return cla.ErrMissingUnit
  } else {
    if unit.Type == cla.Army {
      if !v.Graph().Flags(self.targets[1])[cla.Land] {
        return cla.ErrIllegalDestination
      }
    } else if unit.Type == cla.Fleet {
      if !v.Graph().Flags(self.targets[1])[cla.Sea] {
        return cla.ErrIllegalDestination
      }
    } else {
      panic(fmt.Errorf("Unknown unit type %v", unit.Type))
    }
  }
  found, path := v.Graph().Path(self.targets[0], self.targets[1], nil)
  if !found {
    return cla.ErrMissingPath
  }
  if len(path) > 1 {
    if unit, _ := v.Unit(self.targets[0]); unit.Type == cla.Army {
      if found, _ = v.Graph().Path(self.targets[0], self.targets[1], func(name dip.Province, flags map[dip.Flag]bool, sc *dip.Nationality) bool {
        return name == self.targets[0] || name == self.targets[1] || !flags[cla.Land]
      }); !found {
        return cla.ErrMissingSeaPath
      }
      if found, _ = v.Graph().Path(self.targets[0], self.targets[1], func(name dip.Province, flags map[dip.Flag]bool, sc *dip.Nationality) bool {
        unit, ok := v.Unit(name)
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
