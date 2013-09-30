package orders

import (
	cla "github.com/zond/godip/classical/common"
	dip "github.com/zond/godip/common"
)

func Types() []dip.Order {
	return []dip.Order{
		&build{},
		&convoy{},
		&disband{},
		&hold{},
		&move{},
		&move{
			flags: map[dip.Flag]bool{
				cla.ViaConvoy: true,
			},
		},
		&support{},
	}
}
