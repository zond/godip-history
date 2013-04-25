package classical

import (
  "fmt"
  cla "github.com/zond/godip/classical/common"
  "github.com/zond/godip/classical/orders"
  "github.com/zond/godip/classical/start"
  dip "github.com/zond/godip/common"
  "github.com/zond/godip/state"
)

func Blank(phase dip.Phase) *state.State {
  return state.New(start.Graph(), phase, BackupRule, DefaultOrderGenerator)
}

func Start() *state.State {
  return state.New(start.Graph(), &phase{1901, cla.Spring, cla.Movement}, BackupRule, DefaultOrderGenerator).
    SetUnits(start.Units()).
    SetSupplyCenters(start.SupplyCenters())
}

func DefaultOrderGenerator(prov dip.Province) dip.Order {
  return orders.Hold(prov)
}

/*
BackupRule will make sets of only Move orders succeed, while orders with at least one Convoy all fail.
Any other alternative will cause a panic.
*/
func BackupRule(resolver dip.Resolver, prov dip.Province, deps map[dip.Province]bool) error {
  only_moves := true
  convoys := false
  for prov, _ := range deps {
    if order, _, ok := resolver.Order(prov); ok {
      if order.Type() != cla.Move {
        only_moves = false
      }
      if order.Type() == cla.Convoy {
        convoys = true
      }
    }
  }

  if only_moves {
    return nil
  }
  if convoys {
    return cla.ErrConvoyParadox
  }
  panic(fmt.Errorf("Unknown circular dependency between %v", deps))
}
