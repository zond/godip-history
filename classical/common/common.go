package common

import (
  graph "github.com/zond/godip/graph"
)

var Coast = []graph.Flag{Sea, Land}

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
)
