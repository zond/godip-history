package judge

import (
  . "github.com/zond/godip/common"
)

type resolver struct {
  *Judge
  visited map[Province]bool
  guesses map[Province]bool
}

/*
Will recursively visit the Order dependency graph by calling adjudicate for its Order.

If the recursion hits the same Order again, the result will be guessed twice, once for each outcome.
If only one of the guesses was consistent with the Order#Adjudicate result, that result will be returned, 
otherwise a BackupRule will be invoced.

Make sure never to call Order#Adjudicate from another Order! Only call Resolver#Resolve from Orders.
*/
func (self *resolver) Resolve(prov Province) (result bool, err error) {
  if result, ok := self.guesses[prov]; !ok { // Already guessed
    if self.visited[prov] { // Not yet guessed, but visited before, introduce a guess (the default false result) and return it.
      self.guesses[prov] = result
    } else { // Not yet visited, do a proper adjudication.
      self.visited[prov] = true

      result, err = self.Judge.orders[prov].Adjudicate(self) // Ask order to adjudicate itself.
      if _, ok := self.guesses[prov]; ok {                   // We were visited again, and depend on our guess.
        self.guesses[prov] = true // Switch the guess to true.
        var second_result bool
        second_result, err = self.Judge.orders[prov].Adjudicate(self) // Ask order to adjudicate itself with the new guess.
        if result != second_result {                                  // If the results are the same, it means that exactly one of them were consistent (and any one of them could be returned). If not, none or both are consistent.
          result, err = self.Judge.backupRule(self, prov, self.visited) // So, run the BackupRule on the orders we visited and let it decide.
        }
      }
    }
  }
  return
}
