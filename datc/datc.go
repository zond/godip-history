package datc

import (
  "bufio"
  "fmt"
  "github.com/zond/godip/common"
  "io"
)

type Parser struct {
  Variant string
  State   common.State
}

func (self Parser) Parse(r io.Reader) {
  lr := bufio.NewReader(r)
  var err error
  var line string
  for line, err := lr.ReadString('\n'); err == nil; line, err = lr.ReadString('\n') {
    fmt.Println(line)
  }
}
