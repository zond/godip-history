package state

import (
	"fmt"
	"github.com/zond/godip/common"
)

type resolver struct {
	*State
	visited map[common.Province]bool
	guesses map[common.Province]error
}

/*
Will recursively visit the Order dependency graph by calling adjudicate for its Order.

If the recursion hits the same Order again, the result will be guessed twice, once for each outcome.
If only one of the guesses was consistent with the Order#Adjudicate result, that result will be returned, 
otherwise a BackupRule will be invoced.

Make sure never to call Order#Adjudicate from another Order! Only call Resolver#Resolve from Orders.
*/
func (self *resolver) Resolve(prov common.Province) (err error) {
	common.Logf("Res(%v)", prov)
	common.Indent("  ")
	var ok bool
	if !self.State.successes[prov] { // Already resolved
		if err, ok = self.State.errors[prov]; !ok { // Already found error
			if err, ok = self.guesses[prov]; !ok { // Already guessed
				if self.visited[prov] { // Not yet guessed, but visited before, introduce a negative guess and return it.
					err = fmt.Errorf("Negative guess")
					common.Logf("Visited, guessing negative")
					self.guesses[prov] = err
				} else { // Not yet visited, do a proper adjudication.
					self.visited[prov] = true

					common.Logf("Adj(%v)", prov)
					common.Indent("  ")
					err = self.State.orders[prov].Adjudicate(self) // Ask order to adjudicate itself.
					common.DeIndent()
					if err != nil {
						common.Logf("%v", err)
					} else {
						common.Logf("T")
					}
					if _, ok := self.guesses[prov]; ok { // We were visited again, and depend on our guess.
						common.Logf("Guess made, switching to positive guess")
						self.guesses[prov] = nil // Switch the guess to success.
						common.Logf("Adj(%v)", prov)
						common.Indent("  ")
						second_err := self.State.orders[prov].Adjudicate(self) // Ask order to adjudicate itself with the new guess.
						common.DeIndent()
						if second_err != nil {
							common.Logf("%v", second_err)
						} else {
							common.Logf("T")
						}
						if (err == nil && second_err != nil) || (err != nil && second_err == nil) { // If the results are the same (in regards to success), it means that exactly one of them were consistent (and any one of them could be returned). If not, none or both are consistent.
							err = self.State.backupRule(self, prov, self.visited) // So, run the BackupRule on the orders we visited and let it decide.
							if err != nil {
								common.Logf("Backup:%v", err)
							} else {
								common.Logf("Backup:T")
							}
						}
					}
					delete(self.visited, prov)
				}
			}
		} else {
			common.Logf("Cached")
		}
	} else {
		common.Logf("Cached")
	}
	common.DeIndent()
	if err != nil {
		common.Logf("%v", err)
	} else {
		common.Logf("T")
	}
	return
}
