package common

import (
	"fmt"
	. "github.com/zond/godip/common"
	"sort"
)

const (
	Sea  Flag = "Sea"
	Land Flag = "Land"

	Army  UnitType = "Army"
	Fleet UnitType = "Fleet"

	England Nation = "England"
	France  Nation = "France"
	Germany Nation = "Germany"
	Russia  Nation = "Russia"
	Austria Nation = "Austria"
	Italy   Nation = "Italy"
	Turkey  Nation = "Turkey"
	Neutral Nation = "Neutral"

	Spring Season = "Spring"
	Fall   Season = "Fall"

	Movement   PhaseType = "Movement"
	Retreat    PhaseType = "Retreat"
	Adjustment PhaseType = "Adjustment"

	Build   OrderType = "Build"
	Move    OrderType = "Move"
	Hold    OrderType = "Hold"
	Convoy  OrderType = "Convoy"
	Support OrderType = "Support"
	Disband OrderType = "Disband"

	ViaConvoy Flag = "C"
)

var Coast = []Flag{Sea, Land}
var Nations = []Nation{Austria, England, France, Germany, Italy, Turkey, Russia}

var ErrInvalidSource = fmt.Errorf("ErrInvalidSource")
var ErrInvalidDestination = fmt.Errorf("ErrInvalidDestination")
var ErrInvalidTarget = fmt.Errorf("ErrInvalidTarget")
var ErrInvalidPhase = fmt.Errorf("ErrInvalidPhase")
var ErrMissingUnit = fmt.Errorf("ErrMissingUnit")
var ErrIllegalDestination = fmt.Errorf("ErrIllegalDestination")
var ErrMissingPath = fmt.Errorf("ErrMissingPath")
var ErrMissingSeaPath = fmt.Errorf("ErrMissingSeaPath")
var ErrMissingConvoyPath = fmt.Errorf("ErrMissignConvoyPath")
var ErrIllegalDistance = fmt.Errorf("ErrIllegalDistance")
var ErrConvoyParadox = fmt.Errorf("ErrConvoyParadox")
var ErrMissingConvoy = fmt.Errorf("ErrMissingConvoy")
var ErrIllegalSupportPosition = fmt.Errorf("ErrIllegalSupportPosition")
var ErrIllegalSupportDestination = fmt.Errorf("ErrIllegalSupportDestination")
var ErrIllegalSupportDestinationNation = fmt.Errorf("ErrIllegalSupportDestinationNation")
var ErrMissingSupportUnit = fmt.Errorf("ErrMissingSupportUnit")
var ErrInvalidSupportMove = fmt.Errorf("ErrInvalidSupportMove")
var ErrIllegalConvoyUnit = fmt.Errorf("ErrIllegalConvoyUnit")
var ErrIllegalConvoyMove = fmt.Errorf("ErrIllegalConvoyMove")
var ErrMissingConvoyee = fmt.Errorf("ErrMissingConvoyee")
var ErrIllegalBuild = fmt.Errorf("ErrIllegalBuild")
var ErrIllegalDisband = fmt.Errorf("ErrIllegalDisband")
var ErrOccupiedSupplyCenter = fmt.Errorf("ErrOccupiedSupplyCenter")
var ErrMissingSupplyCenter = fmt.Errorf("ErrMissingSupplyCenter")
var ErrMissingSurplus = fmt.Errorf("ErrMissingSurplus")
var ErrIllegalUnitType = fmt.Errorf("ErrIllegalUnitType")
var ErrMissingDeficit = fmt.Errorf("ErrMissingDeficit")
var ErrOccupiedDestination = fmt.Errorf("ErrOccupiedDestination")
var ErrIllegalRetreat = fmt.Errorf("ErrIllegalRetreat")
var ErrForcedDisband = fmt.Errorf("ErrForcedDisband")
var ErrHostileSupplyCenter = fmt.Errorf("ErrHostileSupplyCenter")

type ErrConvoyDislodged struct {
	Province Province
}

func (self ErrConvoyDislodged) Error() string {
	return fmt.Sprintf("ErrConvoyDislodged:%v", self.Province)
}

type ErrSupportBroken struct {
	Province Province
}

func (self ErrSupportBroken) Error() string {
	return fmt.Sprintf("ErrSupportBroken:%v", self.Province)
}

type ErrBounce struct {
	Province Province
}

func (self ErrBounce) Error() string {
	return fmt.Sprintf("ErrBounce:%v", self.Province)
}

func convoyPathExists(v Validator, src, dst, realSrc, realDst Province, checkOrders bool) bool {
	return v.Graph().Path(src, dst, func(name Province, edgeFlags, nodeFlags map[Flag]bool, sc *Nation) bool {
		if name.Contains(dst) {
			return true
		}
		if nodeFlags[Land] {
			return false
		}
		if u, _, ok := v.Unit(name); ok && u.Type == Fleet {
			if !checkOrders {
				return true
			}
			if order, prov, ok := v.Order(name); ok && order.Type() == Convoy && order.Targets()[1].Contains(realSrc) && order.Targets()[2].Contains(realDst) {
				if r, ok := v.(Resolver); ok {
					if err := r.Resolve(prov); err != nil {
						return false
					}
				}
				return true
			}
		}
		return false
	}) != nil
}

/*
convoyPossible will return whether it is possible to convoy from src to dst in v.

It will validate the presence of successful and relevant convoy orders if checkOrders.
*/
func convoyPossible(v Validator, src, dst Province, checkOrders bool, mustNation *Nation) error {
	unit, _, ok := v.Unit(src)
	if !ok {
		return ErrMissingUnit
	}
	if unit.Type != Army {
		return ErrIllegalConvoyUnit
	}
	if mustNation == nil {
		if convoyPathExists(v, src, dst, src, dst, checkOrders) {
			return nil
		}
	} else {
		waypoints, _, _ := v.Find(func(p Province, o Order, u *Unit) bool {
			if u.Nation == *mustNation && u.Type == Fleet && o.Type() == Convoy {
				if !checkOrders {
					return true
				}
				if o.Type() == Convoy && o.Targets()[1].Contains(src) && o.Targets()[2].Contains(dst) {
					if r, ok := v.(Resolver); ok {
						if err := r.Resolve(p); err != nil {
							return false
						}
					}
					return true
				}
			}
			return false
		})
		for _, waypoint := range waypoints {
			if convoyPathExists(v, src, waypoint, src, dst, false) && convoyPathExists(v, waypoint, dst, src, dst, false) {
				return nil
			}
		}
	}
	return ErrMissingConvoyPath
}

/*
AnyConvoyPossible will return whether it is possible to convoy from any coast in src to any coast in dst.

It will validate the presence of successful and relevant convoy orders if checkConvoyOrders.
*/
func AnyConvoyPossible(v Validator, src, dst Province, checkConvoyOrders bool, mustNation *Nation) (err error) {
	if err = convoyPossible(v, src, dst, checkConvoyOrders, mustNation); err == nil {
		return
	}
	for _, srcCoast := range v.Graph().Coasts(src) {
		for _, dstCoast := range v.Graph().Coasts(dst) {
			if err = convoyPossible(v, srcCoast, dstCoast, checkConvoyOrders, mustNation); err == nil {
				return
			}
		}
	}
	return
}

/*
AnySupportPossible will return whether the src can support the dst province by checking if movement is possible from src to any coast of dst.
*/
func AnySupportPossible(v Validator, src, dst Province) (err error) {
	if err = movePossible(v, src, dst, false, false); err == nil {
		return
	}
	for _, coast := range v.Graph().Coasts(dst) {
		if err = movePossible(v, src, coast, false, false); err == nil {
			return
		}
	}
	return
}

/*
AnyMovePossible returns true if movePossible would return true for any movement between src and any coast of dst.

It will check sub provinces (coasts) of dstonly if lax or if dst is the name of the super province.

It will allow convoys if allowConvoy.

It will validate presence of successful and relevant convoy orders along the path if checkConvoyOrders.
*/
func AnyMovePossible(v Validator, src, dst Province, lax, allowConvoy, checkConvoyOrders bool) (dstCoast Province, err error) {
	dstCoast = dst
	if err = movePossible(v, src, dst, allowConvoy, checkConvoyOrders); err == nil {
		return
	}
	if lax || dst.Super() == dst {
		var options []Province
		for _, coast := range v.Graph().Coasts(dst) {
			if err2 := movePossible(v, src, coast, allowConvoy, checkConvoyOrders); err2 == nil {
				options = append(options, coast)
			}
		}
		if len(options) > 0 {
			if lax || len(options) == 1 {
				dstCoast, err = options[0], nil
			}
		}
	}
	return
}

/*
movePossible returns true if a move from src to dst is possible in v.

It will validate that the move is theoretically possible without privileged information.

It will (if allowConvoy and the need for convoying) validate the presence of fleets along the path.

It will (if allowConvoy, the need for convoying and resolveConvoy) validate presence of successful and relevant convoy orders along the path.
*/
func movePossible(v Validator, src, dst Province, allowConvoy, checkConvoyOrders bool) error {
	if !v.Graph().Has(src) {
		return ErrInvalidSource
	}
	if !v.Graph().Has(dst) {
		return ErrInvalidDestination
	}
	unit, _, ok := v.Unit(src)
	if !ok {
		return ErrMissingUnit
	}
	var filter PathFilter
	if unit.Type == Army {
		if !v.Graph().Flags(dst)[Land] {
			return ErrIllegalDestination
		}
		filter = func(p Province, ef, nf map[Flag]bool, sc *Nation) bool {
			return ef[Land] && nf[Land]
		}
	} else if unit.Type == Fleet {
		if !v.Graph().Flags(dst)[Sea] {
			return ErrIllegalDestination
		}
		filter = func(p Province, ef, nf map[Flag]bool, sc *Nation) bool {
			return ef[Sea] && nf[Sea]
		}
	} else {
		panic(fmt.Errorf("Unknown unit type %v", unit.Type))
	}
	if path := v.Graph().Path(src, dst, filter); path == nil || len(path) > 1 {
		if allowConvoy {
			return AnyConvoyPossible(v, src, dst, checkConvoyOrders, nil)
		}
		if path == nil {
			return ErrMissingPath
		} else {
			return ErrIllegalDistance
		}
	}
	return nil
}

func AdjustmentStatus(v Validator, me Nation) (builds Orders, disbands Orders, balance int) {
	scs := 0
	for prov, nat := range v.SupplyCenters() {
		if nat == me {
			scs += 1
			if order, _, ok := v.Order(prov); ok && order.Type() == Build {
				builds = append(builds, order)
			}
		}
	}

	units := 0
	for prov, unit := range v.Units() {
		if unit.Nation == me {
			units += 1
			if order, _, ok := v.Order(prov); ok && order.Type() == Disband {
				disbands = append(disbands, order)
			}
		}
	}

	sort.Sort(builds)
	sort.Sort(disbands)

	change := scs - units
	if change > 0 {
		disbands = nil
		builds = builds[:Max(0, Min(len(builds)-1, change))]
	} else if change < 0 {
		builds = nil
		disbands = disbands[:Max(0, Min(len(disbands)-1, change))]
	} else {
		builds = nil
		disbands = nil
	}

	return
}

func IsConvoyed(r Resolver, order Order) (result bool, err error) {
	if order.Type() != Move {
		panic(fmt.Errorf("%v is not a Move order", order))
	}

	unit, _, _ := r.Unit(order.Targets()[0])
	// is convoyed?
	if unit.Type == Army {
		steps := r.Graph().Path(order.Targets()[0], order.Targets()[1], nil)
		if order.Flags()[ViaConvoy] || len(steps) > 1 || AnyConvoyPossible(r, order.Targets()[0], order.Targets()[1], true, &unit.Nation) == nil {
			Logf("Conv(%v)", order)
			Indent("  ")
			err = AnyConvoyPossible(r, order.Targets()[0], order.Targets()[1], true, nil)
			if err != nil {
				DeIndent()
				Logf("%v", err)
				if len(steps) == 1 {
					err = nil
				}
			} else {
				DeIndent()
				Logf("T")
				result = true
			}
		}
	}
	return
}

/*
HoldSupport returns successful supports of a hold in prov.
*/
func HoldSupport(r Resolver, prov Province) int {
	_, supports, _ := r.Find(func(p Province, o Order, u *Unit) bool {
		if o != nil && u != nil && o.Type() == Support && p.Super() != prov.Super() && len(o.Targets()) == 2 && o.Targets()[1].Super() == prov.Super() {
			if err := r.Resolve(p); err == nil {
				return true
			}
		}
		return false
	})
	return len(supports)
}

/*
MoveSupport returns the successful supports of movement from src to dst, discounting the nations in forbiddenSupports.
*/
func MoveSupport(r Resolver, src, dst Province, forbiddenSupports []Nation) int {
	_, supports, _ := r.Find(func(p Province, o Order, u *Unit) bool {
		if o != nil && u != nil {
			if o.Type() == Support && len(o.Targets()) == 3 && o.Targets()[1].Contains(src) && o.Targets()[2].Contains(dst) {
				for _, ban := range forbiddenSupports {
					if ban == u.Nation {
						return false
					}
				}
				if err := r.Resolve(p); err == nil {
					return true
				}
			}
		}
		return false
	})
	return len(supports)
}
