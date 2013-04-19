package judge

import (
  . "github.com/zond/godip/common"
  "github.com/zond/godip/graph"
  "testing"
)

type testOrder int

func (self testOrder) Type() OrderType {
  return ""
}
func (self testOrder) Targets() []Province {
  return nil
}
func (self testOrder) Adjudicate(Resolver) (bool, error) {
  return false, nil
}
func (self testOrder) Validate(Validator) error {
  return nil
}
func (self testOrder) Execute(State) {
}

/*
     C
 A B 
     D
*/
func testGraph() Graph {
  return graph.New().
    Prov("a").Conn("b").Conn("b/sc").Conn("b/nc").
    Prov("b").Conn("a").Conn("c").Conn("d").
    Prov("b/sc").Conn("a").Conn("d").
    Prov("b/nc").Conn("a").Conn("c").
    Prov("b/ec").Conn("c").Conn("d").
    Prov("c").Conn("b/nc").Conn("b/ec").
    Prov("d").Conn("b/sc").Conn("b/ec").
    Done()
}

func assertOrderLocation(t *testing.T, j *Judge, prov Province, order Order, ok bool) {
  if o, k := j.Order(prov); o != order || k != ok {
    t.Errorf("Wrong order, wanted %v, %v at %v but got %v, %v", order, ok, prov, o, k)
  }
}

func TestJudgeLocations(t *testing.T) {
  j := New(testGraph(), nil, nil)
  j.SetOrders(map[Province]Order{
    "a":    testOrder(1),
    "b/ec": testOrder(2),
  })
  j.SetOrders(map[Province]Order{
    "b": testOrder(2),
  })
  assertOrderLocation(t, j, "a", nil, false)
  assertOrderLocation(t, j, "b", testOrder(2), true)
  assertOrderLocation(t, j, "b/sc", testOrder(2), true)
  assertOrderLocation(t, j, "b/ec", testOrder(2), true)
  assertOrderLocation(t, j, "b/nc", testOrder(2), true)
}
