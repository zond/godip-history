package datc

import (
  "bufio"
  "fmt"
  "github.com/zond/godip/common"
  "io"
  "regexp"
  "strings"
)

var clearCommentReg = regexp.MustCompile("(?m)^\\s*([^#\n\t]+?)\\s*(#.*)?$")

var variantReg = regexp.MustCompile("^VARIANT_ALL\\s+(\\S*)\\s*$")
var caseReg = regexp.MustCompile("^CASE\\s+(.*)$")

var prestateSetPhaseReg = regexp.MustCompile("^PRESTATE_SETPHASE\\s+(\\S+)\\s+(\\d+),\\s+(\\S+)\\s*$")

var stateReg = regexp.MustCompile("^([^:\\s]+):?\\s+(\\S+)\\s+(\\S+)\\s*$")

var ordersReg = regexp.MustCompile("^([^:]+):\\s+(.*)$")

const (
  prestate                   = "PRESTATE"
  orders                     = "ORDERS"
  poststateSame              = "POSTSTATE_SAME"
  end                        = "END"
  poststate                  = "POSTSTATE"
  poststateDislodged         = "POSTSTATE_DISLODGED"
  prestateSupplycenterOwners = "PRESTATE_SUPPLYCENTER_OWNERS"
)

func newState() *State {
  return &State{
    SCs:        make(map[common.Province]common.Nationality),
    Units:      make(map[common.Province]common.Unit),
    Dislodgeds: make(map[common.Province]common.Unit),
    Orders:     make(map[common.Province]common.Adjudicator),
  }
}

type State struct {
  SCs        map[common.Province]common.Nationality
  Units      map[common.Province]common.Unit
  Dislodgeds map[common.Province]common.Unit
  Orders     map[common.Province]common.Adjudicator
  Phase      common.Phase
}

func (self *State) copyFrom(o *State) {
  for prov, unit := range self.Units {
    o.Units[prov] = unit
  }
  for prov, dislodged := range self.Dislodgeds {
    o.Dislodgeds[prov] = dislodged
  }
  for prov, nation := range self.SCs {
    o.SCs[prov] = nation
  }
}

func newStatePair() *StatePair {
  return &StatePair{
    Before: newState(),
    After:  newState(),
  }
}

type StatePair struct {
  Case   string
  Before *State
  After  *State
}

func (self *StatePair) copyBeforeToAfter() {
  self.After.copyFrom(self.Before)
}

type StatePairHandler func(states *StatePair)

type OrderParser func(nation common.Nationality, text string) (province common.Province, order common.Adjudicator)

type PhaseParser func(season string, year int, typ string) common.Phase

type NationalityParser func(nation string) common.Nationality

type UnitTypeParser func(typ string) common.UnitType

type ProvinceParser func(prov string) common.Province

type Parser struct {
  Variant           string
  OrderParser       OrderParser
  PhaseParser       PhaseParser
  NationalityParser NationalityParser
  UnitTypeParser    UnitTypeParser
  ProvinceParser    ProvinceParser
}

const (
  waiting = iota
  inCase
  inPrestate
  inOrders
  inPoststate
  inPoststateDislodged
  inPrestateSupplycenterOwners
)

func (self Parser) Parse(r io.Reader, handler StatePairHandler) {
  lr := bufio.NewReader(r)
  var match []string
  state := waiting
  statePair := newStatePair()
  for line, err := lr.ReadString('\n'); err == nil; line, err = lr.ReadString('\n') {
    if match = clearCommentReg.FindStringSubmatch(line); match != nil {
      line = strings.TrimSpace(match[1])
      switch state {
      case waiting:
        if match = variantReg.FindStringSubmatch(line); match != nil {
          if match[1] != self.Variant {
            panic(fmt.Errorf("%+v only supports DATC files for %v variant", self, self.Variant))
          }
        } else if match = caseReg.FindStringSubmatch(line); match != nil {
          state = inCase
          statePair.Case = match[1]
        } else {
          panic(fmt.Errorf("Unrecognized line for state waiting: %#v", line))
        }
      case inPrestateSupplycenterOwners:
        if match = stateReg.FindStringSubmatch(line); match != nil {
          statePair.Before.SCs[self.ProvinceParser(match[3])] = self.NationalityParser(match[1])
        } else if line == prestate {
          state = inPrestate
        } else {
          panic(fmt.Errorf("Unrecognized line for state inPrestateSupplycenterOwners: %#v", line))
        }
      case inCase:
        if match = prestateSetPhaseReg.FindStringSubmatch(line); match != nil {
          statePair.Before.Phase = self.PhaseParser(match[1], common.MustParseInt(match[2]), match[3])
        } else if line == prestate {
          state = inPrestate
        } else if line == prestateSupplycenterOwners {
          state = inPrestateSupplycenterOwners
        } else {
          panic(fmt.Errorf("Unrecognized line for state inCase: %#v", line))
        }
      case inPrestate:
        if match = stateReg.FindStringSubmatch(line); match != nil {
          statePair.Before.Units[self.ProvinceParser(match[3])] = common.Unit{
            self.UnitTypeParser(match[2]),
            self.NationalityParser(match[1]),
          }
        } else if line == orders {
          state = inOrders
        } else {
          panic(fmt.Errorf("Unrecognized line for state inPrestate: %#v", line))
        }
      case inPoststate:
        if match = stateReg.FindStringSubmatch(line); match != nil {
          statePair.After.Units[self.ProvinceParser(match[3])] = common.Unit{
            self.UnitTypeParser(match[2]),
            self.NationalityParser(match[1]),
          }
        } else if line == end {
          handler(statePair)
          statePair = newStatePair()
          state = waiting
        } else if line == poststateDislodged {
          state = inPoststateDislodged
        } else {
          panic(fmt.Errorf("Unrecognized line for state inPoststate: %#v", line))
        }
      case inPoststateDislodged:
        if match = stateReg.FindStringSubmatch(line); match != nil {
          statePair.After.Dislodgeds[self.ProvinceParser(match[3])] = common.Unit{
            self.UnitTypeParser(match[2]),
            self.NationalityParser(match[1]),
          }
        } else if line == end {
          handler(statePair)
          statePair = newStatePair()
          state = waiting
        } else {
          panic(fmt.Errorf("Unrecognized line for state inPoststateDislodged: %#v", line))
        }
      case inOrders:
        if match = ordersReg.FindStringSubmatch(line); match != nil {
          prov, order := self.OrderParser(self.NationalityParser(match[1]), match[2])
          statePair.Before.Orders[prov] = order
        } else if line == poststateSame {
          statePair.copyBeforeToAfter()
        } else if line == poststate {
          state = inPoststate
        } else if line == end {
          handler(statePair)
          statePair = newStatePair()
          state = waiting
        } else {
          panic(fmt.Errorf("Unrecognized line for state inOrders: %#v", line))
        }
      }
    }
  }
}
