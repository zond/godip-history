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
  Autumn = "A"

  Movement = "M"
  Build    = "B"
  Retreat  = "R"

  Move   = "M"
  Convoy = "C"
)

var Coast = []Flag{Sea, Land}

var ErrTargetLength = fmt.Errorf("ErrTargetLength")
var ErrInvalidSource = fmt.Errorf("ErrInvalidSource")
var ErrInvalidTarget = fmt.Errorf("ErrInvalidTarget")
var ErrInvalidDestination = fmt.Errorf("ErrInvalidDestination")
var ErrInvalidPhase = fmt.Errorf("ErrInvalidPhase")
var ErrMissingUnit = fmt.Errorf("ErrMissingUnit")
var ErrIllegalDestination = fmt.Errorf("ErrIllegalDestination")
var ErrIllegalDistance = fmt.Errorf("ErrIllegalDistance")
