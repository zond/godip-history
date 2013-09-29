package orders

import (
	"fmt"
	cla "github.com/zond/godip/classical/common"
	dip "github.com/zond/godip/common"
	"time"
)

func Hold(source dip.Province) *hold {
	return &hold{
		targets: []dip.Province{source},
	}
}

type hold struct {
	targets []dip.Province
}

func (self *hold) String() string {
	return fmt.Sprintf("%v %v", self.targets[0], cla.Hold)
}

func (self *hold) Flags() map[dip.Flag]bool {
	return nil
}

func (self *hold) Type() dip.OrderType {
	return cla.Hold
}

func (self *hold) Targets() []dip.Province {
	return self.targets
}

func (self *hold) At() time.Time {
	return time.Now()
}

func (self *hold) Adjudicate(r dip.Resolver) error {
	return nil
}

func (self *hold) Options(v dip.Validator, src dip.Province) (nation *dip.Nation, result *dip.Option, found bool) {
	if v.Phase().Type() == cla.Movement {
		if v.Graph().Has(src) {
			var unit dip.Unit
			var ok bool
			if unit, src, ok = v.Unit(src); ok {
				found = true
				nation = &unit.Nation
				result = &dip.Option{
					Value: src,
				}
			}
		}
	}
	return
}

func (self *hold) Validate(v dip.Validator) error {
	if v.Phase().Type() != cla.Movement {
		return cla.ErrInvalidPhase
	}
	if !v.Graph().Has(self.targets[0]) {
		return cla.ErrInvalidTarget
	}
	var ok bool
	if _, self.targets[0], ok = v.Unit(self.targets[0]); !ok {
		return cla.ErrMissingUnit
	}
	return nil
}

func (self *hold) Execute(state dip.State) {
}
