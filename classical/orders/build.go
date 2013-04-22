package orders

import (
  cla "github.com/zond/godip/classical/common"
  dip "github.com/zond/godip/common"
  "sort"
  "time"
)

func Build(source dip.Province, typ dip.UnitType, at time.Time) *build {
  return &build{
    targets:   []dip.Province{source},
    typ:       typ,
    createdAt: at,
  }
}

type builds []*build

func (self builds) Less(a, b int) bool {
  return self[a].createdAt.Before(self[b].createdAt)
}

func (self builds) Swap(a, b int) {
  self[a], self[b] = self[b], self[a]
}

func (self builds) Len() int {
  return len(self)
}

type build struct {
  targets   []dip.Province
  typ       dip.UnitType
  createdAt time.Time
}

func (self *build) Type() dip.OrderType {
  return cla.Build
}

func (self *build) Targets() []dip.Province {
  return self.targets
}

func (self *build) Adjudicate(r dip.Resolver) error {
  me, ok := r.SupplyCenters()[self.targets[0]]
  if !ok {
    return cla.ErrIllegalBuild
  }

  scs := 0
  for _, nat := range r.SupplyCenters() {
    if nat == me {
      scs += 1
    }
  }

  units := 0
  var buildOrders builds
  r.Find(func(p dip.Province, o dip.Order, u dip.Unit) bool {
    if u.Nationality == me {
      units += 1
    }
    if r.SupplyCenters()[p] == me {
      if buildOrder, ok := o.(*build); ok {
        buildOrders = append(buildOrders, buildOrder)
      }
    }
    return false
  })
  sort.Sort(buildOrders)

  allowed := scs - units
  buildOrders = buildOrders[:allowed]

  if !self.createdAt.After(buildOrders[len(buildOrders)-1].createdAt) {
    return cla.ErrIllegalBuild
  }

  return nil
}

func (self *build) Validate(v dip.Validator) error {
  if v.Phase().Type() != cla.Build {
    return cla.ErrInvalidPhase
  }
  if v.Unit(self.targets[0]) != nil {
    return cla.ErrOccupiedSupplyCenter
  }
  return nil
}

func (self *build) Execute(state dip.State) {
  var me dip.Nationality
  for prov, nat := range state.SupplyCenters() {
    if prov == self.targets[0] {
      me = nat
    }
  }
  state.SetUnit(self.targets[0], dip.Unit{self.typ, me})
}
