package droidippy

import (
	"bufio"
	"fmt"
	"github.com/zond/godip/classical"
	cla "github.com/zond/godip/classical/common"
	dip "github.com/zond/godip/common"
	"github.com/zond/godip/state"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var phaseReg = regexp.MustCompile("^PHASE (\\d+) (\\S+) (\\S+)$")
var posReg = regexp.MustCompile("^(\\S+): (F|S|A) (\\S+)$")

const (
	positions = "POSITIONS"
	orders    = "ORDERS"
)

const (
	inNothing = iota
	inPositions
	inOrders
)

func verifyReversePositions(t *testing.T, s *state.State, scCollector map[dip.Province]dip.Nation, unitCollector map[dip.Province]dip.Unit) {
	for prov, nation1 := range s.SupplyCenters() {
		if nation2, ok := scCollector[prov]; !ok || nation1 != nation2 {
			t.Errorf("Found %v in %v, expected %v, %v", nation1, prov, nation2, ok)
		}
	}
	for prov, unit1 := range s.Units() {
		if unit2, ok := unitCollector[prov]; !ok || unit2.Nation != unit1.Nation || unit1.Type != unit2.Type {
			t.Errorf("Found %v in %v, expected %v, %v", unit1, prov, unit2, ok)
		}
	}
}

func verifyPosition(t *testing.T, s *state.State, match []string, scCollector map[dip.Province]dip.Nation, unitCollector map[dip.Province]dip.Unit) {
	if match[2] == "S" {
		if nation, _, ok := s.SupplyCenter(dip.Province(match[3])); ok && nation == dip.Nation(match[1]) {
			scCollector[dip.Province(match[3])] = nation
		} else {
			t.Errorf("Expected %v to own SC in %v, but found %v, %v", match[1], match[3], nation, ok)
		}
	} else if match[2] == "A" {
		if unit, _, ok := s.Unit(dip.Province(match[3])); ok && unit.Nation == dip.Nation(match[1]) && unit.Type == cla.Army {
			unitCollector[dip.Province(match[3])] = unit
		} else {
			t.Errorf("Expected to find %v %v in %v, but found %v, %v", match[1], cla.Army, match[3], unit, ok)
		}
	} else if match[2] == "F" {
		if unit, _, ok := s.Unit(dip.Province(match[3])); ok && unit.Nation == dip.Nation(match[1]) && unit.Type == cla.Fleet {
			unitCollector[dip.Province(match[3])] = unit
		} else {
			t.Errorf("Expected to find %v %v in %v, but found %v, %v", match[1], cla.Fleet, match[3], unit, ok)
		}
	} else {
		panic(fmt.Errorf("Unknown position description %v", match))
	}
}

func setPhase(s *state.State, match []string) {
	year, err := strconv.Atoi(match[1])
	if err != nil {
		panic(err)
	}
	season := match[2]
	typ := match[3]
	for (s.Phase().Year() <= year && (string(s.Phase().Season()) != season || string(s.Phase().Type()) != typ)) || s.Phase().Year() != year {
		s.Next()
	}
	if s.Phase().Year() > year {
		panic(fmt.Errorf("What the, we wanted %v but ended up with %v", match, s.Phase()))
	}
	fmt.Printf("In %v", s.Phase())
}

func assertGame(t *testing.T, name string) {
	file, err := os.Open(fmt.Sprintf("games/%v", name))
	if err != nil {
		panic(err)
	}
	s := classical.Start()
	lines := bufio.NewReader(file)
	var match []string
	state := inNothing
	scCollector := make(map[dip.Province]dip.Nation)
	unitCollector := make(map[dip.Province]dip.Unit)
	for line, err := lines.ReadString('\n'); err == nil; line, err = lines.ReadString('\n') {
		line = strings.TrimSpace(line)
		switch state {
		case inNothing:
			if match = phaseReg.FindStringSubmatch(line); match != nil {
				setPhase(s, match)
			} else if line == positions {
				state = inPositions
			} else {
				panic(fmt.Errorf("Unknown line for state inNothing: %v", line))
			}
		case inPositions:
			if match = posReg.FindStringSubmatch(line); match != nil {
				verifyPosition(t, s, match, scCollector, unitCollector)
			} else if line == orders {
				verifyReversePositions(t, s, scCollector, unitCollector)
				state = inOrders
			} else {
				panic(fmt.Errorf("Unknown line for state inPositions: %v", line))
			}
		case inOrders:
			panic("fixme")
		default:
			panic(fmt.Errorf("Unknown state %v", state))
		}
	}
}

func TestDroidippyGames(t *testing.T) {
	gamedir, err := os.Open("games")
	if err != nil {
		panic(err)
	}
	defer gamedir.Close()
	gamefiles, err := gamedir.Readdirnames(0)
	if err != nil {
		panic(err)
	}
	for _, name := range gamefiles {
		fmt.Println(name)
		assertGame(t, name)
	}
}
