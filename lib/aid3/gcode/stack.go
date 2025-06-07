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
	top *Node
}

func (s *Stk) Push(t *Tok) {
	n := &Node{
		t:    t,
		next: s.top,
	}
	s.top = n
}

func (s *Stk) Pop() *Tok {
	n := s.top
	if n == nil {
		return nil
	}
	s.top = n.next
	return n.t
}
