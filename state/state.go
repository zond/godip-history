package state

import (
  "github.com/zond/godip/common"
)

type State struct {
  units     map[string]common.Unit
  ownership map[string]string
  graph     common.Graph
}

func (self State) Resolve(orders []common.Order) (result State) {
  // check orders valid for node
  // check orders valid for unit
  return self
}
