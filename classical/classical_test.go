package classical

import (
  "github.com/zond/godip/classical/common"
  "github.com/zond/godip/classical/orders"
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
  // Happy path
  assertOrderValidity(t, state, orders.Move("bre", "mid"), nil)
  // Fleet on land
  assertOrderValidity(t, state, orders.Move("bre", "par"), common.ErrIllegalDestination)
  // Army at sea
  assertOrderValidity(t, state, orders.Move("smy", "eas"), common.ErrIllegalDestination)
  // Unknown source
  assertOrderValidity(t, state, orders.Move("a", "mid"), common.ErrInvalidSource)
  // Unknown destination
  assertOrderValidity(t, state, orders.Move("bre", "a"), common.ErrInvalidDestination)
  // Too far
  assertOrderValidity(t, state, orders.Move("bre", "kie"), common.ErrIllegalDistance)
  // No unit
  assertOrderValidity(t, state, orders.Move("spa", "por"), common.ErrMissingUnit)
  // Bad phase
  state.Phase, _ = state.Phase.Next()
  assertOrderValidity(t, state, orders.Move("bre", "mid"), common.ErrInvalidPhase)
}
