package judge

import (
  "bytes"
  "fmt"
  . "github.com/zond/godip/common"
  "github.com/zond/godip/graph"
)

type Resolution int

const (
  Success Resolution = iota
  Failure
)

type Order interface {
  Type() OrderType
  Targets() []Province
  Adjudicate(State) Resolution
}

type State struct {
  Orders        map[Province]Order
  Units         map[Province]Unit
  SupplyCenters map[Province]Nationality
  Graph         *graph.Graph
  Phase         Phase
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

func (self *State) Resolve() *State {
  return self
}
