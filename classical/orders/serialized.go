package orders

import (
	dip "github.com/zond/godip/common"
	"time"
)

type serializedOrder struct {
	Targets []dip.Province
	At      time.Time
	Typ     dip.UnitType
	Flags   map[dip.Flag]bool
}
