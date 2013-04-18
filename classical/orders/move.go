package orders

import (
  "fmt"
  cla "github.com/zond/godip/classical/common"
  dip "github.com/zond/godip/common"
  "github.com/zond/godip/judge"
)

func Move(source, dest dip.Province) *move {
  return &move{
    targets: []dip.Province{source, dest},
  }
}

type move struct {
  targets []dip.Province
}

func (self *move) Type() dip.OrderType {
  return cla.Move
}

func (self *move) Targets() []dip.Province {
  return self.targets
}

func (self *move) Adjudicate(state *judge.State) bool {
  // if head to head: defend strength of h2h < attack strength
  // else: hold strength of target < attack strength
  return false
}

func (self *move) Validate(state *judge.State) error {
  if state.Phase.Type() != cla.Movement {
    return cla.ErrInvalidPhase
  }
  if len(self.targets) != 2 {
    return cla.ErrTargetLength
  }
  if !state.Graph.Has(self.targets[0]) {
    return cla.ErrInvalidSource
  }
  if !state.Graph.Has(self.targets[1]) {
    return cla.ErrInvalidDestination
  }
  if unit, ok := state.Units[self.targets[0]]; !ok {
    return cla.ErrMissingUnit
  } else {
    if unit.Type == cla.Army {
      if !state.Graph.Flags(self.targets[1])[cla.Land] {
        return cla.ErrIllegalDestination
      }
    } else if unit.Type == cla.Fleet {
      if !state.Graph.Flags(self.targets[1])[cla.Sea] {
        return cla.ErrIllegalDestination
      }
    } else {
      panic(fmt.Errorf("Unknown unit type %v", unit.Type))
    }
  }
  if !state.Graph.Edges(self.targets[0])[self.targets[1]] {
    return cla.ErrIllegalDistance
  }
  // support convoys
  return nil
}

func (self *move) Execute(state *judge.State) {
  unit := state.Units[self.targets[0]]
  delete(state.Units, self.targets[0])
  if dislodged, ok := state.Units[self.targets[1]]; ok {
    delete(state.Units, self.targets[1])
    state.Dislodged[self.targets[1]] = dislodged
  }
  state.Units[self.targets[1]] = unit
}
