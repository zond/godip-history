package graph

import (
  "bytes"
  "fmt"
)

func New() *Graph {
  return &Graph{
    nodes: make(map[string]*Node),
  }
}

type Graph struct {
  nodes map[string]*Node
}

func (self *Graph) String() string {
  buf := new(bytes.Buffer)
  for _, n := range self.nodes {
    fmt.Fprintf(buf, "%v", n)
  }
  return string(buf.Bytes())
}

func (self *Graph) Node(n string) *Node {
  if self.nodes[n] == nil {
    self.nodes[n] = &Node{
      name:  n,
      subs:  make(map[string]*SubNode),
      graph: self,
    }
  }
  return self.nodes[n]
}

type Node struct {
  name  string
  subs  map[string]*SubNode
  graph *Graph
}

func (self *Node) String() string {
  buf := new(bytes.Buffer)
  fmt.Fprintf(buf, "%v\n", self.name)
  for _, s := range self.subs {
    fmt.Fprintf(buf, "  %v\n", s)
  }
  return string(buf.Bytes())
}

func (self *Node) Sub(s string) *SubNode {
  if self.subs[s] == nil {
    self.subs[s] = &SubNode{
      name:  s,
      edges: make(map[string]*SubNode),
      node:  self,
    }
  }
  return self.subs[s]
}

func (self *Node) Conn(n, s string) *SubNode {
  sub := self.Sub("")
  return sub.Conn(n, s)
}

func (self *Node) Con(n string) *SubNode {
  sub := self.Sub("")
  return sub.Con(n)
}

type SubNode struct {
  name  string
  edges map[string]*SubNode
  node  *Node
  flags int
}

func (self *SubNode) String() string {
  buf := new(bytes.Buffer)
  fmt.Fprintf(buf, "%v => ", self.name)
  dests := make([]string, 0, len(self.edges))
  for n, _ := range self.edges {
    dests = append(dests, n)
  }
  fmt.Fprintf(buf, "%v", dests)
  return string(buf.Bytes())
}

func (self *SubNode) getName() string {
  if self.name == "" {
    return fmt.Sprintf("%v", self.node.name)
  }
  return fmt.Sprintf("%v/%v", self.node.name, self.name)
}

func (self *SubNode) Conn(n, s string) *SubNode {
  target := self.node.graph.Node(n).Sub(s)
  self.edges[target.getName()] = target
  return self
}

func (self *SubNode) Con(n string) *SubNode {
  target := self.node.graph.Node(n).Sub("")
  self.edges[target.getName()] = target
  return self
}

func (self *SubNode) Done() *Graph {
  return self.node.graph
}

func (self *SubNode) Flag(flags int) *SubNode {
  self.flags |= flags
  return self
}
