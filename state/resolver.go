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
			common.Logf("Success")
		} else {
			common.Logf("Failure: %v", err)
		}
	}()

	var ok bool
	if err, ok = self.State.resolutions[prov]; ok { // Already resolved
		return
	}

	if err, ok = self.guesses[prov]; ok { // Already guessed
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

	err = self.State.orders[prov].Adjudicate(self)

	if len(self.path) == pathLength { // No new guesses made
		self.State.resolutions[prov] = err // Resolve
		return
	}

	if self.path[pathLength] != prov { // We are not the first new order added
		self.path = append(self.path, prov) // Add us as dep
		self.guesses[prov] = err            // Add the result as guess
		return
	}

	for _, p := range self.path[pathLength:] { // Clear new path bits
		delete(self.guesses, p)
	}
	self.path = self.path[:pathLength]

	self.guesses[prov] = nil // Make successful guess

	secondError := self.State.orders[prov].Adjudicate(self)

	if (err == nil && secondError == nil) || (err != nil && secondError != nil) { // Results are the same, and thus only one is consistent
		for _, p := range self.path[pathLength:] { // Clear new path bits
			delete(self.guesses, p)
		}
		self.path = self.path[:pathLength]
		self.State.resolutions[prov] = err // Resolve
		return
	}

	self.State.backupRule(self, self.path[pathLength:]) // Backup rule

	for _, p := range self.path[pathLength:] {
		delete(self.guesses, p)
	}
	self.path = self.path[pathLength:]

	return self.Resolve(prov)
}
