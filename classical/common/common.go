package common

const (
  Sea = 1 << iota
  Land
)

const (
  Army = iota
  Fleet
)

const (
  Coast = Sea | Land

  SC = "SC"

  England = "E"
  France  = "F"
  Germany = "G"
  Russia  = "R"
  Austria = "A"
  Italy   = "I"
  Turkey  = "T"

  Neutral = "N"
)
