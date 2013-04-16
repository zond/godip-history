package state

import (
  "github.com/zond/godip/common"
)

type State struct {
  orders []common.Order
  graph  common.Graph
  phase  common.Phase
}

func (self State) Phase() common.Phase {
  return self.phase
}

func (self State) Resolve() (State, error) {
  return self, nil
}
