package judge

import (
  "bytes"
  "fmt"
  . "github.com/zond/godip/common"
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
  Validate(*State) error
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

/*
Remember to clean up self.visited after calling this for a top level order!
*/
func (self *State) Resolve(prov Province) (result bool) {
  if result, ok := self.guesses[prov]; !ok { // Already guessed
    if self.visited[prov] { // Not yet guessed, but visited before, introduce a guess (the default false result) and return it.
      self.guesses[prov] = result
    } else { // Not yet visited, do a proper adjudication.
      self.visited[prov] = true

      result = self.Orders[prov].Adjudicate(self) // Ask order to adjudicate itself.
      if _, ok := self.guesses[prov]; ok {        // We were visited again, and depend on our guess.
        self.guesses[prov] = true                           // Switch the guess to true.
        second_result := self.Orders[prov].Adjudicate(self) // Ask order to adjudicate itself with the new guess.
        if result != second_result {                        // If the results are the same, it means that exactly one of them were consistent (and any one of them could be returned). If not, none or both are consistent.
          result = self.BackupRule(self, prov, self.visited) // So, run the BackupRule on the orders we visited and let it decide.
        }
      }
    }
  }
  return
}
