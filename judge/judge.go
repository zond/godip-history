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
    dislodgeds:            make(map[Province]Unit),
    supplyCenters:         make(map[Province]Nationality),
    errors:                make(map[Province]error),
  }
}

type Judge struct {
  orders                map[Province]Adjudicator
  units                 map[Province]Unit
  dislodgeds            map[Province]Unit
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

func (self *Judge) SetDislodgeds(dislodgeds map[Province]Unit) *Judge {
  self.dislodgeds = make(map[Province]Unit)
  for prov, unit := range dislodgeds {
    self.SetDislodged(prov, unit)
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
  fmt.Fprintln(buf, "Dislodgeds", self.dislodgeds)
  fmt.Fprintln(buf, "Phase", self.phase)
  fmt.Fprintln(buf, "Orders", self.orders)
  fmt.Fprintln(buf, "Errors", self.errors)
  return string(buf.Bytes())
}

func (self *Judge) Resolver() *resolver {
  return &resolver{
    Judge:   self,
    visited: make(map[Province]bool),
    guesses: make(map[Province]error),
  }
}

func (self *Judge) Errors() map[Province]error {
  return self.errors
}

func (self *Judge) SupplyCenters() map[Province]Nationality {
  return self.supplyCenters
}

func (self *Judge) SetDislodged(prov Province, unit Unit) {
  if found := self.Dislodged(prov); found != nil {
    panic(fmt.Errorf("%v is already at %v", found, prov))
  }
  self.dislodgeds[prov] = unit
}

func (self *Judge) SetUnit(prov Province, unit Unit) {
  if found := self.Unit(prov); found != nil {
    panic(fmt.Errorf("%v is already at %v", found, prov))
  }
  self.units[prov] = unit
}

func (self *Judge) SetOrder(prov Province, order Adjudicator) {
  if found := self.findOrder(prov); found != nil {
    panic(fmt.Errorf("%v is already at %v", found, prov))
  }
  self.orders[prov] = order
}

func (self *Judge) RemoveDislodged(prov Province) {
  if _, p, ok := self.findDislodged(prov); ok {
    delete(self.dislodgeds, p)
  }
}

func (self *Judge) findDislodged(prov Province) (u Unit, p Province, ok bool) {
  if u, ok = self.dislodgeds[prov]; ok {
    p = prov
    return
  }
  sup, _ := prov.Split()
  if u, ok = self.dislodgeds[sup]; ok {
    p = sup
    return
  }
  for name, _ := range self.graph.Coasts(prov) {
    if u, ok = self.dislodgeds[name]; ok {
      p = name
      return
    }
  }
  return
}

func (self *Judge) Dislodged(prov Province) *Unit {
  if u, _, ok := self.findDislodged(prov); ok {
    return &u
  }
  return nil
}

func (self *Judge) findUnit(prov Province) (u Unit, p Province, ok bool) {
  if u, ok = self.units[prov]; ok {
    p = prov
    return
  }
  sup, _ := prov.Split()
  if u, ok = self.units[sup]; ok {
    p = sup
    return
  }
  for name, _ := range self.graph.Coasts(prov) {
    if u, ok = self.units[name]; ok {
      p = name
      return
    }
  }
  return
}

func (self *Judge) Unit(prov Province) *Unit {
  if u, _, ok := self.findUnit(prov); ok {
    return &u
  }
  return nil
}

func (self *Judge) findOrder(prov Province) (result Order) {
  var ok bool
  if result, ok = self.orders[prov]; ok {
    return
  }
  sup, _ := prov.Split()
  if result, ok = self.orders[sup]; ok {
    return
  }
  for name, _ := range self.graph.Coasts(prov) {
    if result, ok = self.orders[name]; ok {
      return
    }
  }
  return
}

func (self *Judge) Order(prov Province) (result Order) {
  result = self.findOrder(prov)
  if result == nil {
    if unit := self.Unit(prov); unit != nil {
      result = self.defaultOrderGenerator(prov)
    }
  }
  return
}

func (self *Judge) Move(src, dst Province) {
  unit := self.Unit(src)
  if unit == nil {
    panic(fmt.Errorf("No unit at %v?", src))
  }
  if dislodged := self.Unit(dst); dislodged != nil {
    delete(self.units, dst)
    self.dislodgeds[dst] = *dislodged
  }
  delete(self.units, src)
  self.units[dst] = *unit
}

func (self *Judge) Graph() Graph {
  return self.graph
}

func (self *Judge) Find(filter StateFilter) (provinces []Province, orders []Order, units []Unit) {
  visitedProvinces := make(map[Province]bool)
  for prov, unit := range self.units {
    visitedProvinces[prov] = true
    order := self.defaultOrderGenerator(prov)
    if ord := self.Order(prov); ord != nil {
      order = ord
    }
    if filter(prov, order, unit) {
      provinces = append(provinces, prov)
      orders = append(orders, order)
      units = append(units, unit)
    }
  }
  for prov, order := range self.orders {
    if !visitedProvinces[prov] {
      if filter(prov, order, Unit{}) {
        provinces = append(provinces, prov)
        orders = append(orders, order)
        units = append(units, Unit{})
      }
    }
  }
  return
}

func (self *Judge) Next() (err error) {
  for prov, order := range self.orders {
    if err := order.Validate(self); err != nil {
      self.errors[prov] = err
      delete(self.orders, prov)
    }
  }
  for prov, _ := range self.orders {
    err := self.Resolver().Resolve(prov)
    if err != nil {
      self.errors[prov] = err
      delete(self.orders, prov)
    }
  }
  for _, order := range self.orders {
    order.Execute(self)
  }
  self.orders = make(map[Province]Adjudicator)
  self.phase = self.phase.Next()
  return
}

func (self *Judge) Phase() Phase {
  return self.phase
}
