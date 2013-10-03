package orders

import (
	"fmt"
	cla "github.com/zond/godip/classical/common"
	dip "github.com/zond/godip/common"
	"time"
)

func Convoy(source, from, to dip.Province) *convoy {
	return &convoy{
		targets: []dip.Province{source, from, to},
	}
}

type convoy struct {
	targets []dip.Province
}

func (self *convoy) String() string {
	return fmt.Sprintf("%v %v %v", self.targets[0], cla.Convoy, self.targets[1:])
}

func (self *convoy) Flags() map[dip.Flag]bool {
	return nil
}

func (self *convoy) At() time.Time {
	return time.Now()
}

func (self *convoy) Type() dip.OrderType {
	return cla.Convoy
}

func (self *convoy) DisplayType() dip.OrderType {
	return cla.Convoy
}

func (self *convoy) Targets() []dip.Province {
	return self.targets
}

func (self *convoy) Adjudicate(r dip.Resolver) error {
	unit, _, _ := r.Unit(self.targets[0])
	if breaks, _, _ := r.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
		return (o.Type() == cla.Move && // move
			o.Targets()[1] == self.targets[0] && // against us
			u.Nation != unit.Nation && // not friendly
			r.Resolve(p) == nil)
	}); len(breaks) > 0 {
		return cla.ErrConvoyDislodged{breaks[0]}
	}
	return nil
}

func (self *convoy) Options(v dip.Validator, src dip.Province) (nation dip.Nation, actualSrc dip.Province, result dip.Options, found bool) {
	if v.Phase().Type() == cla.Movement {
		if v.Graph().Has(src) {
			var convoyer dip.Unit
			var ok bool
			if convoyer, actualSrc, ok = v.Unit(src); ok && convoyer.Type == cla.Fleet {
				nation = convoyer.Nation
				for mvSrc, unit := range v.Units() {
					if unit.Type == cla.Army {
						for _, mvDst := range v.Graph().Provinces() {
							if mvDst.Sub() == "" {
								if part1 := v.Graph().Path(mvSrc, src, cla.PossibleConvoyPathFilter(v, mvSrc, mvDst, false, false)); part1 != nil {
									if part2 := v.Graph().Path(src, mvDst, cla.PossibleConvoyPathFilter(v, mvSrc, mvDst, false, true)); part2 != nil {
										found = true
										if result == nil {
											result = dip.Options{}
										}
										opt, f := result[mvSrc]
										if !f {
											opt = dip.Options{}
											result[mvSrc] = opt
										}
										opt[mvDst] = nil
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

func (self *convoy) Validate(v dip.Validator) error {
	if v.Phase().Type() != cla.Movement {
		return cla.ErrInvalidPhase
	}
	if !v.Graph().Has(self.targets[0]) {
		return cla.ErrInvalidSource
	}
	if !v.Graph().Has(self.targets[1]) {
		return cla.ErrInvalidTarget
	}
	if !v.Graph().Has(self.targets[2]) {
		return cla.ErrInvalidTarget
	}
	var convoyer dip.Unit
	var ok bool
	if convoyer, self.targets[0], ok = v.Unit(self.targets[0]); !ok {
		return cla.ErrMissingUnit
	} else if convoyer.Type != cla.Fleet {
		return cla.ErrIllegalConvoyer
	}
	var convoyee dip.Unit
	if convoyee, self.targets[1], ok = v.Unit(self.targets[1]); !ok {
		return cla.ErrMissingConvoyee
	} else if convoyee.Type != cla.Army {
		return cla.ErrIllegalConvoyee
	}
	if cla.AnyConvoyPath(v, self.targets[1], self.targets[2], false, nil) == nil {
		return cla.ErrIllegalConvoyMove
	}
	return nil
}

func (self *convoy) Execute(state dip.State) {
}
