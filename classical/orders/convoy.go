package orders

import (
  cla "github.com/zond/godip/classical/common"
  dip "github.com/zond/godip/common"
  "time"
)

func Convoy(source, from, to dip.Province) *convoy {
  return &convoy{
    targets: []dip.Province{source, from, to},
  }
}

type convoy struct {
  targets []dip.Province
}

func (self *convoy) At() time.Time {
  return time.Now()
}

func (self *convoy) Type() dip.OrderType {
  return cla.Convoy
}

func (self *convoy) Targets() []dip.Province {
  return self.targets
}

func (self *convoy) Adjudicate(r dip.Resolver) error {
  unit := r.Unit(self.targets[0])
  if breaks, _, _ := r.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
    return (o.Type() == cla.Move && // move
      o.Targets()[1] == self.targets[0] && // against us
      u.Nationality != unit.Nationality && // not friendly
      r.Resolve(p) == nil)
  }); len(breaks) > 0 {
    return cla.ErrConvoyDislodged{breaks[0]}
  }
  return nil
}

func (self *convoy) Validate(v dip.Validator) error {
  if v.Phase().Type() != cla.Movement {
    return cla.ErrInvalidPhase
  }
  if !v.Graph().Has(self.targets[0]) {
    return cla.ErrInvalidSource
  }
  if !v.Graph().Has(self.targets[1]) {
    return cla.ErrInvalidTarget
  }
  if !v.Graph().Has(self.targets[2]) {
    return cla.ErrInvalidTarget
  }
  if unit := v.Unit(self.targets[0]); unit == nil {
    return cla.ErrMissingUnit
  }
  if unit := v.Unit(self.targets[1]); unit == nil {
    return cla.ErrMissingConvoyee
  }
  if cla.MovePossible(v, self.targets[1], self.targets[2], true, false) != nil {
    return cla.ErrIllegalConvoyMove
  }
  return nil
}

func (self *convoy) Execute(state dip.State) {
}
