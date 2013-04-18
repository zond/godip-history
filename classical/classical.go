package classical

import (
  "fmt"
  . "github.com/zond/godip/classical/common"
  "github.com/zond/godip/classical/start"
  "github.com/zond/godip/common"
  "github.com/zond/godip/judge"
)

type phase struct {
  year   int
  season string
  typ    common.PhaseType
}

func (self phase) Year() int {
  return self.year
}

func (self phase) Season() string {
  return self.season
}

func (self phase) Type() common.PhaseType {
  return self.typ
}

func (self phase) Prev() (result common.Phase, err error) {
  if self.typ == Retreat {
    result = phase{
      year:   self.year,
      season: self.season,
      typ:    Movement,
    }
  } else if self.typ == Movement {
    if self.season == Spring {
      if self.year == 1901 {
        err = fmt.Errorf("No year earlier than 1901")
        return
      }
      result = phase{
        year:   self.year - 1,
        season: Winter,
        typ:    Build,
      }
    } else {
      result = phase{
        year:   self.year,
        season: Spring,
        typ:    Retreat,
      }
    }
  } else {
    result = phase{
      year:   self.year,
      season: Autumn,
      typ:    Retreat,
    }
  }
  return
}

func (self phase) Next() (result common.Phase, err error) {
  if self.typ == Movement {
    result = phase{
      year:   self.year,
      season: self.season,
      typ:    Retreat,
    }
  } else if self.typ == Retreat {
    if self.season == Spring {
      result = phase{
        year:   self.year,
        season: Autumn,
        typ:    Movement,
      }
    } else {
      result = phase{
        year:   self.year,
        season: Winter,
        typ:    Build,
      }
    }
  } else {
    result = phase{
      year:   self.year + 1,
      season: Spring,
      typ:    Movement,
    }
  }
  return
}

func Blank() *judge.State {
  return &judge.State{
    Graph:         start.Graph(),
    SupplyCenters: make(map[common.Province]common.Nationality),
    Units:         make(map[common.Province]common.Unit),
  }
}

func Start() (result *judge.State) {
  return &judge.State{
    Graph:         start.Graph(),
    SupplyCenters: start.SupplyCenters(),
    Units:         start.Units(),
    Phase: phase{
      1901,
      Spring,
      Movement,
    },
  }
}
