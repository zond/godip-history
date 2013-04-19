package common

import (
  "fmt"
  "strings"
)

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
  Next() (Phase, error)
  Prev() (Phase, error)
}

type PathFilter func(n Province, flags map[Flag]bool, sc *Nationality) bool

type Flag string

type Graph interface {
  Has(Province) bool
  Flags(Province) map[Flag]bool
  SC(Province) *Nationality
  Path(src, dst Province, filter PathFilter) (found bool, path []Province)
  Coasts(Province) map[Province]bool
}

type Order interface {
  Type() OrderType
  Targets() []Province
  Adjudicate(Resolver) (bool, error)
  Validate(Validator) error
  Execute(State)
}

/*
The BackupRule takes a state, a Province causing an inconsistency and set of all Provinces visited while finding the inconsistency, 
and returns whether the Order provided Province ought to succeed.
*/
type BackupRule func(Resolver, Province, map[Province]bool) (bool, error)

type Validator interface {
  Order(Province) (Order, bool)
  Unit(Province) (Unit, bool)
  Graph() Graph
  Phase() Phase
}

type Resolver interface {
  Validator
  Resolve(Province) (bool, error)
}

type State interface {
  Move(Province, Province)
}
