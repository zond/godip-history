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

func (self *move) calcAttackSupport(r dip.Resolver, src, dst dip.Province, forbidden dip.Nation) int {
  _, supports, _ := r.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
    if o != nil && u != nil && o.Type() == cla.Support && len(o.Targets()) == 3 && u.Nation != forbidden && o.Targets()[1].Contains(src) && o.Targets()[2].Contains(dst) {
      if err := r.Resolve(p); err == nil {
        return true
      }
    }
    return false
  })
  return len(supports)
}

func (self *move) calcHoldSupport(r dip.Resolver, prov dip.Province) int {
  _, supports, _ := r.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
    if o != nil && u != nil && o.Type() == cla.Support && p.Super() != prov.Super() && len(o.Targets()) == 2 && o.Targets()[1].Super() == prov.Super() {
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
  unit, _, _ := r.Unit(self.targets[0])
  // competing moves for the same destination
  _, competingOrders, competingUnits := r.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
    return o != nil && u != nil && o.Type() == cla.Move && o.Targets()[0] != self.targets[0] && self.targets[1].Super() == o.Targets()[1].Super()
  })
  for index, competingOrder := range competingOrders {
    attackStrength := self.calcAttackSupport(r, self.targets[0], self.targets[1], competingUnits[index].Nation) + 1
    dip.Logf("%v:vs %v: %v", self, competingOrder, attackStrength)
    if as := self.calcAttackSupport(r, competingOrder.Targets()[0], competingOrder.Targets()[1], unit.Nation) + 1; as >= attackStrength {
      dip.Indent(fmt.Sprintf("D(%v):", competingOrder.Targets()[0]))
      if dislodgers, _, _ := r.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
        res := o != nil && // is an order
          u != nil && // is a unit
          o.Type() == cla.Move && // move
          o.Targets()[1].Super() == competingOrder.Targets()[0].Super() && // against the competition
          o.Targets()[0].Super() == self.targets[1].Super() && // is from our destination
          u.Nation != competingUnits[index].Nation && // not from themselves
          r.Resolve(p) == nil // and it succeeded
        return res
      }); len(dislodgers) == 0 {
        dip.DeIndent()
        dip.Logf("%v:vs %v: %v", competingOrder, self, as)
        return cla.ErrBounce{competingOrder.Targets()[0]}
      } else {
        dip.DeIndent()
      }
    }
  }

  convoyed := false
  if unit.Type == cla.Army {
    steps := r.Graph().Path(self.targets[0], self.targets[1], nil)
    if self.viaConvoy || len(steps) > 1 {
      dip.Indent(fmt.Sprintf("C(%v):", self.targets))
      err := cla.AnyConvoyPossible(r, self.targets[0], self.targets[1], true)
      if err != nil {
        dip.Logf("%v", err)
        dip.DeIndent()
        if len(steps) > 1 {
          return err
        }
      } else {
        dip.Logf("T")
        dip.DeIndent()
        convoyed = true
      }
    }
  }

  if blocker, _, ok := r.Unit(self.targets[1]); ok {
    attackStrength := self.calcAttackSupport(r, self.targets[0], self.targets[1], blocker.Nation) + 1
    order, prov, _ := r.Order(self.targets[1])
    dip.Logf("%v:vs %v: %v", self, order, attackStrength)
    if !convoyed && order.Type() == cla.Move && order.Targets()[1] == self.targets[0] { // head to head
      as := self.calcAttackSupport(r, order.Targets()[0], order.Targets()[1], unit.Nation) + 1
      dip.Logf("%v:vs %v: %v", order, self, as)
      if blocker.Nation == unit.Nation || as >= attackStrength {
        return cla.ErrBounce{self.targets[1]}
      }
    } else if order.Type() == cla.Move { // attack against something that moves away
      if err := r.Resolve(prov); err != nil {
        if blocker.Nation == unit.Nation || 1 >= attackStrength {
          return cla.ErrBounce{self.targets[1]}
        }
      }
    } else { // simple attack
      hs := self.calcHoldSupport(r, self.targets[1]) + 1
      dip.Logf("%v: %v", order, hs)
      if blocker.Nation == unit.Nation || hs >= attackStrength {
        return cla.ErrBounce{self.targets[1]}
      }
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
  if !v.Graph().Has(self.targets[0]) {
    return cla.ErrInvalidSource
  }
  if !v.Graph().Has(self.targets[1]) {
    return cla.ErrInvalidDestination
  }
  var unit dip.Unit
  var ok bool
  if unit, self.targets[0], ok = v.Dislodged(self.targets[0]); !ok {
    return cla.ErrMissingUnit
  }
  var err error
  if self.targets[1], err = cla.AnyMovePossible(v, self.targets[0], self.targets[1], unit.Type == cla.Army, false, false); err != nil {
    return err
  }
  if v.IsDislodger(self.targets[1], self.targets[0]) {
    return cla.ErrIllegalRetreat
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
  var unit dip.Unit
  var ok bool
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
