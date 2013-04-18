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

type Flag string

type Resolution int

const (
  Success Resolution = iota
  Failure
)

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
