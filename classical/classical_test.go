package classical

import (
  "fmt"
  cla "github.com/zond/godip/classical/common"
  "github.com/zond/godip/classical/orders"
  dip "github.com/zond/godip/common"
  "github.com/zond/godip/datc"
  "github.com/zond/godip/state"
  "os"
  "reflect"
  "testing"
)

func assertOrderValidity(t *testing.T, validator dip.Validator, order dip.Order, err error) {
  if e := order.Validate(validator); e != err {
    t.Errorf("%v should validate to %v, but got %v", order, err, e)
  }
}

func assertMove(t *testing.T, j *state.State, src, dst dip.Province, success bool) {
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
  assertOrderValidity(t, judge, orders.Move("bre", "wes"), cla.ErrIllegalConvoyUnit)
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
}

func TestMoveAdjudication(t *testing.T) {
  assertMove(t, Start(), "bre", "mid", true)
  assertMove(t, Start(), "stp/sc", "bot", true)
  assertMove(t, Start(), "vie", "bud", false)
  assertMove(t, Start(), "mid", "nat", false)
}

func testDATC(t *testing.T, statePair *datc.StatePair) {
  s := Blank(statePair.Before.Phase)
  for prov, unit := range statePair.Before.Units {
    s.SetUnit(prov, unit)
  }
  for prov, dislodged := range statePair.Before.Dislodgeds {
    s.SetDislodged(prov, dislodged)
  }
  for prov, nation := range statePair.Before.SCs {
    s.SetSC(prov, nation)
  }
  for prov, order := range statePair.Before.Orders {
    s.SetOrder(prov, order)
  }
  s.Next()
  for prov, unit := range statePair.After.Units {
    if found, ok := s.Units()[prov]; ok {
      if !found.Equal(unit) {
        t.Errorf("%v: Expected %v in %v, but found %v", statePair.Case, unit, prov, found)
      }
    } else {
      t.Errorf("%v: Expected %v in %v, but found nothing", statePair.Case, unit, prov)
    }
  }
}

func assertDATC(t *testing.T, file string) {
  in, err := os.Open(file)
  if err != nil {
    panic(err)
  }
  parser := datc.Parser{
    Variant:           "Standard",
    OrderParser:       DATCOrder,
    PhaseParser:       DATCPhase,
    NationalityParser: DATCNationality,
    UnitTypeParser:    DATCUnitType,
    ProvinceParser:    DATCProvince,
  }
  parser.Parse(in, func(statePair *datc.StatePair) {
    fmt.Printf("Running %v\n", statePair.Case)
    testDATC(t, statePair)
  })
}

func TestDATC(t *testing.T) {
  assertDATC(t, "datc/datc_v2.4_06.txt")
}
