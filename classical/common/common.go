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
var ErrIllegalSupport = fmt.Errorf("ErrIllegalSupport")
var ErrMissingSupportee = fmt.Errorf("ErrMissingSupportee")
var ErrInvalidSupportedMove = fmt.Errorf("ErrInvalidSupportedMove")

type ErrBounce struct {
  Province Province
}

func (self ErrBounce) Error() string {
  return fmt.Sprintf("ErrBounce:%v", self.Province)
}
