package orders

import (
  "fmt"
  cla "github.com/zond/godip/classical/common"
  dip "github.com/zond/godip/common"
  "time"
)

func Support(targets ...dip.Province) *support {
  if len(targets) < 2 || len(targets) > 3 {
    panic(fmt.Errorf("Support orders must either be support Hold with two targets, or support Move with three targets."))
  }
  return &support{
    targets: targets,
  }
}

type support struct {
  targets []dip.Province
}

func (self *support) At() time.Time {
  return time.Now()
}

func (self *support) Type() dip.OrderType {
  return cla.Support
}

func (self *support) Targets() []dip.Province {
  return self.targets
}

func (self *support) Adjudicate(r dip.Resolver) error {
  unit := r.Unit(self.targets[0])
  if len(self.targets) == 3 {
    if victim := r.Unit(self.targets[2]); victim != nil && victim.Nationality == unit.Nationality {
      return cla.ErrIllegalSupportDestinationNationality
    }
  }
  if breaks, _, _ := r.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
    return (o.Type() == cla.Move && // move
      o.Targets()[1] == self.targets[0] && // against us
      (len(self.targets) == 2 || o.Targets()[0] != self.targets[2]) && // not from something we support attacking
      u.Nationality != unit.Nationality && // not friendly
      cla.MovePossible(r, o.Targets()[0], o.Targets()[1], true, true) == nil) // and legal move counting convoy success
  }); len(breaks) > 0 {
    return cla.ErrSupportBroken{breaks[0]}
  }
  return nil
}

func (self *support) Validate(v dip.Validator) error {
  if v.Phase().Type() != cla.Movement {
    return cla.ErrInvalidPhase
  }
  if !v.Graph().Has(self.targets[0]) {
    return cla.ErrInvalidSource
  }
  if !v.Graph().Has(self.targets[1]) {
    return cla.ErrInvalidTarget
  }
  if v.Unit(self.targets[0]) == nil {
    return cla.ErrMissingUnit
  }
  if v.Unit(self.targets[1]) == nil {
    return cla.ErrMissingSupportUnit
  }
  if len(self.targets) == 2 {
    if cla.AnyMovePossible(v, self.targets[0], self.targets[1]) != nil {
      return cla.ErrIllegalSupportPosition
    }
  } else {
    if !v.Graph().Has(self.targets[2]) {
      return cla.ErrInvalidTarget
    }
    if cla.AnyMovePossible(v, self.targets[0], self.targets[2]) != nil {
      return cla.ErrIllegalSupportDestination
    }
    if cla.AnyMovePossible(v, self.targets[1], self.targets[2]) != nil {
      return cla.ErrInvalidSupportMove
    }
  }
  return nil
}

func (self *support) Execute(state dip.State) {
}
