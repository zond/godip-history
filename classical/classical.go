package classical

import (
  "bytes"
  "fmt"
  "github.com/zond/godip/classical/start"
  "github.com/zond/godip/common"
  "github.com/zond/godip/graph"
)

type State struct {
  graph         *graph.Graph
  supplyCenters map[string]string
  units         map[string]common.Unit
}

func (self *State) String() string {
  buf := new(bytes.Buffer)
  fmt.Fprintln(buf, self.graph)
  fmt.Fprintln(buf, self.supplyCenters)
  fmt.Fprintln(buf, self.units)
  return string(buf.Bytes())
}

func Blank() *State {
  return &State{
    graph:         start.Graph(),
    supplyCenters: make(map[string]string),
    units:         make(map[string]common.Unit),
  }
}

func Start() (result *State) {
  return &State{
    graph:         start.Graph(),
    supplyCenters: start.SupplyCenters(),
    units:         start.Units(),
  }
}
