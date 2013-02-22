package godip

const (
	Hold OrderType = iota
	Move
	Support
	Convoy
	Build
	Disband
)

const (
	ViaConvoy OrderFlag = 1 << iota
)

type OrderFlag int

type OrderType int

type Order interface {
	Type() OrderType
	GetFlags() OrderFlag
	SetFlags(f OrderFlag)
}
