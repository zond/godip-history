package state

import (
	"fmt"
	"github.com/zond/godip/common"
)

type resolver struct {
	*State
	path    []common.Province
	guesses map[common.Province]error
}

func (self *resolver) Resolve(prov common.Province) (err error) {

	common.Logf("Res(%v)", prov)
	common.Indent("  ")
	defer func() {
		common.DeIndent()
		if err == nil {
			common.Logf("%v: Success", prov)
		} else {
			common.Logf("%v: Failure: %v", prov, err)
		}
	}()

	var ok bool
	if err, ok = self.State.resolutions[prov]; ok { // Already resolved
		common.Logf("Resolved")
		return
	}

	if err, ok = self.guesses[prov]; ok { // Already guessed
		common.Logf("Guessed")
		for _, ancestor := range self.path { // Add to path if missing
			if ancestor == prov {
				return
			}
		}
		self.path = append(self.path, prov)
		return // Return old guess
	}

	pathLength := len(self.path) // Store old path length

	self.guesses[prov] = fmt.Errorf("Negative guess") // Make negative guess

	common.Logf("Adj(%v)", prov)
	common.Indent("  ")
	err = self.State.orders[prov].Adjudicate(self)
	common.DeIndent()
	if err == nil {
		common.Logf("%v: Success", prov)
	} else {
		common.Logf("%v: Failure: %v", prov, err)
	}

	if len(self.path) == pathLength { // No new guesses made
		delete(self.guesses, prov)
		common.Logf("Resolving, no new guesses made")
		self.State.resolutions[prov] = err // Resolve
		return
	}

	if self.path[pathLength] != prov { // We are not the first new order added
		self.path = append(self.path, prov) // Add us as dep
		self.guesses[prov] = err            // Add the result as guess
		return
	}

	for _, p := range self.path[pathLength:] { // Clear new path bits
		delete(self.State.resolutions, p)
		delete(self.guesses, p)
	}
	self.path = self.path[:pathLength]

	self.guesses[prov] = nil // Make successful guess

	common.Logf("Guess made")
	common.Logf("Adj(%v)", prov)
	common.Indent("  ")
	secondError := self.State.orders[prov].Adjudicate(self)
	common.DeIndent()
	if err == nil {
		common.Logf("%v: Success", prov)
	} else {
		common.Logf("%v: Failure: %v", prov, err)
	}

	if (err == nil && secondError == nil) || (err != nil && secondError != nil) { // Results are the same, and thus only one is consistent
		for _, p := range self.path[pathLength:] { // Clear new path bits
			delete(self.State.resolutions, p)
			delete(self.guesses, p)
		}
		self.path = self.path[:pathLength]
		delete(self.guesses, prov)
		common.Logf("Resolving, exactly one guess consistent")
		self.State.resolutions[prov] = err // Resolve
		return
	}

	common.Logf("BackupRule(%v)", self.path[pathLength:])
	self.State.backupRule(self, self.path[pathLength:]) // Backup rule

	for _, p := range self.path[pathLength:] {
		delete(self.guesses, p)
	}
	self.path = self.path[pathLength:]

	return self.Resolve(prov)
}
