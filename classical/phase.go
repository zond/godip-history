package classical

import (
  . "github.com/zond/godip/classical/common"
  "github.com/zond/godip/common"
)

type phase struct {
  year   int
  season string
  typ    common.PhaseType
}

func (self *phase) PostProcess(s common.State) {
  if self.typ == Retreat {
  } else if self.typ == Build {
  } else if self.typ == Movement && self.season == Fall {
  }
}

func (self *phase) Year() int {
  return self.year
}

func (self *phase) Season() string {
  return self.season
}

func (self *phase) Type() common.PhaseType {
  return self.typ
}

func (self *phase) Prev() common.Phase {
  if self.typ == Retreat {
    return &phase{
      year:   self.year,
      season: self.season,
      typ:    Movement,
    }
  } else if self.typ == Movement {
    if self.season == Spring {
      if self.year == 1901 {
        return nil
      }
      return &phase{
        year:   self.year - 1,
        season: Winter,
        typ:    Build,
      }
    } else {
      return &phase{
        year:   self.year,
        season: Spring,
        typ:    Retreat,
      }
    }
  } else {
    return &phase{
      year:   self.year,
      season: Fall,
      typ:    Retreat,
    }
  }
  return nil
}

func (self *phase) Next() common.Phase {
  if self.typ == Movement {
    return &phase{
      year:   self.year,
      season: self.season,
      typ:    Retreat,
    }
  } else if self.typ == Retreat {
    if self.season == Spring {
      return &phase{
        year:   self.year,
        season: Fall,
        typ:    Movement,
      }
    } else {
      return &phase{
        year:   self.year,
        season: Winter,
        typ:    Build,
      }
    }
  } else {
    return &phase{
      year:   self.year + 1,
      season: Spring,
      typ:    Movement,
    }
  }
  return nil
}
