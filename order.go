package godip

type OrderType int

const (
  Hold OrderType = iota
  Move
  Support
  Convoy
  Build
  Disband
)

type OrderFlags int

const (
  ViaConvoyOrder OrderFlags = 1 << iota
)

type Orders []*Order

func (self *Orders) hasConvoyPath(unitType UnitType, from, to string, w *World, deps *Orders) bool {
  if from == to {
    return true
  }
  for typ, conns := range w.slots[from].connections {
    if typ.canConvoy(unitType) {
      for _, slot := range conns {
	if slot.unit != nil && unit.Type.canConvoy(unitType) && slot.canConvoy(unitType) {
	  if order := self.getByOrigin(slot.name); order != nil && order.Type == Convoy && order.Dest1 == from && order.Dest2 == to && order.resolve(w, self, deps
	}
      }
    }
  }
}

func (self *Orders) getByDest1(pos string) (result Orders) {
  for _, order := range self {
    if order.Dest1 == pos {
      result = append(result, order)
    }
  }
  return
}

func (self *Orders) getByDest2(pos string) (result Orders) {
  for _, order := range self {
    if order.Dest2 == pos {
      result = append(result, order)
    }
  }
  return
}

func (self *Orders) getByOrigin(pos string) (result *Order) {
  for _, order := range self {
    if order.Origin == pos {
      if result != nil {
        panic(fmt.Errorf("Only one order per origin is allowed"))
      }
      result = order
    }
  }
  return
}

// Used only on dependencies, which are also order slices
func (self *Orders) backupRule(old_deps int) {
  only_moves := true
  convoys := false
  for i := old_deps; i < len(self); i++ {
    if self[i].Type != Move {
      only_moves = false
    }
    if self[i].Type == Convoy {
      convoys = true
    }
  }
  var popped *Order
  if only_moves {
    for len(self) > old_deps {
      popped = self.Pop()
      popped.Resolution = Success
      popped.state = resolved
    }
    return
  }
  if convoys {
    for len(self) > old_deps {
      popped = self.Pop()
      if popped.Type == Convoy {
        popped.Resolution = Fails
        popped.state = resolved
      }
    }
    return
  }
  panic(fmt.Errorf("Unknown circular dependency"))
}

func (self *Orders) Add(order *Order) {
  for _, o := range *self {
    if o == order {
      return
    }
  }
  *self = append(*self, order)
}

func (self *Orders) Pop() (result *Order) {
  result = (*self)[len(self)-1]
  *self = (*self)[:len(self)-1]
  return
}

func (self *Orders) Resolve(w *World) {
  for _, order := range self {
    order.state = unresolved
  }
  for _, order := range self {
    order.resolve(w, self, &Orders{})
  }
}

type Order struct {
  Origin string
  Dest1  string
  Dest2  string
  Type   OrderType
  Flags  OrderFlags
  // For each order we maintain the resolution.
  Resolution Resolution
  Reason     ResolutionReason
  state      state
}

/*
 Function: resolve(nr)
 nr - The number of the order to be resolved.
 Returns the resolution for that order.
*/
func (self *Order) resolve(w *World, orders *Orders, deps *Orders) (Resolution, ResolutionReason) {
  /* 
     If order is already resolved, just return
     the resolution.
  */
  if self.state == resolved {
    return self.Resolution, self.ResolutionReason
  } else if self.state == guessing {
    /* 
       Order is in guess state. Add the order
       nr to the dependencies list, if it isn't
       there yet and return the guess.
    */
    deps.Add(self)
    return self.Resolution, self.ResolutionReason
  }

  /* 
     Remember how big the dependency list is before we
     enter recursion.
  */
  old_deps := len(deps)

  // Set order in guess state.
  self.Resolution = Fails
  self.state = guessing

  // Adjudicate order.    
  first_result, first_reason := self.adjudicate(w, orders, deps)

  if len(deps) == old_deps {
    /* 
       No orders were added to the dependency list.
       This means that the result is not dependent
       on a guess.
    */

    /* Set the resolution (ignoring the initial
       guess). The order may already have the state
       RESOLVED, due to the backup rule, acting
       outside the cycle.
    */
    if self.state != resolved {
      self.Resolution, self.ResolutionReason, self.state = first_result, first_reason, resolved
    }
    return first_result, first_reason
  }

  if deps[old_deps] != self {
    /* 
       The order is dependent on a guess, but not our
       own guess, because it would be the first
       dependency added. Add to dependency list,
       update result, but state remains guessing
    */
    deps.Add(self)
    self.Resolution, self.ResolutionReason = first_result, first_result
    return first_result, first_reason
  }

  /*
     Result is dependent on our own guess. Set all
     orders in dependency list to UNRESOLVED and reset
     dependency list.
  */
  for len(deps) > old_deps {
    deps.Pop().state = unresolved
  }

  // Do the other guess.
  self.Resolution = success
  self.state = guessing

  // Adjudicate with the other guess.
  second_result, second_reson = self.adjudicate(w, orders, deps)

  if first_result == second_result {
    /*
       Although there is a cycle, there is only
       one resolution. Cleanup dependency list first.
    */
    for len(deps) > old_deps {
      deps.Pop().state = unresolved
    }
    // Now set the final result and return.
    self.Resolution, self.ResolutionReason, self.state = first_result, first_reason, resolved
    return first_result, first_reason
  }

  /*
     There are two or no resolutions for the cycle.
      Pass dependencies to the backup rule.
      These are dependencies with index in range
      [old_nr_of_dep, nr_of_dep - 1]
      The backup_rule, should clean up the dependency
      list (setting nr_of_dep to old_nr_of_dep). Any
      order in the dependency list that is not set to
      RESOLVED should be set to UNRESOLVED.
  */
  deps.backupRule(old_deps)

  /*
     The backup_rule may not have resolved all
     orders in the cycle. For instance, the
     Szykman rule, will not resolve the orders
     of the moves attacking the convoys. To deal
     with this, we start all over again.
  */
  return self.resolve(w, orders, deps)
}

func (self *Order) adjudicate(w *World, orders *Orders, deps *Orders) (Resolution, ResolutionReason) {
  switch self.Type {
  case Move:
    return self.adjudicateMove(w, orders, deps)
  case Support:
    return self.adjudicateSupport(w, orders, deps)
  case Convoy:
    return self.adjudicateConvoy(w, orders, deps)
  case Hold:
    return self.adjudicateHold(w, orders, deps)
  case Build:
    return self.adjudicateBuild(w, orders, deps)
  case Disband:
    return self.adjudicateDisband(w, orders, deps)
  }
  panic(fmt.Errorf("Unknown order type for %+v", self))
}

func (self *Order) attackStrength(w *World, orders *Orders, deps *Orders) (result int) {
  result = 1
  orderAtDest := orders.getByOrigin(self.Dest1)
  unitAtDest := w.getUnitAt(self.Dest1)
  myUnit := w.getUnitAt(self.Origin)
  if unitAtDest == nil || (unitAtDest != nil && orderAtDest.Type == Move && orderAtDest.Dest1 != self.Origin && orderAtDest.resolve(w, orders, deps) == Succeeds) {
    for _, order := range orders.getByDest2(self.Dest1) {
      if order.Type == Support && order.Dest1 == self.Origin && order.resolve(w, orders, deps) == Succeeds {
        result += 1
      }
    }
  } else {
    if unitAtDest.Nationality == myUnit.Nationality {
      result = 0
    } else {
      for _, order := range orders.getByDest2(self.Dest1) {
        if order.Type == Support && order.Dest1 == self.Origin && order.resolve(w, orders, deps) == Succeeds {
          result += 1
        }
      }
    }
  }
  return
}

func (self *Order) hasMovePath(w *World, orders *Orders, deps *Orders) bool {
  return w.hasPath(self.Type, self.Origin, self.Dest1)
}

func (self *Order) hasConvoyPath(w *World, orders *Orders, deps *Orders) bool {
  return orders.hasConvoyPath(w, deps, self.Origin, self.Dest1)
}

func (self *Order) adjudicateMove(w *World, orders *Orders, deps *Orders) (Resolution, ResolutionReason) {
  var competingStrengths []int

  if !self.hasMovePath(w, orders, deps) {
    if self.hasConvoyPath(w, orders, deps) {
      competingStrengths = append(competingStrengths, w.holdStrength(self.Dest1, orders, deps))
    } else {
      return Failed, NoConvoyReason
    }
  }

  if self.Flags&ViaConvoyOrder == ViaConvoyOrder && self.hasConvoyPath(w, orders, deps) {
    competingStrenghts = append(competingStrenghts, w.holdStrenght(self.Dest1, orders, deps))
  } else {

    var headToHead *Order
    if prospect := orders.getByOrigin(self.Dest1); prospect != nil && prospect.Type == Move && prospect.Dest1 == self.Origin {
      headToHead = prospect
    }

    competingStrengths = append(competingStrenghts, headToHead.defendStrength(w, orders, deps))

  }
}

func (self *Order) adjudicateSupport(w *World, orders *Orders, deps *Orders) (Resolution, ResolutionReason) {
}

func (self *Order) adjudicateConvoy(w *World, orders *Orders, deps *Orders) (Resolution, ResolutionReason) {
}

func (self *Order) adjudicateHold(w *World, orders *Orders, deps *Orders) (Resolution, ResolutionReason) {
}

func (self *Order) adjudicateBuild(w *World, orders *Orders, deps *Orders) (Resolution, ResolutionReason) {
}

func (self *Order) adjudicateDisband(w *World, orders *Orders, deps *Orders) (Resolution, ResolutionReason) {
}
