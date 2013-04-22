package common

import (
  "fmt"
  "strings"
  "time"
)

func Max(is ...int) (result int) {
  for index, i := range is {
    if index == 0 || i > result {
      result = i
    }
  }
  return
}

func Min(is ...int) (result int) {
  for index, i := range is {
    if index == 0 || i < result {
      result = i
    }
  }
  return
}

type UnitType string

type Nationality string

type OrderType string

type PhaseType string

type Province string

func (self Province) Split() (sup Province, sub Province) {
  split := strings.Split(string(self), "/")
  if len(split) > 0 {
    sup = Province(split[0])
  }
  if len(split) > 1 {
    sub = Province(split[1])
  }
  return
}

func (self Province) Join(n Province) (result Province) {
  if n != "" {
    result = Province(fmt.Sprintf("%v/%v", self, n))
  } else {
    result = self
  }
  return
}

type Unit struct {
  Type        UnitType
  Nationality Nationality
}

type Phase interface {
  Year() int
  Season() string
  Type() PhaseType
  Next() Phase
  Prev() Phase
}

type PathFilter func(n Province, flags map[Flag]bool, sc *Nationality) bool

type Flag string

type Graph interface {
  Has(Province) bool
  Flags(Province) map[Flag]bool
  SC(Province) *Nationality
  Path(src, dst Province, filter PathFilter) []Province
  Coasts(Province) map[Province]bool
}

type Orders []Order

func (self Orders) Less(a, b int) bool {
  return self[a].At().Before(self[b].At())
}

func (self Orders) Swap(a, b int) {
  self[a], self[b] = self[b], self[a]
}

func (self Orders) Len() int {
  return len(self)
}

type Order interface {
  Type() OrderType
  Targets() []Province
  Validate(Validator) error
  Execute(State)
  At() time.Time
}

type Adjudicator interface {
  Order
  Adjudicate(Resolver) error
}

/*
The BackupRule takes a state, a Province causing an inconsistency and set of all Provinces visited while finding the inconsistency, 
and returns whether the Order provided Province ought to succeed.
*/
type BackupRule func(Resolver, Province, map[Province]bool) error

type StateFilter func(n Province, o Order, u Unit) bool

type Validator interface {
  Order(Province) Order
  Unit(Province) *Unit
  Dislodged(Province) *Unit
  Graph() Graph
  Phase() Phase
  SupplyCenters() map[Province]Nationality
  Find(StateFilter) (provinces []Province, orders []Order, units []Unit)
}

type Resolver interface {
  Validator
  Resolve(Province) error
}

type State interface {
  Validator
  Move(Province, Province)
  SetUnit(Province, Unit)
  RemoveDislodged(Province)
}

type OrderGenerator func(prov Province) Order
