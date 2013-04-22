package common

import (
  "fmt"
  . "github.com/zond/godip/common"
)

const (
  Sea  = "S"
  Land = "L"

  Army  = "A"
  Fleet = "F"

  England = "E"
  France  = "F"
  Germany = "G"
  Russia  = "R"
  Austria = "A"
  Italy   = "I"
  Turkey  = "T"

  Neutral = "N"

  Spring = "S"
  Winter = "W"
  Fall   = "F"

  Movement = "M"
  Build    = "B"
  Retreat  = "R"

  Move    = "M"
  Hold    = "H"
  Convoy  = "C"
  Support = "S"
)

var Coast = []Flag{Sea, Land}

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
var ErrIllegalSupportDestinationNationality = fmt.Errorf("ErrIllegalSupportDestinationNationality")
var ErrMissingSupportUnit = fmt.Errorf("ErrMissingSupportUnit")
var ErrInvalidSupportMove = fmt.Errorf("ErrInvalidSupportMove")
var ErrIllegalConvoy = fmt.Errorf("ErrIllegalConvoy")
var ErrMissingConvoyee = fmt.Errorf("ErrMissingConvoyee")
var ErrIllegalBuild = fmt.Errorf("ErrIllegalBuild")
var ErrOccupiedSupplyCenter = fmt.Errorf("ErrOccupiedSupplyCenter")

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

func ConvoyPossible(v Validator, src, dst Province, checkOrders bool) error {
  unit := v.Unit(src)
  if unit == nil {
    return ErrMissingUnit
  }
  if unit.Type != Army {
    return ErrIllegalConvoy
  }
  if path := v.Graph().Path(src, dst, func(name Province, flags map[Flag]bool, sc *Nationality) bool {
    if u := v.Unit(name); u != nil && u.Type == Fleet {
      if !checkOrders {
        return true
      }
      if order := v.Order(name); order != nil && order.Type() == Convoy && order.Targets()[1] == src && order.Targets()[2] == dst {
        if r, ok := v.(Resolver); ok {
          if err := r.Resolve(name); err == nil {
            return true
          }
        } else {
          return true
        }
      }
    }
    return false
  }); path == nil {
    return ErrMissingConvoyPath
  }
  return nil
}

/*
AnyMovePossible returns true if MovePossible would return true for any movement between src and any coast of dst.
*/
func AnyMovePossible(v Validator, src, dst Province) error {
  var err error
  for coast, _ := range v.Graph().Coasts(dst) {
    if err = MovePossible(v, src, coast, false, false); err == nil {
      return nil
    }
  }
  return err
}

/*
PossibleMove returns true if a move from src to dst is possible in v.

It will validate that the move is theoretically possible without privileged information.

It will (if allowConvoy and the need for convoying) validate the presence of fleets along the path.

It will (if allowConvoy, the need for convoying and resolveConvoy) validate presence of successful and relevant convoy orders along the path.
*/
func MovePossible(v Validator, src, dst Province, allowConvoy, checkConvoyOrders bool) error {
  if !v.Graph().Has(src) {
    return ErrInvalidSource
  }
  if !v.Graph().Has(dst) {
    return ErrInvalidDestination
  }
  unit := v.Unit(src)
  if unit == nil {
    return ErrMissingUnit
  }
  if unit.Type == Army {
    if !v.Graph().Flags(dst)[Land] {
      return ErrIllegalDestination
    }
  } else if unit.Type == Fleet {
    if !v.Graph().Flags(dst)[Sea] {
      return ErrIllegalDestination
    }
  } else {
    panic(fmt.Errorf("Unknown unit type %v", unit.Type))
  }
  path := v.Graph().Path(src, dst, nil)
  if path == nil {
    return ErrMissingPath
  }
  if len(path) > 1 {
    if allowConvoy {
      return ConvoyPossible(v, src, dst, checkConvoyOrders)
    }
    return ErrIllegalDistance
  }
  return nil
}
