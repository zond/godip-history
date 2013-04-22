package orders

import (
  "fmt"
  cla "github.com/zond/godip/classical/common"
  dip "github.com/zond/godip/common"
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

func (self *support) Type() dip.OrderType {
  return cla.Support
}

func (self *support) Targets() []dip.Province {
  return self.targets
}

func (self *support) Adjudicate(r dip.Resolver) (result bool, err error) {
  return true, nil
}

func (self *support) simpleMoveAllowed(v dip.Validator, src, dst dip.Province) bool {
  unit, _ := v.Unit(src)
  var flag dip.Flag
  if unit.Type == cla.Army {
    flag = cla.Land
  } else {
    flag = cla.Sea
  }
  for coast, _ := range v.Graph().Coasts(dst) {
    if found, steps := v.Graph().Path(src, coast, func(p dip.Province, f map[dip.Flag]bool, sc *dip.Nationality) bool {
      return f[flag]
    }); found && len(steps) == 1 {
      return true
    }
  }
  return false
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
  if _, ok := v.Unit(self.targets[0]); !ok {
    return cla.ErrMissingUnit
  }
  if _, ok := v.Unit(self.targets[1]); !ok {
    return cla.ErrMissingSupportee
  }
  if len(self.targets) == 2 {
    if !self.simpleMoveAllowed(v, self.targets[0], self.targets[1]) {
      return cla.ErrIllegalSupport
    }
  } else {
    if !v.Graph().Has(self.targets[2]) {
      return cla.ErrInvalidTarget
    }
    if !self.simpleMoveAllowed(v, self.targets[0], self.targets[2]) {
      return cla.ErrIllegalSupport
    }
    if Move(self.targets[1], self.targets[2]).Validate(v) != nil {
      return cla.ErrInvalidSupportedMove
    }
  }
  return nil
}

func (self *support) Execute(state dip.State) {
  state.Move(self.targets[0], self.targets[1])
}
