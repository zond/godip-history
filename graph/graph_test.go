package graph

import (
  "fmt"
  "testing"
)

func TestGraphBuilding(t *testing.T) {
  g := New().
    // spa
    Node("spa").Sub("").Conn("por", "").Conn("gas", "").Conn("mar", "").Done().
    // spa/sc
    Node("spa").Sub("sc").Conn("por", "").Conn("mid", "").Conn("gol", "").Conn("mar", "").Conn("wes", "").Done().
    // spa/nc
    Node("spa").Sub("nc").Conn("mid", "").Conn("mid", "").Conn("gas", "").Done().
    // por
    Node("por").Sub("").Conn("mid", "").Conn("spa", "nc").Conn("spa", "").Conn("spa", "sc").Done()
  fmt.Println(g)
}
