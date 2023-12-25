package main

type Rectangle struct {
	Extent [2]int
}

type Node struct {
	Entire    Rectangle
	Left      *Node
	Right     *Node
	IsMatched bool
}

type Path []bool

func (n *Node) setLeft(rect Rectangle) {
	n.Left = &Node{rect, nil, nil, false}
}

func (n *Node) setRight(rect Rectangle) {
	n.Right = &Node{rect, nil, nil, false}
}

func (n *Node) cut(dirEq1 bool, length int) {
	extent := n.Entire.Extent
	if dirEq1 {
		n.setLeft(Rectangle{[2]int{extent[0], length}})
		n.setRight(Rectangle{[2]int{extent[0], extent[1] - length}})
	} else {
		n.setLeft(Rectangle{[2]int{length, extent[1]}})
		n.setRight(Rectangle{[2]int{extent[0] - length, extent[1]}})
	}
}
