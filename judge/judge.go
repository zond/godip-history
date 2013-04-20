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
}

func (self *Judge) SetOrders(orders map[Province]Adjudicator) *Judge {
  self.orders = orders
  return self
}

func (self *Judge) SetUnits(units map[Province]Unit) *Judge {
  self.units = units
  return self
}

func (self *Judge) SetDislodged(dislodged map[Province]Unit) *Judge {
  self.dislodged = dislodged
  return self
}

func (self *Judge) SetSupplyCenters(supplyCenters map[Province]Nationality) *Judge {
  self.supplyCenters = supplyCenters
  return self
}

func (self *Judge) String() string {
  buf := new(bytes.Buffer)
  fmt.Fprintln(buf, self.graph)
  fmt.Fprintln(buf, self.supplyCenters)
  fmt.Fprintln(buf, self.units)
  fmt.Fprintln(buf, self.phase)
  fmt.Fprintln(buf, self.orders)
  return string(buf.Bytes())
}

func (self *Judge) Resolver() *resolver {
  return &resolver{
    Judge:   self,
    visited: make(map[Province]bool),
    guesses: make(map[Province]bool),
  }
}

func (self *Judge) SetUnit(prov Province, unit Unit) {
  if found, ok := self.Unit(prov); ok {
    panic(fmt.Errorf("%v is already at %v", found, prov))
  }
  self.units[prov] = unit
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

func (self *Judge) Graph() Graph {
  return self.graph
}

func (self *Judge) Next() (err error) {
  self.phase, err = self.phase.Next()
  return
}

func (self *Judge) Phase() Phase {
  return self.phase
}
