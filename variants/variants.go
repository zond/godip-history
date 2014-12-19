package variants

import (
	"github.com/zond/godip/classical"
	cla "github.com/zond/godip/classical/common"
	"github.com/zond/godip/classical/orders"
	dip "github.com/zond/godip/common"
	"github.com/zond/godip/state"
)

const (
	Classical = "classical"
	FleetRome = "fleetrome"
)

type Variant struct {
	Start       func() (*state.State, error)
	BlankStart  func() (*state.State, error)
	Blank       func(dip.Phase) *state.State
	Phase       func(int, dip.Season, dip.PhaseType) dip.Phase
	Nations     func() []dip.Nation
	PhaseTypes  func() []dip.PhaseType
	Seasons     func() []dip.Season
	UnitTypes   func() []dip.UnitType
	Orders      func() []dip.Order
	ParseOrders func(map[dip.Nation]map[dip.Province][]string) (map[dip.Province]dip.Adjudicator, error)
	ParseOrder  func([]string) (dip.Adjudicator, error)
	Name        string
}

var Variants = map[string]Variant{
	Classical: Variant{
		Name:  Classical,
		Start: classical.Start,
		Blank: classical.Blank,
		BlankStart: func() (result *state.State, err error) {
			result = classical.Blank(classical.Phase(1900, cla.Fall, cla.Adjustment))
			return
		},
		Phase:       classical.Phase,
		Orders:      orders.Orders,
		ParseOrders: orders.ParseAll,
		ParseOrder:  orders.Parse,
		Nations:     func() []dip.Nation { return cla.Nations },
		PhaseTypes:  func() []dip.PhaseType { return cla.PhaseTypes },
		Seasons:     func() []dip.Season { return cla.Seasons },
		UnitTypes:   func() []dip.UnitType { return cla.UnitTypes },
	},
	FleetRome: Variant{
		Name: FleetRome,
		Start: func() (result *state.State, err error) {
			if result, err = classical.Start(); err != nil {
				return
			}
			result.RemoveUnit(dip.Province("rom"))
			if err = result.SetUnit(dip.Province("rom"), dip.Unit{
				Type:   cla.Fleet,
				Nation: cla.Italy,
			}); err != nil {
				return
			}
			return
		},
		Blank:       classical.Blank,
		Phase:       classical.Phase,
		Orders:      orders.Orders,
		ParseOrders: orders.ParseAll,
		ParseOrder:  orders.Parse,
		Nations:     func() []dip.Nation { return cla.Nations },
		PhaseTypes:  func() []dip.PhaseType { return cla.PhaseTypes },
		Seasons:     func() []dip.Season { return cla.Seasons },
		UnitTypes:   func() []dip.UnitType { return cla.UnitTypes },
	},
}
