package orders

import (
  cla "github.com/zond/godip/classical/common"
  dip "github.com/zond/godip/common"
  "time"
)

func Retreat(src, dst dip.Province) *retreat {
  return &retreat{
    targets: []dip.Province{src, dst},
  }
}

type retreat struct {
  targets []dip.Province
}

func (self *retreat) Type() dip.OrderType {
  return cla.Retreat
}

func (self *retreat) Targets() []dip.Province {
  return self.targets
}

func (self *retreat) At() time.Time {
  return time.Now()
}

func (self *retreat) Adjudicate(r dip.Resolver) error {
  _, competingOrders, _ := r.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
    return p != self.targets[0] && o.Type() == cla.Retreat && o.Targets()[1] == self.targets[1]
  })
  if len(competingOrders) > 0 {
    return cla.ErrBounce{competingOrders[0].Targets()[0]}
  }
  return nil
}

func (self *retreat) Validate(v dip.Validator) error {
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
  return nil
}

func (self *retreat) Execute(state dip.State) {
  state.Retreat(self.targets[0], self.targets[1])
}
