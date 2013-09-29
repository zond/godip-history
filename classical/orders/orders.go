package orders

import (
	dip "github.com/zond/godip/common"
)

func Types() []dip.Order {
	return []dip.Order{
		&build{},
		&convoy{},
		&disband{},
		&hold{},
		&move{},
		&support{},
	}
}
