package orders

import (
  "fmt"
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
  targets   []dip.Province
  viaConvoy bool
}

func (self *move) String() string {
  return fmt.Sprintf("%v %v %v", self.targets[0], cla.Move, self.targets[1])
}

func (self *move) ViaConvoy() *move {
  self.viaConvoy = true
  return self
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
    if o.Type() == cla.Support && len(o.Targets()) == 3 && o.Targets()[1].Contains(src) && o.Targets()[2].Contains(dst) {
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
    if o.Type() == cla.Support && len(o.Targets()) == 2 && o.Targets()[1].Contains(p) {
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
    return o.Type() == cla.Move && o.Targets()[0] != self.targets[0] && self.targets[1].Super() == o.Targets()[1].Super()
  })
  for _, competingOrder := range competingOrders {
    if self.calcAttackSupport(r, competingOrder.Targets()[0], competingOrder.Targets()[1])+1 >= attackStrength {
      return cla.ErrBounce{competingOrder.Targets()[0]}
    }
  }

  convoyed := false
  if unit, _, ok := r.Unit(self.targets[0]); ok && unit.Type == cla.Army {
    steps := r.Graph().Path(self.targets[0], self.targets[1], nil)
    if self.viaConvoy || len(steps) > 1 {
      err := cla.AnyConvoyPossible(r, self.targets[0], self.targets[1], true)
      if err != nil {
        if len(steps) > 1 {
          return err
        }
      } else {
        convoyed = true
      }
    }
  }

  if atDest, _, ok := r.Order(self.targets[1]); ok {
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

func (self *move) validateRetreatPhase(v dip.Validator) error {
  if v.Phase().Type() != cla.Retreat {
    return cla.ErrInvalidPhase
  }
  var ok bool
  var unit dip.Unit
  if unit, self.targets[0], ok = v.Dislodged(self.targets[0]); !ok {
    return cla.ErrMissingUnit
  }
  if _, _, ok = v.Unit(self.targets[1]); ok {
    return cla.ErrOccupiedDestination
  }
  if v.IsDislodger(self.targets[1], self.targets[0]) {
    return cla.ErrIllegalRetreat
  }
  var err error
  if self.targets[1], err = cla.AnyMovePossible(v, self.targets[0], self.targets[1], unit.Type == cla.Army, false, false); err != nil {
    return err
  }
  return nil
}

func (self *move) validateMovementPhase(v dip.Validator) error {
  if !v.Graph().Has(self.targets[0]) {
    return cla.ErrInvalidSource
  }
  if !v.Graph().Has(self.targets[1]) {
    return cla.ErrInvalidDestination
  }
  var ok bool
  var unit dip.Unit
  if unit, self.targets[0], ok = v.Unit(self.targets[0]); !ok {
    return cla.ErrMissingUnit
  }
  var err error
  if self.targets[1], err = cla.AnyMovePossible(v, self.targets[0], self.targets[1], unit.Type == cla.Army, true, false); err != nil {
    return err
  }
  return nil
}

func (self *move) Execute(state dip.State) {
  if state.Phase().Type() == cla.Retreat {
    state.Retreat(self.targets[0], self.targets[1])
  } else {
    state.Move(self.targets[0], self.targets[1])
  }
}
