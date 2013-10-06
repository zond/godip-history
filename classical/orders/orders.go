package orders

import (
	"fmt"
	cla "github.com/zond/godip/classical/common"
	dip "github.com/zond/godip/common"
	"time"
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

func Parse(bits []string) (result dip.Order, err error) {
	if len(bits) > 1 {
		switch dip.OrderType(bits[1]) {
		case (&build{}).DisplayType():
			if len(bits) == 3 {
				result = Build(dip.Province(bits[0]), dip.UnitType(bits[2]), time.Now())
			}
		case (&convoy{}).DisplayType():
			if len(bits) == 4 {
				result = Convoy(dip.Province(bits[0]), dip.Province(bits[2]), dip.Province(bits[3]))
			}
		case (&disband{}).DisplayType():
			if len(bits) == 2 {
				result = Disband(dip.Province(bits[0]), time.Now())
			}
		case (&hold{}).DisplayType():
			if len(bits) == 2 {
				result = Hold(dip.Province(bits[0]))
			}
		case (&move{}).DisplayType():
			if len(bits) == 3 {
				result = Move(dip.Province(bits[0]), dip.Province(bits[2]))
			}
		case (&move{flags: map[dip.Flag]bool{cla.ViaConvoy: true}}).DisplayType():
			if len(bits) == 3 {
				result = Move(dip.Province(bits[0]), dip.Province(bits[2])).ViaConvoy()
			}
		case (&support{}).DisplayType():
			if len(bits) == 4 {
				result = Support(dip.Province(bits[0]), dip.Province(bits[2]), dip.Province(bits[3]))
			}
		}
	}
	if result == nil {
		err = fmt.Errorf("Invalid order %+v", bits)
	}
	return
}
