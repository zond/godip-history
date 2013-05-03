package droidippy

import (
	"bufio"
	"fmt"
	"github.com/zond/godip/classical"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var phaseReg = regexp.MustCompile("^PHASE (\\d+) (\\S+) (\\S+)$")

const (
	positions = "POSITIONS"
	orders    = "ORDERS"
)

func assertGame(t *testing.T, name string) {
	file, err := os.Open(fmt.Sprintf("games/%v", name))
	if err != nil {
		panic(err)
	}
	s := classical.Start()
	lines := bufio.NewReader(file)
	var match []string
	for line, err := lines.ReadString('\n'); err == nil; line, err = lines.ReadString('\n') {
		line = strings.TrimSpace(line)
		if match = phaseReg.FindStringSubmatch(line); match != nil {
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
				panic(fmt.Errorf("What the, we wanted %v but ended up with %v", line, s.Phase()))
			}
			fmt.Println(s.Phase())
		} else if line == positions {

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
