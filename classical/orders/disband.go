package orders

import (
  "fmt"
  cla "github.com/zond/godip/classical/common"
  dip "github.com/zond/godip/common"
  "time"
)

func Disband(source dip.Province, at time.Time) *disband {
  return &disband{
    targets: []dip.Province{source},
    at:      at,
  }
}

type disband struct {
  targets []dip.Province
  at      time.Time
}

func (self *disband) String() string {
  return fmt.Sprintf("%v %v", self.targets[0], cla.Disband)
}

func (self *disband) Type() dip.OrderType {
  return cla.Disband
}

func (self *disband) Targets() []dip.Province {
  return self.targets
}

func (self *disband) adjudicateBuildPhase(r dip.Resolver) error {
  unit, _, _ := r.Unit(self.targets[0])
  _, disbands, _ := cla.AdjustmentStatus(r, unit.Nation)
  if self.at.After(disbands[len(disbands)-1].At()) {
    return cla.ErrIllegalDisband
  }
  return nil
}

func (self *disband) adjudicateRetreatPhase(r dip.Resolver) error {
  return nil
}

func (self *disband) Adjudicate(r dip.Resolver) error {
  if r.Phase().Type() == cla.Adjustment {
    return self.adjudicateBuildPhase(r)
  }
  return self.adjudicateRetreatPhase(r)
}

func (self *disband) validateRetreatPhase(v dip.Validator) error {
  if _, _, ok := v.Dislodged(self.targets[0]); !ok {
    return cla.ErrMissingUnit
  }
  return nil
}

func (self *disband) validateBuildPhase(v dip.Validator) error {
  unit, _, ok := v.Unit(self.targets[0])
  if !ok {
    return cla.ErrMissingUnit
  }
  if _, _, balance := cla.AdjustmentStatus(v, unit.Nation); balance > -1 {
    return cla.ErrMissingDeficit
  }
  return nil
}

func (self *disband) sanitizeBuildPhase(v dip.Validator) error {
  if !v.Graph().Has(self.targets[0]) {
    return cla.ErrInvalidTarget
  }
  return nil
}

func (self *disband) sanitizeRetreatPhase(v dip.Validator) error {
  if !v.Graph().Has(self.targets[0]) {
    return cla.ErrInvalidTarget
  }
  return nil
}

func (self *disband) Sanitize(v dip.Validator) error {
  if v.Phase().Type() == cla.Adjustment {
    return self.sanitizeBuildPhase(v)
  } else if v.Phase().Type() == cla.Retreat {
    return self.sanitizeRetreatPhase(v)
  }
  return cla.ErrInvalidPhase
}

func (self *disband) Validate(v dip.Validator) error {
  if v.Phase().Type() == cla.Adjustment {
    return self.validateBuildPhase(v)
  }
  return self.validateRetreatPhase(v)
}

func (self *disband) Execute(state dip.State) {
  if state.Phase().Type() == cla.Adjustment {
    state.RemoveUnit(self.targets[0])
  } else {
    state.RemoveDislodged(self.targets[0])
  }
}
