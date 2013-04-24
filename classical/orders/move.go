package orders

import (
  cla "github.com/zond/godip/classical/common"
  dip "github.com/zond/godip/common"
  "time"
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

func (self *move) At() time.Time {
  return time.Now()
}

func (self *move) calcAttackSupport(r dip.Resolver, src, dst dip.Province) int {
  _, supports, _ := r.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
    if o.Type() == cla.Support && len(o.Targets()) == 3 && o.Targets()[1] == src && o.Targets()[2] == dst {
      if err := r.Resolve(p); err == nil {
        return true
      }
    }
    return false
  })
  return len(supports)
}

func (self *move) calcHoldSupport(r dip.Resolver, p dip.Province) int {
  _, supports, _ := r.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
    if o.Type() == cla.Support && len(o.Targets()) == 2 && o.Targets()[1] == p {
      if err := r.Resolve(p); err == nil {
        return true
      }
    }
    return false
  })
  return len(supports)
}

func (self *move) Adjudicate(r dip.Resolver) error {
  if r.Phase().Type() == cla.Movement {
    return self.adjudicateMovementPhase(r)
  }
  return self.adjudicateRetreatPhase(r)
}

func (self *move) adjudicateRetreatPhase(r dip.Resolver) error {
  _, competingOrders, _ := r.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
    return p != self.targets[0] && o.Type() == cla.Move && o.Targets()[1] == self.targets[1]
  })
  if len(competingOrders) > 0 {
    return cla.ErrBounce{competingOrders[0].Targets()[0]}
  }
  return nil
}

func (self *move) adjudicateMovementPhase(r dip.Resolver) error {
  // my power
  attackStrength := self.calcAttackSupport(r, self.targets[0], self.targets[1]) + 1

  // competing moves for the same destination
  _, competingOrders, _ := r.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
    return p != self.targets[0] && o.Type() == cla.Move && o.Targets()[1] == self.targets[1]
  })
  for _, competingOrder := range competingOrders {
    if self.calcAttackSupport(r, competingOrder.Targets()[0], competingOrder.Targets()[1])+1 >= attackStrength {
      return cla.ErrBounce{competingOrder.Targets()[0]}
    }
  }

  convoyed := false
  if unit := r.Unit(self.targets[0]); unit.Type == cla.Army {
    if steps := r.Graph().Path(self.targets[0], self.targets[1], nil); len(steps) > 1 {
      convoyed = true
      if err := cla.ConvoyPossible(r, self.targets[0], self.targets[1], true); err != nil {
        return err
      }
    }
  }

  if atDest := r.Order(self.targets[1]); atDest != nil {
    if !convoyed && atDest.Type() == cla.Move && atDest.Targets()[1] == self.targets[0] { // head to head
      if self.calcAttackSupport(r, atDest.Targets()[0], atDest.Targets()[1])+1 >= attackStrength {
        return cla.ErrBounce{self.targets[1]}
      }
    } else if atDest.Type() == cla.Move {
      if err := r.Resolve(self.targets[1]); err != nil && 1 >= attackStrength { // attack against something that moves away
        return cla.ErrBounce{self.targets[1]}
      }
    } else if self.calcHoldSupport(r, self.targets[1])+1 >= attackStrength { // simple attack
      return cla.ErrBounce{self.targets[1]}
    }
  }
  return nil
}

func (self *move) Validate(v dip.Validator) error {
  if v.Phase().Type() == cla.Movement {
    return self.validateMovementPhase(v)
  } else if v.Phase().Type() == cla.Retreat {
    return self.validateRetreatPhase(v)
  }
  return cla.ErrInvalidPhase
}

func (self *move) validateMovePossible(v dip.Validator, useConvoys bool) (err error) {
  if err = cla.MovePossible(v, self.targets[0], self.targets[1], useConvoys, false); err == nil {
    return
  }
  var possibilities []dip.Province
  for _, coast := range v.Graph().Coasts(self.targets[1]) {
    if err2 := cla.MovePossible(v, self.targets[0], coast, useConvoys, false); err2 == nil {
      possibilities = append(possibilities, coast)
    }
  }
  if len(possibilities) == 1 {
    self.targets[1] = possibilities[0]
    err = nil
  }
  return
}

func (self *move) validateRetreatPhase(v dip.Validator) error {
  if v.Phase().Type() != cla.Retreat {
    return cla.ErrInvalidPhase
  }
  if v.Dislodged(self.targets[0]) == nil {
    return cla.ErrMissingUnit
  }
  if v.Unit(self.targets[1]) != nil {
    return cla.ErrOccupiedDestination
  }
  if v.IsDislodger(self.targets[1], self.targets[0]) {
    return cla.ErrIllegalRetreat
  }
  return self.validateMovePossible(v, false)
}

func (self *move) validateMovementPhase(v dip.Validator) error {
  if !v.Graph().Has(self.targets[0]) {
    return cla.ErrInvalidSource
  }
  if !v.Graph().Has(self.targets[1]) {
    return cla.ErrInvalidDestination
  }
  return self.validateMovePossible(v, true)
}

func (self *move) Execute(state dip.State) {
  if state.Phase().Type() == cla.Retreat {
    state.Retreat(self.targets[0], self.targets[1])
  } else {
    state.Move(self.targets[0], self.targets[1])
  }
}
