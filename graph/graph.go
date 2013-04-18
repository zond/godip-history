package graph

import (
  "bytes"
  "fmt"
  . "github.com/zond/godip/common"
)

type Connectable interface {
  Prov(Province) Connectable
  Conn(Province) Connectable
  Flag([]Flag) Connectable
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
  if sub, ok := self.subs[""]; ok {
    fmt.Fprintf(buf, " %v\n", sub)
  }
  for _, s := range self.subs {
    if s.name != "" {
      fmt.Fprintf(buf, "  %v\n", s)
    }
  }
  return string(buf.Bytes())
}

func (self *node) sub(n Province) *subNode {
  if self.subs[n] == nil {
    self.subs[n] = &subNode{
      name:  n,
      edges: make(map[Province]*subNode),
      node:  self,
      flags: make(map[Flag]bool),
    }
  }
  return self.subs[n]
}

type subNode struct {
  name  Province
  edges map[Province]*subNode
  node  *node
  flags map[Flag]bool
}

func (self *subNode) String() string {
  buf := new(bytes.Buffer)
  if self.name != "" {
    fmt.Fprintf(buf, "%v ", self.name)
  }
  flags := make([]Flag, 0, len(self.flags))
  for flag, _ := range self.flags {
    flags = append(flags, flag)
  }
  if len(flags) > 0 {
    fmt.Fprintf(buf, "%v ", flags)
  }
  dests := make([]string, 0, len(self.edges))
  for n, _ := range self.edges {
    dests = append(dests, string(n))
  }
  fmt.Fprintf(buf, "=> %v", dests)
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

func (self *subNode) Flag(flags ...Flag) *subNode {
  for _, flag := range flags {
    self.flags[flag] = true
  }
  return self
}

func (self *subNode) Prov(n Province) *subNode {
  return self.node.graph.Prov(n)
}

func (self *subNode) Done() *Graph {
  return self.node.graph
}
