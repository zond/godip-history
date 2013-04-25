package orders

import (
  "fmt"
  cla "github.com/zond/godip/classical/common"
  dip "github.com/zond/godip/common"
  "time"
)

func Build(source dip.Province, typ dip.UnitType, at time.Time) *build {
  return &build{
    targets: []dip.Province{source},
    typ:     typ,
    at:      at,
  }
}

type build struct {
  targets []dip.Province
  typ     dip.UnitType
  at      time.Time
}

func (self *build) Type() dip.OrderType {
  return cla.Build
}

func (self *build) String() string {
  return fmt.Sprintf("%v %v %v", self.targets[0], cla.Build, self.typ)
}

func (self *build) Targets() []dip.Province {
  return self.targets
}

func (self *build) At() time.Time {
  return self.at
}

func (self *build) Adjudicate(r dip.Resolver) error {
  me := r.Graph().SC(self.targets[0])
  builds, _, _ := cla.AdjustmentStatus(r, *me)
  if self.at.After(builds[len(builds)-1].At()) {
    return cla.ErrIllegalBuild
  }
  return nil
}

func (self *build) Validate(v dip.Validator) error {
  if v.Phase().Type() != cla.Adjustment {
    return cla.ErrInvalidPhase
  }
  if _, _, ok := v.Unit(self.targets[0]); ok {
    return cla.ErrOccupiedSupplyCenter
  }
  var me dip.Nation
  var ok bool
  if me, self.targets[0], ok = v.SupplyCenter(self.targets[0]); !ok {
    return cla.ErrMissingSupplyCenter
  }
  if owner := v.Graph().SC(self.targets[0]); owner == nil {
    return cla.ErrMissingSupplyCenter
  } else if *owner != me {
    return cla.ErrHostileSupplyCenter
  }
  if _, _, balance := cla.AdjustmentStatus(v, me); balance < 1 {
    return cla.ErrMissingSurplus
  }
  if self.typ == cla.Army && !v.Graph().Flags(self.targets[0])[cla.Land] {
    return cla.ErrIllegalUnitType
  }
  if self.typ == cla.Fleet && !v.Graph().Flags(self.targets[0])[cla.Sea] {
    return cla.ErrIllegalUnitType
  }
  return nil
}

func (self *build) Execute(state dip.State) {
  var me dip.Nation
  for prov, nat := range state.SupplyCenters() {
    if prov == self.targets[0] {
      me = nat
    }
  }
  state.SetUnit(self.targets[0], dip.Unit{self.typ, me})
}
