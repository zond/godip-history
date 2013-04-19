package classical

import (
  cla "github.com/zond/godip/classical/common"
  "github.com/zond/godip/classical/orders"
  dip "github.com/zond/godip/common"
  "testing"
)

func assertOrderValidity(t *testing.T, validator dip.Validator, order dip.Order, err error) {
  if e := order.Validate(validator); e != err {
    t.Errorf("%v should validate to %v, but got %v", order, err, e)
  }
}

func TestMoveOrderValidation(t *testing.T) {
  judge := Start()
  // Happy path fleet
  assertOrderValidity(t, judge, orders.Move("bre", "mid"), nil)
  // Happy path army
  assertOrderValidity(t, judge, orders.Move("mun", "ruh"), nil)
  // Too far
  assertOrderValidity(t, judge, orders.Move("bre", "wes"), cla.ErrIllegalDistance)
  // Fleet on land
  assertOrderValidity(t, judge, orders.Move("bre", "par"), cla.ErrIllegalDestination)
  // Army at sea
  assertOrderValidity(t, judge, orders.Move("smy", "eas"), cla.ErrIllegalDestination)
  // Unknown source
  assertOrderValidity(t, judge, orders.Move("a", "mid"), cla.ErrInvalidSource)
  // Unknown destination
  assertOrderValidity(t, judge, orders.Move("bre", "a"), cla.ErrInvalidDestination)
  // Missing sea path
  assertOrderValidity(t, judge, orders.Move("par", "mos"), cla.ErrMissingSeaPath)
  // No unit
  assertOrderValidity(t, judge, orders.Move("spa", "por"), cla.ErrMissingUnit)
  // Working convoy
  judge.SetUnit("eng", dip.Unit{cla.Fleet, cla.England})
  judge.SetUnit("wal", dip.Unit{cla.Army, cla.England})
  assertOrderValidity(t, judge, orders.Move("wal", "bre"), nil)
  // Missing convoy
  assertOrderValidity(t, judge, orders.Move("wal", "gas"), cla.ErrMissingConvoyPath)
  // Bad phase
  judge.Next()
  assertOrderValidity(t, judge, orders.Move("bre", "mid"), cla.ErrInvalidPhase)
}
