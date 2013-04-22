package classical

import (
  cla "github.com/zond/godip/classical/common"
  "github.com/zond/godip/classical/orders"
  dip "github.com/zond/godip/common"
  "github.com/zond/godip/judge"
  "reflect"
  "testing"
)

func assertOrderValidity(t *testing.T, validator dip.Validator, order dip.Order, err error) {
  if e := order.Validate(validator); e != err {
    t.Errorf("%v should validate to %v, but got %v", order, err, e)
  }
}

func assertMove(t *testing.T, j *judge.Judge, src, dst dip.Province, success bool) {
  if success {
    unit := j.Unit(src)
    if unit == nil {
      t.Errorf("Should be a unit at %v", src)
    }
    j.SetOrder(src, orders.Move(src, dst))
    j.Next()
    if err, ok := j.Errors()[src]; ok {
      t.Errorf("Move from %v to %v should have worked, got %v", src, dst, err)
    }
    if now := j.Unit(src); now != nil && reflect.DeepEqual(now, unit) {
      t.Errorf("%v should have moved from %v", now, src)
    }
    if now := j.Unit(dst); now == nil || !reflect.DeepEqual(now, unit) {
      t.Errorf("%v should be at %v now", unit, dst)
    }
  } else {
    unit := j.Unit(src)
    j.SetOrder(src, orders.Move(src, dst))
    j.Next()
    if _, ok := j.Errors()[src]; !ok {
      t.Errorf("Move from %v to %v should not have worked", src, dst)
    }
    if now := j.Unit(src); now == nil && !reflect.DeepEqual(now, unit) {
      t.Errorf("%v should not have moved from %v", now, src)
    }
  }
}

func TestSupportValidation(t *testing.T) {
  judge := Start()
  // Happy paths
  assertOrderValidity(t, judge, orders.Support("bre", "par", "gas"), nil)
  assertOrderValidity(t, judge, orders.Support("par", "bre"), nil)
  assertOrderValidity(t, judge, orders.Support("par", "bre", "gas"), nil)
  judge.SetUnit("spa/sc", dip.Unit{cla.Fleet, cla.France})
  judge.SetUnit("por", dip.Unit{cla.Fleet, cla.France})
  judge.SetUnit("gol", dip.Unit{cla.Fleet, cla.France})
  assertOrderValidity(t, judge, orders.Support("spa/sc", "por", "mid"), nil)
  assertOrderValidity(t, judge, orders.Support("gol", "mar", "spa"), nil)
  // Missing unit
  assertOrderValidity(t, judge, orders.Support("ruh", "kie", "hol"), cla.ErrMissingUnit)
  // Missing supportee
  assertOrderValidity(t, judge, orders.Support("ber", "sil"), cla.ErrMissingSupportUnit)
  // Illegal support
  assertOrderValidity(t, judge, orders.Support("bre", "par"), cla.ErrIllegalSupportPosition)
  assertOrderValidity(t, judge, orders.Support("mar", "spa/nc", "por"), cla.ErrIllegalSupportDestination)
  assertOrderValidity(t, judge, orders.Support("spa/nc", "mar", "gol"), cla.ErrIllegalSupportDestination)
  // Illegal moves
  assertOrderValidity(t, judge, orders.Support("mar", "spa/nc", "bur"), cla.ErrInvalidSupportMove)
}

func TestMoveValidation(t *testing.T) {
  judge := Start()
  // Happy path fleet
  assertOrderValidity(t, judge, orders.Move("bre", "mid"), nil)
  // Happy path army
  assertOrderValidity(t, judge, orders.Move("mun", "ruh"), nil)
  // Too far
  assertOrderValidity(t, judge, orders.Move("bre", "wes"), cla.ErrIllegalConvoy)
  // Fleet on land
  assertOrderValidity(t, judge, orders.Move("bre", "par"), cla.ErrIllegalDestination)
  // Army at sea
  assertOrderValidity(t, judge, orders.Move("smy", "eas"), cla.ErrIllegalDestination)
  // Unknown source
  assertOrderValidity(t, judge, orders.Move("a", "mid"), cla.ErrInvalidSource)
  // Unknown destination
  assertOrderValidity(t, judge, orders.Move("bre", "a"), cla.ErrInvalidDestination)
  // Missing sea path
  assertOrderValidity(t, judge, orders.Move("par", "mos"), cla.ErrMissingConvoyPath)
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

func TestMoveAdjudication(t *testing.T) {
  assertMove(t, Start(), "bre", "mid", true)
  assertMove(t, Start(), "stp/sc", "bot", true)
  assertMove(t, Start(), "vie", "bud", false)
  assertMove(t, Start(), "mid", "nat", false)
}
