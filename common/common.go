package common

import (
  "fmt"
  "strings"
)

type UnitType string

type Nationality string

type OrderType string

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

type ProvinceFlag int

type Resolution int

const (
  Success Resolution = iota
  Failure
)

type Unit struct {
  Type        UnitType
  Nationality Nationality
}

type Order interface {
  Type() OrderType
  Targets() []Province
  Adjudicate(State) Resolution
}

type State interface {
  Resolve(orders []Order) State
  Orders() map[Province]Order
  Units() map[Province]Unit
  SupplyCenters() map[Province]Nationality
}
