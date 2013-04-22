package orders

import (
  cla "github.com/zond/godip/classical/common"
  dip "github.com/zond/godip/common"
)

func Hold(source dip.Province) *hold {
  return &hold{
    targets: []dip.Province{source},
  }
}

type hold struct {
  targets []dip.Province
}

func (self *hold) Type() dip.OrderType {
  return cla.Hold
}

func (self *hold) Targets() []dip.Province {
  return self.targets
}

func (self *hold) Adjudicate(r dip.Resolver) (result bool, err error) {
  return true, nil
}

func (self *hold) Validate(v dip.Validator) error {
  if v.Phase().Type() != cla.Movement {
    return cla.ErrInvalidPhase
  }
  return nil
}

func (self *hold) Execute(state dip.State) {
}
