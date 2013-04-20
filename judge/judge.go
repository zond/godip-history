package judge

import (
  "bytes"
  "fmt"
  . "github.com/zond/godip/common"
)

func New(graph Graph, phase Phase, backupRule BackupRule, defaultOrderGenerator OrderGenerator) *Judge {
  return &Judge{
    graph:                 graph,
    phase:                 phase,
    backupRule:            backupRule,
    defaultOrderGenerator: defaultOrderGenerator,
    orders:                make(map[Province]Adjudicator),
    units:                 make(map[Province]Unit),
    dislodged:             make(map[Province]Unit),
    supplyCenters:         make(map[Province]Nationality),
    errors:                make(map[Province]error),
  }
}

type Judge struct {
  orders                map[Province]Adjudicator
  units                 map[Province]Unit
  dislodged             map[Province]Unit
  supplyCenters         map[Province]Nationality
  graph                 Graph
  phase                 Phase
  backupRule            BackupRule
  defaultOrderGenerator OrderGenerator
  errors                map[Province]error
}

func (self *Judge) SetOrders(orders map[Province]Adjudicator) *Judge {
  self.orders = make(map[Province]Adjudicator)
  for prov, order := range orders {
    self.SetOrder(prov, order)
  }
  return self
}

func (self *Judge) SetUnits(units map[Province]Unit) *Judge {
  self.units = make(map[Province]Unit)
  for prov, unit := range units {
    self.SetUnit(prov, unit)
  }
  return self
}

func (self *Judge) SetDislodged(dislodged map[Province]Unit) *Judge {
  self.dislodged = make(map[Province]Unit)
  for prov, unit := range dislodged {
    self.SetDislodge(prov, unit)
  }
  return self
}

func (self *Judge) SetSupplyCenters(supplyCenters map[Province]Nationality) *Judge {
  self.supplyCenters = supplyCenters
  return self
}

func (self *Judge) String() string {
  buf := new(bytes.Buffer)
  fmt.Fprintln(buf, self.graph)
  fmt.Fprintln(buf, "SC", self.supplyCenters)
  fmt.Fprintln(buf, "Units", self.units)
  fmt.Fprintln(buf, "Dislodged", self.dislodged)
  fmt.Fprintln(buf, "Phase", self.phase)
  fmt.Fprintln(buf, "Orders", self.orders)
  fmt.Fprintln(buf, "Errors", self.errors)
  return string(buf.Bytes())
}

func (self *Judge) Resolver() *resolver {
  return &resolver{
    Judge:   self,
    visited: make(map[Province]bool),
    guesses: make(map[Province]bool),
  }
}

func (self *Judge) SetDislodge(prov Province, unit Unit) {
  if found, ok := self.dislodged[prov]; ok {
    panic(fmt.Errorf("%v is already at %v", found, prov))
  }
  self.dislodged[prov] = unit
}

func (self *Judge) SetUnit(prov Province, unit Unit) {
  if found, ok := self.Unit(prov); ok {
    panic(fmt.Errorf("%v is already at %v", found, prov))
  }
  self.units[prov] = unit
}

func (self *Judge) SetOrder(prov Province, order Adjudicator) {
  if found, ok := self.Order(prov); ok {
    panic(fmt.Errorf("%v is already at %v", found, prov))
  }
  self.orders[prov] = order
}

func (self *Judge) Unit(prov Province) (unit Unit, ok bool) {
  if unit, ok = self.units[prov]; ok {
    return
  }
  sup, _ := prov.Split()
  if unit, ok = self.units[sup]; ok {
    return
  }
  for name, _ := range self.graph.Coasts(prov) {
    if unit, ok = self.units[name]; ok {
      return
    }
  }
  return
}

func (self *Judge) Order(prov Province) (order Order, ok bool) {
  if order, ok = self.orders[prov]; ok {
    return
  }
  sup, _ := prov.Split()
  if order, ok = self.orders[sup]; ok {
    return
  }
  for name, _ := range self.graph.Coasts(prov) {
    if order, ok = self.orders[name]; ok {
      return
    }
  }
  if !ok {
    if _, ok := self.Unit(prov); ok {
      order = self.defaultOrderGenerator(prov)
    }
  }
  return
}

func (self *Judge) Move(src, dst Province) {
  unit, ok := self.Unit(src)
  if !ok {
    panic(fmt.Errorf("No unit at %v?", src))
  }
  if dislodged, ok := self.Unit(dst); ok {
    delete(self.units, dst)
    self.dislodged[dst] = dislodged
  }
  delete(self.units, src)
  self.units[dst] = unit
}

func (self *Judge) Graph() Graph {
  return self.graph
}

func (self *Judge) Next() (err error) {
  for prov, order := range self.orders {
    if err := order.Validate(self); err != nil {
      self.errors[prov] = err
      delete(self.orders, prov)
    }
  }
  for prov, _ := range self.orders {
    success, err := self.Resolver().Resolve(prov)
    if err != nil {
      self.errors[prov] = err
    }
    if !success {
      delete(self.orders, prov)
    }
  }
  for _, order := range self.orders {
    order.Execute(self)
  }
  self.orders = make(map[Province]Adjudicator)
  self.phase, err = self.phase.Next()
  return
}

func (self *Judge) Phase() Phase {
  return self.phase
}
