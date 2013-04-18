package start

import (
  . "github.com/zond/godip/classical/common"
  "github.com/zond/godip/common"
)

func SupplyCenters() map[common.Province]common.Nationality {
  return map[common.Province]common.Nationality{
    "edi": England,
    "lvp": England,
    "lon": England,
    "bre": France,
    "par": France,
    "mar": France,
    "kie": Germany,
    "ber": Germany,
    "mun": Germany,
    "ven": Italy,
    "rom": Italy,
    "nap": Italy,
    "tri": Austria,
    "vie": Austria,
    "bud": Austria,
    "con": Turkey,
    "ank": Turkey,
    "smy": Turkey,
    "sev": Russia,
  }
}
