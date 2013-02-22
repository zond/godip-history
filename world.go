package godip

// Possible resolutions of an order. 
type resolution int

const (
  fails resolution = iota
  succeeds
)

// The resolution of an order, can be in three states.
type state int

const (
  unresolved state = iota
  guessing
  resolved
)

type World map[string]Province

type OrderType int

const (
  Hold OrderType = iota
  Move
  Support
  Convoy
  Build
  Disband
)

type Order struct {
  Location    string
  Source      string
  Destination string
  Type        OrderType
  // For each order we maintain the resolution.
  resolution resolution
  state      state
}

/*
 Function: resolve(nr)
 nr - The number of the order to be resolved.
 Returns the resolution for that order.
*/
func (self *Order) resolve(w World, orders *Orders, deps *Orders) resolution {
  /* 
     If order is already resolved, just return
     the resolution.
  */
  if self.state == resolved {
    return self.resolution
  } else if self.state == guessing {
    /* 
       Order is in guess state. Add the order
       nr to the dependencies list, if it isn't
       there yet and return the guess.
    */
    deps.Add(self)
    return self.resolution
  }

  /* 
     Remember how big the dependency list is before we
     enter recursion.
  */
  old_deps := len(deps)

  // Set order in guess state.
  self.resolution = fails
  self.state = guessing

  // Adjudicate order.    
  first_result := self.adjudicate(w, orders, deps)

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
      self.resolution = first_result
      self.state = resolved
    }
    return first_result
  }

  if deps[old_deps] != self {
    /* 
       The order is dependent on a guess, but not our
       own guess, because it would be the first
       dependency added. Add to dependency list,
       update result, but state remains guessing
    */
    deps.Add(self)
    self.resolution = first_result
    return first_result
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
  self.resolution = success
  self.state = guessing

  // Adjudicate with the other guess.
  second_result = self.adjudicate(w, orders, deps)

  if first_result == second_result {
    /*
       Although there is a cycle, there is only
       one resolution. Cleanup dependency list first.
    */
    for len(deps) > old_deps {
      deps.Pop().state = unresolved
    }
    // Now set the final result and return.
    self.resolution = first_result
    self.state = resolved
    return first_result
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
  self.backupRule(old_deps)

  /*
     The backup_rule may not have resolved all
     orders in the cycle. For instance, the
     Szykman rule, will not resolve the orders
     of the moves attacking the convoys. To deal
     with this, we start all over again.
  */
  return self.resolve(w, orders, deps)
}

type Orders []*Order

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

func (self *Orders) Resolve(w World) {
  for _, order := range self {
    order.state = unresolved
  }
  for _, order := range self {
    order.resolve(w, self, &Orders{})
  }
}
