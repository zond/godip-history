package orders

import (
	"fmt"
	cla "github.com/zond/godip/classical/common"
	dip "github.com/zond/godip/common"
	"time"
)

func Move(source, dest dip.Province) *move {
	return &move{
		targets: []dip.Province{source, dest},
		flags:   make(map[dip.Flag]bool),
	}
}

type move struct {
	targets []dip.Province
	flags   map[dip.Flag]bool
}

func (self *move) String() string {
	return fmt.Sprintf("%v %v %v", self.targets[0], cla.Move, self.targets[1])
}

func (self *move) ViaConvoy() *move {
	self.flags[cla.ViaConvoy] = true
	return self
}

func (self *move) Type() dip.OrderType {
	return cla.Move
}

func (self *move) Targets() []dip.Province {
	return self.targets
}

func (self *move) At() time.Time {
	return time.Now()
}

func (self *move) calcAttackSupport(r dip.Resolver, src, dst dip.Province, forbidden []dip.Nation) int {
	_, supports, _ := r.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
		if o != nil && u != nil {
			for _, ban := range forbidden {
				if ban == u.Nation {
					return false
				}
			}
			if o.Type() == cla.Support && len(o.Targets()) == 3 && o.Targets()[1].Contains(src) && o.Targets()[2].Contains(dst) {
				if err := r.Resolve(p); err == nil {
					return true
				}
			}
		}
		return false
	})
	return len(supports)
}

func (self *move) calcHoldSupport(r dip.Resolver, prov dip.Province) int {
	_, supports, _ := r.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
		if o != nil && u != nil && o.Type() == cla.Support && p.Super() != prov.Super() && len(o.Targets()) == 2 && o.Targets()[1].Super() == prov.Super() {
			if err := r.Resolve(p); err == nil {
				return true
			}
		}
		return false
	})
	return len(supports)
}

func (self *move) Adjudicate(r dip.Resolver) error {
	if r.Phase().Type() == cla.Movement {
		return self.adjudicateMovementPhase(r)
	}
	return self.adjudicateRetreatPhase(r)
}

func (self *move) adjudicateRetreatPhase(r dip.Resolver) error {
	_, competingOrders, _ := r.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
		return p != self.targets[0] && o.Type() == cla.Move && o.Targets()[1] == self.targets[1]
	})
	if len(competingOrders) > 0 {
		return cla.ErrBounce{competingOrders[0].Targets()[0]}
	}
	return nil
}

func (self *move) Flags() map[dip.Flag]bool {
	return self.flags
}

func (self *move) adjudicateMovementPhase(r dip.Resolver) error {
	unit, _, _ := r.Unit(self.targets[0])

	convoyed, err := cla.IsConvoyed(r, self)
	if err != nil {
		return err
	}

	var myForbiddenSupporters []dip.Nation
	opposingForbiddenSupporters := []dip.Nation{
		unit.Nation,
	}

	// competing moves for the same destination
	_, competingOrders, competingUnits := r.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
		return o != nil && u != nil && o.Type() == cla.Move && o.Targets()[0] != self.targets[0] && self.targets[1].Super() == o.Targets()[1].Super()
	})
	for index, competingOrder := range competingOrders {
		myForbiddenSupporters = append(myForbiddenSupporters, competingUnits[index].Nation)
		attackStrength := self.calcAttackSupport(r, self.targets[0], self.targets[1], myForbiddenSupporters) + 1
		myForbiddenSupporters = myForbiddenSupporters[:len(myForbiddenSupporters)-1]
		dip.Logf("%v:vs %v: %v", self, competingOrder, attackStrength)
		if as := self.calcAttackSupport(r, competingOrder.Targets()[0], competingOrder.Targets()[1], opposingForbiddenSupporters) + 1; as >= attackStrength {
			conv, _ := cla.IsConvoyed(r, competingOrder)
			if conv {
				dip.Logf("%v:vs %v: %v", competingOrder, self, as)
				return cla.ErrBounce{competingOrder.Targets()[0]}
			} else {
				dip.Logf("H2HDisl(%v)", self.targets[1])
				dip.Indent("  ")
				if dislodgers, _, _ := r.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
					res := o != nil && // is an order
						u != nil && // is a unit
						o.Type() == cla.Move && // move
						o.Targets()[1].Super() == competingOrder.Targets()[0].Super() && // against the competition
						o.Targets()[0].Super() == competingOrder.Targets()[1].Super() && // from their destination
						u.Nation != competingUnits[index].Nation // not from themselves
					if res {
						if conv, _ := cla.IsConvoyed(r, o); !conv && r.Resolve(p) == nil {
							return true
						}
					}
					return false
				}); len(dislodgers) == 0 {
					dip.DeIndent()
					dip.Logf("F")
					dip.Logf("%v:vs %v: %v", competingOrder, self, as)
					return cla.ErrBounce{competingOrder.Targets()[0]}
				} else {
					dip.DeIndent()
					dip.Logf("%v", dislodgers)
				}
			}
		} else {
			dip.Logf("%v:vs %v: %v", competingOrder, self, as)
		}
	}

	// at destination
	victim, _, hasVictim := r.Unit(self.targets[1])
	if hasVictim {
		myForbiddenSupporters = append(myForbiddenSupporters, victim.Nation)
		attackStrength := self.calcAttackSupport(r, self.targets[0], self.targets[1], myForbiddenSupporters) + 1
		order, prov, _ := r.Order(self.targets[1])
		dip.Logf("%v:vs %v: %v", self, order, attackStrength)
		if order.Type() == cla.Move {
			victimConvoyed, _ := cla.IsConvoyed(r, order)
			if !convoyed && !victimConvoyed && order.Targets()[1] == self.targets[0] {
				as := self.calcAttackSupport(r, order.Targets()[0], order.Targets()[1], opposingForbiddenSupporters) + 1
				dip.Logf("%v:vs %v: %v", order, self, as)
				if victim.Nation == unit.Nation || as >= attackStrength {
					return cla.ErrBounce{self.targets[1]}
				}
			} else {
				dip.Logf("Esc(%v)", order.Targets()[0])
				dip.Indent("  ")
				if err := r.Resolve(prov); err == nil {
					dip.DeIndent()
					dip.Logf("T")
					myForbiddenSupporters = myForbiddenSupporters[:len(myForbiddenSupporters)-1]
				} else {
					dip.DeIndent()
					dip.Logf("%v", err)
					if victim.Nation == unit.Nation || 1 >= attackStrength {
						return cla.ErrBounce{self.targets[1]}
					}
				}
			}
		} else {
			hs := self.calcHoldSupport(r, self.targets[1]) + 1
			dip.Logf("%v: %v", order, hs)
			if victim.Nation == unit.Nation || hs >= attackStrength {
				return cla.ErrBounce{self.targets[1]}
			}
		}
	}

	for index, competingOrder := range competingOrders {
		myForbiddenSupporters = append(myForbiddenSupporters, competingUnits[index].Nation)
		attackStrength := self.calcAttackSupport(r, self.targets[0], self.targets[1], myForbiddenSupporters) + 1
		myForbiddenSupporters = myForbiddenSupporters[:len(myForbiddenSupporters)-1]
		dip.Logf("%v:vs %v: %v", self, competingOrder, attackStrength)
		if as := self.calcAttackSupport(r, competingOrder.Targets()[0], competingOrder.Targets()[1], opposingForbiddenSupporters) + 1; as >= attackStrength {
			conv, _ := cla.IsConvoyed(r, competingOrder)
			if conv {
				dip.Logf("%v:vs %v: %v", competingOrder, self, as)
				return cla.ErrBounce{competingOrder.Targets()[0]}
			} else {
				dip.Logf("H2HDisl(%v)", self.targets[1])
				dip.Indent("  ")
				if dislodgers, _, _ := r.Find(func(p dip.Province, o dip.Order, u *dip.Unit) bool {
					res := o != nil && // is an order
						u != nil && // is a unit
						o.Type() == cla.Move && // move
						o.Targets()[1].Super() == competingOrder.Targets()[0].Super() && // against the competition
						o.Targets()[0].Super() == competingOrder.Targets()[1].Super() && // from their destination
						u.Nation != competingUnits[index].Nation // not from themselves
					if res {
						if conv, _ := cla.IsConvoyed(r, o); !conv && r.Resolve(p) == nil {
							return true
						}
					}
					return false
				}); len(dislodgers) == 0 {
					dip.DeIndent()
					dip.Logf("F")
					dip.Logf("%v:vs %v: %v", competingOrder, self, as)
					return cla.ErrBounce{competingOrder.Targets()[0]}
				} else {
					dip.DeIndent()
					dip.Logf("%v", dislodgers)
				}
			}
		} else {
			dip.Logf("%v:vs %v: %v", competingOrder, self, as)
		}
	}

	return nil
}

func (self *move) Validate(v dip.Validator) error {
	if v.Phase().Type() == cla.Movement {
		return self.validateMovementPhase(v)
	} else if v.Phase().Type() == cla.Retreat {
		return self.validateRetreatPhase(v)
	}
	return cla.ErrInvalidPhase
}

func (self *move) validateRetreatPhase(v dip.Validator) error {
	if !v.Graph().Has(self.targets[0]) {
		return cla.ErrInvalidSource
	}
	if !v.Graph().Has(self.targets[1]) {
		return cla.ErrInvalidDestination
	}
	var unit dip.Unit
	var ok bool
	if unit, self.targets[0], ok = v.Dislodged(self.targets[0]); !ok {
		return cla.ErrMissingUnit
	}
	var err error
	if self.targets[1], err = cla.AnyMovePossible(v, self.targets[0], self.targets[1], unit.Type == cla.Army, false, false); err != nil {
		return err
	}
	if v.IsDislodger(self.targets[1], self.targets[0]) {
		return cla.ErrIllegalRetreat
	}
	return nil
}

func (self *move) validateMovementPhase(v dip.Validator) error {
	if !v.Graph().Has(self.targets[0]) {
		return cla.ErrInvalidSource
	}
	if !v.Graph().Has(self.targets[1]) {
		return cla.ErrInvalidDestination
	}
	var unit dip.Unit
	var ok bool
	if unit, self.targets[0], ok = v.Unit(self.targets[0]); !ok {
		return cla.ErrMissingUnit
	}
	var err error
	if self.targets[1], err = cla.AnyMovePossible(v, self.targets[0], self.targets[1], unit.Type == cla.Army, true, false); err != nil {
		return err
	}
	return nil
}

func (self *move) Execute(state dip.State) {
	if state.Phase().Type() == cla.Retreat {
		state.Retreat(self.targets[0], self.targets[1])
	} else {
		state.Move(self.targets[0], self.targets[1])
	}
}
