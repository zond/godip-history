package common

type Unit struct {
  Type        int
  Nationality string
}

type Order struct {
  Type    int
  Targets []string
}

type State interface {
  Resolve(orders []Order) State
}
