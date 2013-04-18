package common

const (
  Sea = 1 << iota
  Land
)

const (
  Coast = Sea | Land

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
)
