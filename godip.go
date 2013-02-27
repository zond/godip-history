package godip

// Possible resolutions of an order. 
type Resolution int

const (
  Fails Resolution = iota
  Succeeds
)

type ResolutionReason int

const (
  NoReason ResolutionReason = iota
  NoConvoyReason
)

// The resolution of an order, can be in three states.
type state int

const (
  unresolved state = iota
  guessing
  resolved
)

type UnitType int

const (
  Army UnitType = iota
  Navy
)

func (self UnitType) canConvoy(u UnitType) bool {
  return self == Navy && u == Army
}

type Unit struct {
  Nationality string
  Type        UnitType
}

type SlotFlags int

const (
  CoastSlot SlotFlags = 1 << iota
)

type SlotType int

const (
  LandSlot = iota
  SeaSlot
)

type slot struct {
  slotName      string
  provinceName  string
  connections   map[UnitType]map[string]*slot
  typ           SlotType
  flags         SlotFlags
  unit          *Unit
  dislodgedUnit *Unit
}

func (self slot) canConvoy(u UnitType) bool {
  return self.typ == SeaSlot && self.flags&CoastSlot == 0 && u == Army
}

type World struct {
  provinces map[string]map[string]*slot
  slots     map[string]*slot
}

func (self *World) getUnitAt(pos string) (result *Unit) {
  if slot := self.slots[pos]; slot != nil {
    return slot.Unit
  }
  return
}

func (self *World) hasMovePath(typ UnitType, from, to string) (result bool) {
  if slot := self.slots[from]; slot != nil {
    if conns := slot.Connections[typ]; conn != nil {
      result = conns[to] != nil
    }
  }
  return
}
