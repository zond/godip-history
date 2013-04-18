package judge

import (
  "bytes"
  "fmt"
  . "github.com/zond/godip/common"
  "github.com/zond/godip/graph"
)

type ErrorCode int

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
}

/*
The BackupRule takes a state and a slice of Provinces, and returns the resolutions for the orders for the given provinces.
*/
type BackupRule func(state *State, deps []Province) []bool

type State struct {
  Orders        map[Province]Order
  Units         map[Province]Unit
  SupplyCenters map[Province]Nationality
  Graph         *graph.Graph
  Phase         Phase

  BackupRule BackupRule

  /* The resolution of an order can be in three states. */
  states map[Province]orderState
  /* For each order we maintain the resolution. */
  resolutions map[Province]bool
  /* A dependency list is maintained, when a cycle is detected. It is initially empty. */
  dep_list []Province
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

func (self *State) Resolve(prov Province) bool {
  if self.states[prov] == resolved {
    /* If order is already resolved, just return the resolution. */
    return self.resolutions[prov]
  }

  if self.states[prov] == guessing {
    /* Order is in guess state. Add the order nr to the dependencies list, if it isn't there yet and return the guess. */
    for i := 0; i < len(self.dep_list); i++ {
      if self.dep_list[i] == prov {
        /* Already in dependency list, just return last resolution. */
        return self.resolutions[prov]
      }
    }
    /* Add it to dependency list and return last resolution. */
    self.dep_list = append(self.dep_list, prov)
    return self.resolutions[prov]
  }

  /* Remember how big the dependency list is before we enter recursion. */
  old_nr_of_deps := len(self.dep_list)

  /*  Set order in guess state. */
  self.resolutions[prov] = false
  self.states[prov] = guessing

  /*  Adjudicate order. */
  first_result := self.Orders[prov].Adjudicate(self)

  if len(self.dep_list) == old_nr_of_deps {
    /* No orders were added to the dependency list. This means that the result is not dependent on a guess. */
    /* Set the resolution (ignoring the initial guess). The order may already have the state RESOLVED, due to the backup rule, acting outside the cycle. */
    if self.states[prov] != resolved {
      self.resolutions[prov] = first_result
      self.states[prov] = resolved
    }
    return first_result
  }

  if self.dep_list[old_nr_of_deps] != prov {
    /*  The order is dependent on a guess, but not our own guess, because it would be the first dependency added. Add to dependency list, update result, but state remains guessing */
    self.dep_list = append(self.dep_list, prov)
    self.resolutions[prov] = first_result
    return first_result
  }

  /* Result is dependent on our own guess. Set all  orders in dependency list to UNRESOLVED and reset dependency list. */
  for i := old_nr_of_deps; i < len(self.dep_list); i++ {
    self.states[self.dep_list[i]] = unresolved
  }
  self.dep_list = self.dep_list[:old_nr_of_deps]

  /* Do the other guess. */
  self.resolutions[prov] = true
  self.states[prov] = guessing

  /* Adjudicate with the other guess. */
  second_result := self.Orders[prov].Adjudicate(self)

  if first_result == second_result {
    /* Although there is a cycle, there is only one resolution. Cleanup dependency list first. */
    for i := old_nr_of_deps; i < len(self.dep_list); i++ {
      self.states[self.dep_list[i]] = unresolved
    }
    self.dep_list = self.dep_list[:old_nr_of_deps]
    /* Now set the final result and return. */
    self.resolutions[prov] = first_result
    self.states[prov] = resolved
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
  resolutions := self.BackupRule(self, self.dep_list[old_nr_of_deps:])
  for index, resolution := range resolutions {
    name := self.dep_list[old_nr_of_deps+index]
    self.resolutions[name] = resolution
    self.states[name] = resolved
  }
  self.dep_list = self.dep_list[:old_nr_of_deps]

  /*
     The backup_rule may not have resolved all
      orders in the cycle. For instance, the
      Szykman rule, will not resolve the orders
      of the moves attacking the convoys. To deal
      with this, we start all over again.
  */
  return self.Resolve(prov)
}
