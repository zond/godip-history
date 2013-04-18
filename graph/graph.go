package graph

import (
  "bytes"
  "fmt"
  . "github.com/zond/godip/common"
)

type Connectable interface {
  Prov(Province) Connectable
  Conn(Province) Connectable
  Flag(ProvinceFlag) Connectable
  SC(Nationality) Connectable
  Done() *Graph
}

func New() *Graph {
  return &Graph{
    nodes: make(map[Province]*node),
  }
}

type Graph struct {
  nodes map[Province]*node
}

func (self *Graph) String() string {
  buf := new(bytes.Buffer)
  for _, n := range self.nodes {
    fmt.Fprintf(buf, "%v", n)
  }
  return string(buf.Bytes())
}

func (self *Graph) Prov(n Province) *subNode {
  p, c := n.Split()
  if self.nodes[p] == nil {
    self.nodes[p] = &node{
      name:  p,
      subs:  make(map[Province]*subNode),
      graph: self,
    }
  }
  return self.nodes[p].sub(c)
}

type node struct {
  name  Province
  subs  map[Province]*subNode
  sc    *Nationality
  graph *Graph
}

func (self *node) String() string {
  buf := new(bytes.Buffer)
  fmt.Fprintf(buf, "%v", self.name)
  if self.sc != nil {
    fmt.Fprintf(buf, " %v", *self.sc)
  }
  fmt.Fprint(buf, "\n")
  for _, s := range self.subs {
    fmt.Fprintf(buf, "  %v\n", s)
  }
  return string(buf.Bytes())
}

func (self *node) sub(n Province) *subNode {
  if self.subs[n] == nil {
    self.subs[n] = &subNode{
      name:  n,
      edges: make(map[Province]*subNode),
      node:  self,
    }
  }
  return self.subs[n]
}

type subNode struct {
  name  Province
  edges map[Province]*subNode
  node  *node
  flags int
}

func (self *subNode) String() string {
  buf := new(bytes.Buffer)
  fmt.Fprintf(buf, "%v (%v) => ", self.name, self.flags)
  dests := make([]string, 0, len(self.edges))
  for n, _ := range self.edges {
    dests = append(dests, string(n))
  }
  fmt.Fprintf(buf, "%v", dests)
  return string(buf.Bytes())
}

func (self *subNode) getName() Province {
  return self.node.name.Join(self.name)
}

func (self *subNode) Conn(n Province) *subNode {
  target := self.node.graph.Prov(n)
  self.edges[target.getName()] = target
  return self
}

func (self *subNode) SC(n Nationality) *subNode {
  self.node.sc = &n
  return self
}

func (self *subNode) Flag(flags int) *subNode {
  self.flags |= flags
  return self
}

func (self *subNode) Prov(n Province) *subNode {
  return self.node.graph.Prov(n)
}

func (self *subNode) Done() *Graph {
  return self.node.graph
}
