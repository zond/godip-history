package graph

import (
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

type SubNode struct {
  name  string
  edges map[string]*SubNode
  node  *Node
  flags int
}

func (self SubNode) getName() string {
  return fmt.Sprintf("%v/%v", self.node.name, self.name)
}

func (self *SubNode) Conn(n, s string) *SubNode {
  target := self.node.graph.Node(n).Sub(s)
  self.edges[target.getName()] = target
  return self
}

func (self SubNode) Done() *Graph {
  return self.node.graph
}
