package orders

import (
  cla "github.com/zond/godip/classical/common"
  dip "github.com/zond/godip/common"
  "github.com/zond/godip/judge"
)

type Move struct {
  targets []dip.Province
}

func (self *Move) Type() dip.OrderType {
  return cla.Move
}

func (self *Move) Targets() []dip.Province {
  return self.targets
}

func (self *Move) Adjudicate(state *judge.State) bool {
  // if head to head: defend strength of h2h < attack strength
  // else: hold strength of target < attack strength
  return false
}

func (self *Move) Validate(state *judge.State) (bool, judge.ErrorCode) {
  if state.Phase.Type() != cla.Movement {
    return false, judge.ErrInvalidPhase
  }
  if len(self.targets) != 2 {
    return false, judge.ErrTargetLength
  }
  for _, target := range self.targets {
    if _, _, found := state.Graph.Find(target); !found {
      return false, judge.ErrInvalidTarget
    }
  }
  if _, ok := state.Units[self.targets[0]]; !ok {
    return false, judge.ErrMissingUnit
  }
  // validate that the found unit can move the asked for distance over the found terrain
  return true, 0
}

func (self *Move) Execute(state *judge.State) {
  unit := state.Units[self.targets[0]]
  delete(state.Units, self.targets[0])
  if dislodged, ok := state.Units[self.targets[1]]; ok {
    delete(state.Units, self.targets[1])
    state.Dislodged[self.targets[1]] = dislodged
  }
  state.Units[self.targets[1]] = unit
}
