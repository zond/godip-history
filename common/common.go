package common

import (
  "fmt"
  "strconv"
  "strings"
  "time"
)

func MustParseInt(s string) (result int) {
  var err error
  if result, err = strconv.Atoi(s); err != nil {
    panic(err)
  }
  return
}

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

type Nation string

type OrderType string

type PhaseType string

type Province string

type Season string

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

func (self Province) Super() (result Province) {
  result, _ = self.Split()
  return
}

func (self Province) Sub() (result Province) {
  _, result = self.Split()
  return
}

func (self Province) Contains(p Province) bool {
  return self == p || (self.Super() == self && self == p.Super())
}

type Unit struct {
  Type   UnitType
  Nation Nation
}

func (self Unit) Equal(o Unit) bool {
  return self.Type == o.Type && self.Nation == o.Nation
}

type Phase interface {
  Year() int
  Season() Season
  Type() PhaseType
  Next() Phase
  Prev() Phase
  PostProcess(State)
}

type PathFilter func(n Province, edgeFlags, provFlags map[Flag]bool, sc *Nation) bool

type Flag string

type Graph interface {
  Has(Province) bool
  Flags(Province) map[Flag]bool
  SC(Province) *Nation
  Path(src, dst Province, filter PathFilter) []Province
  Coasts(Province) []Province
  SCs(Nation) []Province
  Provinces() []Province
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

type StateFilter func(n Province, o Order, u *Unit) bool

type OrderGenerator func(prov Province) Order

type Validator interface {
  Order(Province) (Order, Province, bool)
  Unit(Province) (Unit, Province, bool)
  Dislodged(Province) (Unit, Province, bool)
  SupplyCenter(Province) (Nation, Province, bool)

  SupplyCenters() map[Province]Nation

  IsDislodger(attacker Province, victim Province) bool
  Graph() Graph
  Phase() Phase
  Find(StateFilter) (provinces []Province, orders []Order, units []Unit)
}

type Resolver interface {
  Validator
  Resolve(Province) error
}

type State interface {
  Validator

  Orders() map[Province]Adjudicator
  Units() map[Province]Unit
  Dislodgeds() map[Province]Unit

  Move(Province, Province)
  Retreat(Province, Province)

  RemoveDislodged(Province)
  RemoveUnit(Province)

  SetError(Province, error)
  SetSC(Province, Nation)
  SetOrder(Province, Adjudicator)
  SetUnit(Province, Unit)

  ClearDislodgers()
}
