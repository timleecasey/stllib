package gcode

type Tok struct {
	src     string
	tokType int
	lnPos   int
	stPos   int
}

const (
	MARKER = "_MARKER_"
)

type Node struct {
	t    *Tok
	next *Node
}

type Stk struct {
	top   *Node
	depth int
}

func (s *Stk) Push(t *Tok) {
	n := &Node{
		t:    t,
		next: s.top,
	}
	s.top = n
	s.depth++
}

func (s *Stk) Pop() *Tok {
	n := s.top
	if n == nil {
		return nil
	}
	s.top = n.next
	s.depth--
	return n.t
}

type NodeList struct {
	size int
	head *Node
	last *Node
}

func (l *NodeList) Add(t *Tok) {
	n := &Node{
		t:    t,
		next: nil,
	}

	if l.head == nil {
		l.head = n
		l.last = n
	} else {
		l.last.next = n
		l.last = n
	}

	l.size++

}

func (l *NodeList) Traverse(f func(n *Node) error) error {
	cur := l.head
	for cur != nil {
		//
		// Visit a node and possibly eat tokens
		//
		if err := f(cur); err != nil {
			return err
		}
		cur = cur.next
	}
	return nil
}
