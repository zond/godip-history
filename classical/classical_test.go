package classical

import (
  cla "github.com/zond/godip/classical/common"
  "github.com/zond/godip/classical/orders"
  dip "github.com/zond/godip/common"
  "github.com/zond/godip/judge"
  "testing"
)

func assertOrderValidity(t *testing.T, state *judge.State, order judge.Order, err error) {
  if e := order.Validate(state); e != err {
    t.Errorf("%v should validate to %v, but got %v", order, err, e)
  }
}

func TestMoveOrderValidation(t *testing.T) {
  state := Start()
  // Happy path fleet
  assertOrderValidity(t, state, orders.Move("bre", "mid"), nil)
  // Happy path army
  assertOrderValidity(t, state, orders.Move("mun", "ruh"), nil)
  // Too far
  assertOrderValidity(t, state, orders.Move("bre", "wes"), cla.ErrIllegalDistance)
  // Fleet on land
  assertOrderValidity(t, state, orders.Move("bre", "par"), cla.ErrIllegalDestination)
  // Army at sea
  assertOrderValidity(t, state, orders.Move("smy", "eas"), cla.ErrIllegalDestination)
  // Unknown source
  assertOrderValidity(t, state, orders.Move("a", "mid"), cla.ErrInvalidSource)
  // Unknown destination
  assertOrderValidity(t, state, orders.Move("bre", "a"), cla.ErrInvalidDestination)
  // Missing sea path
  assertOrderValidity(t, state, orders.Move("par", "mos"), cla.ErrMissingSeaPath)
  // No unit
  assertOrderValidity(t, state, orders.Move("spa", "por"), cla.ErrMissingUnit)
  // Working convoy
  state.Units["eng"] = dip.Unit{cla.Fleet, cla.England}
  state.Units["wal"] = dip.Unit{cla.Army, cla.England}
  assertOrderValidity(t, state, orders.Move("wal", "bre"), nil)
  // Missing convoy
  assertOrderValidity(t, state, orders.Move("wal", "gas"), cla.ErrMissingConvoyPath)
  // Bad phase
  state.Phase, _ = state.Phase.Next()
  assertOrderValidity(t, state, orders.Move("bre", "mid"), cla.ErrInvalidPhase)
}
