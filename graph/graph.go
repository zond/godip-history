package graph

import (
  "bytes"
  "fmt"
  "github.com/zond/godip/common"
)

type Connectable interface {
  Prov(common.Province) Connectable
  Conn(common.Province) Connectable
  Flag([]common.Flag) Connectable
  SC(common.Nationality) Connectable
  Done() *Graph
}

func New() *Graph {
  return &Graph{
    nodes: make(map[common.Province]*node),
  }
}

type Graph struct {
  nodes map[common.Province]*node
}

func (self *Graph) String() string {
  buf := new(bytes.Buffer)
  for _, n := range self.nodes {
    fmt.Fprintf(buf, "%v", n)
  }
  return string(buf.Bytes())
}

func (self *Graph) Has(n common.Province) (result bool) {
  p, c := n.Split()
  if node, ok := self.nodes[p]; ok {
    if _, ok := node.subs[c]; ok {
      result = true
    }
  }
  return
}

func (self *Graph) Flags(n common.Province) (result map[common.Flag]bool) {
  p, c := n.Split()
  if node, ok := self.nodes[p]; ok {
    if sub, ok := node.subs[c]; ok {
      result = sub.flags
    }
  }
  return
}

func (self *Graph) SC(n common.Province) (result *common.Nationality) {
  p, _ := n.Split()
  if node, ok := self.nodes[p]; ok {
    result = node.sc
  }
  return
}

func (self *Graph) edges(n common.Province) (result map[common.Province]*subNode) {
  p, c := n.Split()
  if node, ok := self.nodes[p]; ok {
    if sub, ok := node.subs[c]; ok {
      result = sub.edges
    }
  }
  return
}

type pathStep struct {
  path []common.Province
  pos  common.Province
}

func (self *Graph) pathHelper(dst common.Province, queue []pathStep, filter common.PathFilter, seen map[common.Province]bool) []common.Province {
  var newQueue []pathStep
  for _, step := range queue {
    seen[step.pos] = true
    for name, sub := range self.edges(step.pos) {
      if !seen[name] {
        if filter == nil || filter(name, sub.flags, sub.node.sc) {
          thisPath := append(append([]common.Province{}, step.path...), name)
          if name == dst {
            return thisPath
          }
          newQueue = append(newQueue, pathStep{
            path: thisPath,
            pos:  name,
          })
        }
      }
    }
  }
  if len(newQueue) > 0 {
    return self.pathHelper(dst, newQueue, filter, seen)
  }
  return nil
}

func (self *Graph) Path(src, dst common.Province, filter common.PathFilter) []common.Province {
  queue := []pathStep{
    pathStep{
      path: nil,
      pos:  src,
    },
  }
  return self.pathHelper(dst, queue, filter, make(map[common.Province]bool))
}

func (self *Graph) Coasts(prov common.Province) (result map[common.Province]bool) {
  result = make(map[common.Province]bool)
  p, _ := prov.Split()
  if node, ok := self.nodes[p]; ok {
    for _, sub := range node.subs {
      result[sub.getName()] = true
    }
  }
  return
}

func (self *Graph) Prov(n common.Province) *subNode {
  p, c := n.Split()
  if self.nodes[p] == nil {
    self.nodes[p] = &node{
      name:  p,
      subs:  make(map[common.Province]*subNode),
      graph: self,
    }
  }
  return self.nodes[p].sub(c)
}

type node struct {
  name  common.Province
  subs  map[common.Province]*subNode
  sc    *common.Nationality
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

func (self *node) sub(n common.Province) *subNode {
  if self.subs[n] == nil {
    self.subs[n] = &subNode{
      name:  n,
      edges: make(map[common.Province]*subNode),
      node:  self,
      flags: make(map[common.Flag]bool),
    }
  }
  return self.subs[n]
}

type subNode struct {
  name  common.Province
  edges map[common.Province]*subNode
  node  *node
  flags map[common.Flag]bool
}

func (self *subNode) String() string {
  buf := new(bytes.Buffer)
  if self.name != "" {
    fmt.Fprintf(buf, "%v ", self.name)
  }
  flags := make([]common.Flag, 0, len(self.flags))
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

func (self *subNode) getName() common.Province {
  return self.node.name.Join(self.name)
}

func (self *subNode) Conn(n common.Province) *subNode {
  target := self.node.graph.Prov(n)
  self.edges[target.getName()] = target
  return self
}

func (self *subNode) SC(n common.Nationality) *subNode {
  self.node.sc = &n
  return self
}

func (self *subNode) Flag(flags ...common.Flag) *subNode {
  for _, flag := range flags {
    self.flags[flag] = true
  }
  return self
}

func (self *subNode) Prov(n common.Province) *subNode {
  return self.node.graph.Prov(n)
}

func (self *subNode) Done() *Graph {
  return self.node.graph
}
