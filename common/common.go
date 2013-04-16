package common

type PhaseType int

const (
  MovementPhase = iota
  RetreatPhase
  BuildPhase
)

type Phase interface {
  Year() int
  Type() PhaseType
}

type Order interface {
}

type Graph interface {
}

type State interface {
  Resolve() (result State, err error)
  Phase() (result Phase)
}
