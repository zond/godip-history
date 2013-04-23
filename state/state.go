package state

import (
  "bytes"
  "fmt"
  "github.com/zond/godip/common"
)

func New(graph common.Graph, phase common.Phase, backupRule common.BackupRule, defaultOrderGenerator common.OrderGenerator) *State {
  return &State{
    graph:                 graph,
    phase:                 phase,
    backupRule:            backupRule,
    defaultOrderGenerator: defaultOrderGenerator,
    orders:                make(map[common.Province]common.Adjudicator),
    units:                 make(map[common.Province]common.Unit),
    dislodgeds:            make(map[common.Province]common.Unit),
    supplyCenters:         make(map[common.Province]common.Nationality),
    errors:                make(map[common.Province]error),
    dislodgers:            make(map[common.Province]common.Province),
  }
}

type State struct {
  orders                map[common.Province]common.Adjudicator
  units                 map[common.Province]common.Unit
  dislodgeds            map[common.Province]common.Unit
  supplyCenters         map[common.Province]common.Nationality
  graph                 common.Graph
  phase                 common.Phase
  backupRule            common.BackupRule
  defaultOrderGenerator common.OrderGenerator
  errors                map[common.Province]error
  dislodgers            map[common.Province]common.Province
}

func (self *State) SetOrders(orders map[common.Province]common.Adjudicator) *State {
  self.orders = make(map[common.Province]common.Adjudicator)
  for prov, order := range orders {
    self.SetOrder(prov, order)
  }
  return self
}

func (self *State) ClearDislodgers() {
  self.dislodgers = make(map[common.Province]common.Province)
}

func (self *State) SetError(p common.Province, e error) {
  self.errors[p] = e
}

func (self *State) SetSC(p common.Province, n common.Nationality) {
  self.supplyCenters[p] = n
}

func (self *State) SetUnits(units map[common.Province]common.Unit) *State {
  self.units = make(map[common.Province]common.Unit)
  for prov, unit := range units {
    self.SetUnit(prov, unit)
  }
  return self
}

func (self *State) SetDislodgeds(dislodgeds map[common.Province]common.Unit) *State {
  self.dislodgeds = make(map[common.Province]common.Unit)
  for prov, unit := range dislodgeds {
    self.SetDislodged(prov, unit)
  }
  return self
}

func (self *State) SetSupplyCenters(supplyCenters map[common.Province]common.Nationality) *State {
  self.supplyCenters = supplyCenters
  return self
}

func (self *State) String() string {
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

func (self *State) Resolver() *resolver {
  return &resolver{
    State:   self,
    visited: make(map[common.Province]bool),
    guesses: make(map[common.Province]error),
  }
}

func (self *State) Errors() map[common.Province]error {
  return self.errors
}

func (self *State) SupplyCenters() map[common.Province]common.Nationality {
  return self.supplyCenters
}

func (self *State) SetDislodged(prov common.Province, unit common.Unit) {
  if found := self.Dislodged(prov); found != nil {
    panic(fmt.Errorf("%v is already at %v", found, prov))
  }
  self.dislodgeds[prov] = unit
}

func (self *State) SetUnit(prov common.Province, unit common.Unit) {
  if found := self.Unit(prov); found != nil {
    panic(fmt.Errorf("%v is already at %v", found, prov))
  }
  self.units[prov] = unit
}

func (self *State) SetOrder(prov common.Province, order common.Adjudicator) {
  if found := self.findOrder(prov); found != nil {
    panic(fmt.Errorf("%v is already at %v", found, prov))
  }
  self.orders[prov] = order
}

func (self *State) RemoveUnit(prov common.Province) {
  if _, p, ok := self.findUnit(prov); ok {
    delete(self.units, p)
  }
}

func (self *State) RemoveDislodged(prov common.Province) {
  if _, p, ok := self.findDislodged(prov); ok {
    delete(self.dislodgeds, p)
  }
}

func (self *State) findDislodged(prov common.Province) (u common.Unit, p common.Province, ok bool) {
  if u, ok = self.dislodgeds[prov]; ok {
    p = prov
    return
  }
  sup, _ := prov.Split()
  if u, ok = self.dislodgeds[sup]; ok {
    p = sup
    return
  }
  for _, name := range self.graph.Coasts(prov) {
    if u, ok = self.dislodgeds[name]; ok {
      p = name
      return
    }
  }
  return
}

func (self *State) Dislodged(prov common.Province) *common.Unit {
  if u, _, ok := self.findDislodged(prov); ok {
    return &u
  }
  return nil
}

func (self *State) findDislodger(prov common.Province) (p common.Province, ok bool) {
  if p, ok = self.dislodgers[prov]; ok {
    return
  }
  sup, _ := prov.Split()
  if p, ok = self.dislodgers[sup]; ok {
    return
  }
  for _, name := range self.graph.Coasts(prov) {
    if p, ok = self.dislodgers[name]; ok {
      return
    }
  }
  return
}

func (self *State) IsDislodger(attacker, victim common.Province) bool {
  if dislodger, ok := self.findDislodger(victim); ok {
    if dislodger == attacker {
      return true
    }
    sup, _ := dislodger.Split()
    if sup == attacker {
      return true
    }
    for _, name := range self.graph.Coasts(dislodger) {
      if name == attacker {
        return true
      }
    }
  }
  return false
}

func (self *State) findUnit(prov common.Province) (u common.Unit, p common.Province, ok bool) {
  if u, ok = self.units[prov]; ok {
    p = prov
    return
  }
  sup, _ := prov.Split()
  if u, ok = self.units[sup]; ok {
    p = sup
    return
  }
  for _, name := range self.graph.Coasts(prov) {
    if u, ok = self.units[name]; ok {
      p = name
      return
    }
  }
  return
}

func (self *State) Unit(prov common.Province) *common.Unit {
  if u, _, ok := self.findUnit(prov); ok {
    return &u
  }
  return nil
}

func (self *State) findOrder(prov common.Province) (result common.Order) {
  var ok bool
  if result, ok = self.orders[prov]; ok {
    return
  }
  sup, _ := prov.Split()
  if result, ok = self.orders[sup]; ok {
    return
  }
  for _, name := range self.graph.Coasts(prov) {
    if result, ok = self.orders[name]; ok {
      return
    }
  }
  return
}

func (self *State) Order(prov common.Province) (result common.Order) {
  result = self.findOrder(prov)
  if result == nil {
    if unit := self.Unit(prov); unit != nil {
      result = self.defaultOrderGenerator(prov)
    }
  }
  return
}

func (self *State) Move(src, dst common.Province) {
  if unit, prov, ok := self.findUnit(src); !ok {
    panic(fmt.Errorf("No unit at %v?", src))
  } else {
    if d, p, ok := self.findDislodged(dst); ok {
      self.RemoveUnit(p)
      self.SetDislodged(p, d)
      self.dislodgers[dst] = prov
    }
    self.RemoveUnit(prov)
    self.SetUnit(dst, unit)
  }
}

func (self *State) Retreat(src, dst common.Province) {
  if unit, prov, ok := self.findDislodged(src); !ok {
    panic(fmt.Errorf("No dislodged at %v?", src))
  } else {
    self.RemoveDislodged(prov)
    self.SetUnit(dst, unit)
  }
}

func (self *State) Graph() common.Graph {
  return self.graph
}

func (self *State) Find(filter common.StateFilter) (provinces []common.Province, orders []common.Order, units []common.Unit) {
  visitedProvinces := make(map[common.Province]bool)
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
      if filter(prov, order, common.Unit{}) {
        provinces = append(provinces, prov)
        orders = append(orders, order)
        units = append(units, common.Unit{})
      }
    }
  }
  return
}

func (self *State) Next() (err error) {
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
  for prov, order := range self.orders {
    order.Execute(self)
    delete(self.orders, prov)
  }
  self.phase.PostProcess(self)
  self.phase = self.phase.Next()
  return
}

func (self *State) Phase() common.Phase {
  return self.phase
}
