package judge

import (
  "bytes"
  "fmt"
  . "github.com/zond/godip/common"
)

type ErrorCode int

const (
  NoError ErrorCode = iota
  ErrTargetLength
  ErrInvalidTarget
  ErrInvalidPhase
  ErrMissingUnit
)

/* The resolution of an order can be in three states. */
type orderState int

const (
  /* Order is not yet resolved, the resolution has no meaningful value. */
  unresolved orderState = iota
  /* The resolution contains a value, but it is only a guess. */
  guessing
  /* The resolution contains a value, and is final. */
  resolved
)

type Order interface {
  Type() OrderType
  Targets() []Province
  Adjudicate(*State) bool
  Validate(*State) (bool, ErrorCode)
  Execute(*State)
}

/*
The BackupRule takes a state and a slice of Provinces, and returns the resolutions for the orders for the given provinces.
*/
type BackupRule func(state *State, prov Province, deps map[Province]bool) bool

type State struct {
  Orders        map[Province]Order
  Units         map[Province]Unit
  Dislodged     map[Province]Unit
  SupplyCenters map[Province]Nationality
  Graph         Graph
  Phase         Phase

  BackupRule BackupRule

  visited map[Province]bool
  guesses map[Province]bool
}

func (self *State) String() string {
  buf := new(bytes.Buffer)
  fmt.Fprintln(buf, self.Graph)
  fmt.Fprintln(buf, self.SupplyCenters)
  fmt.Fprintln(buf, self.Units)
  fmt.Fprintln(buf, self.Phase)
  fmt.Fprintln(buf, self.Orders)
  return string(buf.Bytes())
}

func (self *State) Next() *State {
  return self
}

func (self *State) Resolve(prov Province) (result bool) {
  if result, ok := self.guesses[prov]; !ok {
    if self.visited[prov] {
      self.guesses[prov] = result
    } else {
      self.visited[prov] = true

      result = self.Orders[prov].Adjudicate(self)
      if _, ok := self.guesses[prov]; ok {
        self.guesses[prov] = true
        second_result := self.Orders[prov].Adjudicate(self)
        if result != second_result {
          result = self.BackupRule(self, prov, self.visited)
        }
      }
    }
  }
  return
}
