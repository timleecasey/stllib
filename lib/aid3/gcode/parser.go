package gcode

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const (
	TOK_WS = iota
	TOK_UNKN
	TOK_NUMBER
	TOK_COMMENT
	TOK_META
	TOK_PERCENT_SCOPE
	TOK_BREAK
	TOK_A
	TOK_B
	TOK_C
	TOK_F
	TOK_G
	TOK_H
	TOK_I
	TOK_J
	TOK_K
	TOK_M
	TOK_N
	TOK_O
	TOK_R
	TOK_S
	TOK_T
	TOK_X
	TOK_Y
	TOK_Z
)
const (
	CLS_WS = iota
	CLS_DIGIT
	CLS_LETTER
	CLS_PUNCT
	CLS_SYMBOL
	CLS_UNKN
)

const (
	CMD_UNKN = iota
	CMD_META
	CMD_ABSOLUTE
	CMD_FAST
	CMD_LINEAR
	CMD_CW_ARC
	CMD_CCW_ARC
	CMD_SPINDLE_OFF
	CMD_FEED_PER_MIN_MODE
	CMD_INVERSE_TIME_FEED
	CMD_FEED_PER_REVOLUTION
	CMD_SPINDLE_SPEED
	CMD_TOOL_CHANGE
	CMD_PLANE_XY
	CMD_PLANE_XZ
	CMD_PLANE_YZ
	CMD_INCH
	CMD_MM
)

var debugTokenize = false
var debugGcode = false

type Cmd struct {
	c      int
	t      *Tok
	sibs   *Tok
	coords *Coords
}

func (c *Cmd) CmdType() int {
	return c.c
}
func (c *Cmd) Src() string {
	return c.t.src
}
func (c *Cmd) Coords() *Coords {
	return c.coords
}

func (c *Cmd) String() string {
	if c == nil {
		return "nil"
	}

	// smells
	tstr := "--"
	if c.t != nil {
		tstr = c.t.src
	}

	coordStr := "--"
	if c.coords != nil {
		coordStr = fmt.Sprintf("X %v Y %v Z %v", c.coords.X, c.coords.Y, c.coords.Z)
	}
	return fmt.Sprintf("%v (%v) %v", tstr, c.c, coordStr)
}

type Settings struct {
	absoluteCoords bool
}

type Coords struct {
	X float64
	Y float64
	Z float64

	A float64
	B float64
	C float64

	F float64
	// G is missing, of course
	H float64

	I float64
	J float64
	K float64
	R float64
}

type ParseTree struct {
	settings *Settings
	nodes    *NodeList
	stk      *Stk
	cmds     *CmdList
	curCmd   *Cmd
}

func (t *ParseTree) TraverseCmds(f func(cn *CmdNode) error) error {
	if err := t.cmds.TraverseCmds(f); err != nil {
		return err
	}
	return nil
}

func (t *ParseTree) AddCmd(c *Cmd) {
	t.cmds.AddCmd(c)
}

type ParseError struct {
	msg string
}

func (pe *ParseError) Error() string {
	return pe.msg
}

func Parse(srcFileNm string) (*ParseTree, error) {

	tree := &ParseTree{
		settings: &Settings{},
		nodes:    &NodeList{},
		stk:      &Stk{},
		cmds:     &CmdList{},
		curCmd:   nil,
	}

	if err := Tokenize(tree, srcFileNm); err != nil {
		return nil, err
	}
	log.Printf("Found %v tokens\n", tree.nodes.size)

	// Maintain x/y/z and this is the coords
	// which are carried forward.
	tree.curCmd = &Cmd{
		coords: &Coords{
			X: 0,
			Y: 0,
			Z: 0,
		},
	}

	if err := MakeGcodeCommands(tree); err != nil {
		return nil, err
	}

	return tree, nil
}

func MakeGcodeCommands(t *ParseTree) error {
	nl := t.nodes
	if err := nl.Traverse(t, HandleToken); err != nil {
		return err
	}
	return nil
}

func HandleToken(tree *ParseTree, n *Node) error {
	t := n.t
	//
	// As an example, G* may expect some amount of codes to follow, including other G*
	// G1 G9 X Y Z F, or move to position with exact stop
	// G9 G1 X Y Z F, the same.
	//
	// If at G*, or other main type (M*), go forward to find the next non-arg token
	// then take the token as a command and the rest as possbily empty siblings
	// Also define an affine for the various pieces
	// Build a slot grammar, run the affine transform as arguments from the slots.
	//
	// Cmd has Coords as a set of slots
	//

	if debugGcode {
		log.Printf("Seeing %v\n", t.src)
	}
	switch t.tokType {
	case TOK_N:
		// This is the Nth part of the line.
		break

	case TOK_BREAK:
		prevType := CMD_UNKN
		var refTok *Tok
		refTok = nil
		if tree.curCmd != nil {
			prevType = tree.curCmd.c
			refTok = tree.curCmd.t
		}
		tree.curCmd = &Cmd{
			c:    prevType,
			t:    refTok,
			sibs: nil,
			coords: &Coords{
				X: tree.curCmd.coords.X,
				Y: tree.curCmd.coords.Y,
				Z: tree.curCmd.coords.Z,
			},
		}
		break

	case TOK_M:
		tree.curCmd = &Cmd{
			c:    CMD_UNKN,
			t:    t,
			sibs: nil,
			coords: &Coords{
				X: tree.curCmd.coords.X,
				Y: tree.curCmd.coords.Y,
				Z: tree.curCmd.coords.Z,
			},
		}
		switch t.src {
		case "M5", "M05": // Spindle off
			tree.curCmd.c = CMD_SPINDLE_OFF
			tree.AddCmd(tree.curCmd)
			break

		case "M0", "M00": // Program stop
		case "M1", "M01": // Optional program stop
		case "M2", "M02": // end of program
		case "M3", "M03": // Spindle on clockwise
		case "M4", "M04": // Spindle on counterclockwise
			break

		case "M6", "M06": // Manual tool change
			break

		case "M7", "M07": // Coolant on (mist)
		case "M8", "M08": // Coolant on
		case "M9", "M09": // Coolant off
		case "M10": // Clamp on
		case "M11": // Clamp off
		case "M19": // Spindle orientation
		case "M30": // Program end, return to start
		case "M40": // Spindle gear at middle
		case "M98": // Subprogram call
		case "M99": // Subprogram end
			break

		default:
			return genErr(fmt.Sprintf("Unknown M code %v @ %v", t.src, t.lnPos))
		}
	case TOK_G:
		tree.curCmd = &Cmd{
			c:    CMD_UNKN,
			t:    t,
			sibs: nil,
			coords: &Coords{
				X: tree.curCmd.coords.X,
				Y: tree.curCmd.coords.Y,
				Z: tree.curCmd.coords.Z,
			},
		}

		switch t.src {
		case "G17":
			tree.curCmd.c = CMD_PLANE_XY
			tree.AddCmd(tree.curCmd)
			break

		case "G18":
			tree.curCmd.c = CMD_PLANE_XZ
			tree.AddCmd(tree.curCmd)
			break
		case "G19":
			tree.curCmd.c = CMD_PLANE_YZ
			tree.AddCmd(tree.curCmd)
			break
		case "G20":
			tree.curCmd.c = CMD_INCH
			tree.AddCmd(tree.curCmd)
			break
		case "G21":
			tree.curCmd.c = CMD_MM
			tree.AddCmd(tree.curCmd)
			break
		case "G00", "G0": // Rapid Positioning of Machine Tool
			tree.curCmd.c = CMD_FAST
			tree.AddCmd(tree.curCmd)
			break

		case "G01", "G1": // Linear Interpolation
			tree.curCmd.c = CMD_LINEAR
			tree.AddCmd(tree.curCmd)
			break

		case "G02", "G2": // Clockwise Arc Interpolation
			tree.curCmd.c = CMD_CW_ARC
			tree.AddCmd(tree.curCmd)
			break

		case "G03", "G3": // Counter-clockwise Interpolation
			tree.curCmd.c = CMD_CCW_ARC
			tree.AddCmd(tree.curCmd)
			break

		case "G90": // Use absolute coordinates
			tree.curCmd.c = CMD_ABSOLUTE
			tree.AddCmd(tree.curCmd)
			break

		case "G08", "G8": // Increment Speed
		case "G09", "G9": // Decrement Speed (exact stop?)
			break
		//
		// Speed
		//
		case "G93": // Linear Feed Units
			tree.curCmd.c = CMD_INVERSE_TIME_FEED
			tree.AddCmd(tree.curCmd)
			break

		case "G94": // Linear Feed Units
			tree.curCmd.c = CMD_FEED_PER_MIN_MODE
			tree.AddCmd(tree.curCmd)
			break

		case "G95": // Linear Feed Units
			tree.curCmd.c = CMD_INVERSE_TIME_FEED
			tree.AddCmd(tree.curCmd)
			break

		case "G96": // Constant Surface Speed
		case "G97": // Constant Spindle Speed
		case "G61": // Exact Stop Mode
		case "G04": // Wait time
		//
		// Drilling
		//
		case "G81": // Simple drilling
		case "G82": // Simple drilling with dwell
		case "G83": // Deep hole drilling
		case "G84": // Tapping
		case "G40", "G41", "G42": // Tool Offset Values
			break

		case "G43": // Tool Offset Values
			tree.curCmd.c = CMD_INVERSE_TIME_FEED
			tree.AddCmd(tree.curCmd)
			break

		case "G44": // Tool Offset Values
			tree.curCmd.c = CMD_INVERSE_TIME_FEED
			tree.AddCmd(tree.curCmd)
			break

		case "G53", "G54", "G55", "G56", "G57", "G58", "G59": // Zero Offset Value
		case "G80", "G85", "G86", "G87", "G88", "G89": // Process Description
			break

		default:
			return genErr(fmt.Sprintf("Unknown G code %v @ %v", t.src, t.lnPos))
		}
	case TOK_O:
		break
	case TOK_COMMENT:
		break
	case TOK_META:
		tree.curCmd.c = CMD_META
		break
	case TOK_F:
		if f, err := strconv.ParseFloat(t.src[1:], 64); tree.curCmd == nil || err != nil {
			return genErr(fmt.Sprintf("Could not parse %v @ %v : %v", t.src, t.lnPos, err))
		} else {
			tree.curCmd.coords.F = f
		}
		break

	case TOK_H:
		if h, err := strconv.ParseFloat(t.src[1:], 64); tree.curCmd == nil || err != nil {
			return genErr(fmt.Sprintf("Could not parse %v @ %v : %v", t.src, t.lnPos, err))
		} else {
			tree.curCmd.coords.H = h
		}
		break

	case TOK_I:
		if i, err := strconv.ParseFloat(t.src[1:], 64); tree.curCmd == nil || err != nil {
			return genErr(fmt.Sprintf("Could not parse %v @ %v : %v", t.src, t.lnPos, err))
		} else {
			tree.curCmd.coords.I = i
		}
		break

	case TOK_J:
		if j, err := strconv.ParseFloat(t.src[1:], 64); tree.curCmd == nil || err != nil {
			return genErr(fmt.Sprintf("Could not parse %v @ %v : %v", t.src, t.lnPos, err))
		} else {
			tree.curCmd.coords.J = j
		}
		break

	case TOK_K:
		if k, err := strconv.ParseFloat(t.src[1:], 64); tree.curCmd == nil || err != nil {
			return genErr(fmt.Sprintf("Could not parse %v @ %v : %v", t.src, t.lnPos, err))
		} else {
			tree.curCmd.coords.K = k
		}
		break

	case TOK_A:
		if a, err := strconv.ParseFloat(t.src[1:], 64); tree.curCmd == nil || err != nil {
			return genErr(fmt.Sprintf("Could not parse %v @ %v : %v", t.src, t.lnPos, err))
		} else {
			tree.curCmd.coords.A = a
		}
		break

	case TOK_B:
		if b, err := strconv.ParseFloat(t.src[1:], 64); tree.curCmd == nil || err != nil {
			return genErr(fmt.Sprintf("Could not parse %v @ %v : %v", t.src, t.lnPos, err))
		} else {
			tree.curCmd.coords.B = b
		}
		break

	case TOK_C:
		if c, err := strconv.ParseFloat(t.src[1:], 64); tree.curCmd == nil || err != nil {
			return genErr(fmt.Sprintf("Could not parse %v @ %v : %v", t.src, t.lnPos, err))
		} else {
			tree.curCmd.coords.C = c
		}
		break

	case TOK_X:
		if x, err := strconv.ParseFloat(t.src[1:], 64); tree.curCmd == nil || err != nil {
			return genErr(fmt.Sprintf("Could not parse %v @ %v : %v", t.src, t.lnPos, err))
		} else {
			tree.curCmd.coords.X = x
		}
		break

	case TOK_Y:
		if y, err := strconv.ParseFloat(t.src[1:], 64); tree.curCmd == nil || err != nil {
			return genErr(fmt.Sprintf("Could not parse %v @ %v : %v", t.src, t.lnPos, err))
		} else {
			tree.curCmd.coords.Y = y
		}
		break

	case TOK_Z:
		if z, err := strconv.ParseFloat(t.src[1:], 64); tree.curCmd == nil || err != nil {
			return genErr(fmt.Sprintf("Could not parse %v @ %v : %v", t.src, t.lnPos, err))
		} else {
			tree.curCmd.coords.Z = z
		}
		break

	case TOK_R:
		if r, err := strconv.ParseFloat(t.src[1:], 64); tree.curCmd == nil || err != nil {
			return genErr(fmt.Sprintf("Could not parse %v @ %v : %v", t.src, t.lnPos, err))
		} else {
			tree.curCmd.coords.R = r
		}
		break

	case TOK_T:
		tree.curCmd.c = CMD_TOOL_CHANGE
		tree.AddCmd(tree.curCmd)
		break

	case TOK_S:
		tree.curCmd.c = CMD_SPINDLE_SPEED
		tree.AddCmd(tree.curCmd)
		break

	default:
		return genErr(fmt.Sprintf("Unknown token type %v @ %v", t.src, t.lnPos))
	}
	if debugGcode {
		log.Printf("State %v\n", tree.curCmd)
	}

	return nil
}

func Tokenize(t *ParseTree, srcFileNm string) error {
	lines, err := readLines(srcFileNm)
	if err != nil {
		log.Printf("Could not open %v: %v\n", srcFileNm, err)
	}

	for i := range lines {
		ln := lines[i]
		err := parseLine(t, ln, i+1)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseLine(tree *ParseTree, ln string, lnMarker int) error {

	ln = strings.ToUpper(ln)

	position := 0
	stPos := 0
	meta := false
	lineComment := false

	//
	// Used for comments
	//
	cur := make([]rune, 10000)
	curI := 0

	nl := tree.nodes

	for _, r := range ln {
		position++
		cls := charClass(r)
		if meta || lineComment {
			cur[curI] = r
			curI++
			if r == ')' {
				meta = false
				buildTok(cur, curI, lnMarker, position, stPos, nl)
				curI = 0
				stPos = 0
			}
			continue
		}
		switch cls {
		case CLS_WS:
			if curI > 0 {
				buildTok(cur, curI, lnMarker, position, stPos, nl)
				curI = 0
				stPos = 0
			}
			break

		case CLS_DIGIT, CLS_LETTER:
			if stPos == 0 {
				stPos = position
			}
			cur[curI] = r
			curI++
			break

		case CLS_PUNCT:
			switch r {
			case ';':
				lineComment = true
				break

			case '-':
				cur[curI] = r
				curI++
				break

			case '.':
				cur[curI] = r
				curI++
				break

			case '%':
				//if lnMarker > 1 {
				//	t := stk.Pop()
				//	if t == nil || t.tokType != TOK_PERCENT_SCOPE {
				//		return genErr(fmt.Sprintf("Miss-matched %v @ %v", string(r), lnMarker))
				//	}
				//} else {
				//	t := &Tok{
				//		src:     "%",
				//		tokType: TOK_PERCENT_SCOPE,
				//		lnPos:   lnMarker,
				//		stPos:   position,
				//	}
				//	nl.Add(t)
				//}
				break
			case '(':
				cur[curI] = r
				curI++
				meta = true
				break

			default:
				return genErr(fmt.Sprintf("Unknown PUNCT %v @ %v", string(r), lnMarker))
			}
			break

		case CLS_SYMBOL:
			return genErr(fmt.Sprintf("Unknown SYMBOL %v @ %v", r, lnMarker))
		case CLS_UNKN:
			return genErr(fmt.Sprintf("Unknown %v @ %v", r, lnMarker))
		}
	}
	if curI > 0 {
		buildTok(cur, curI, lnMarker, position, stPos, nl)
	}
	//
	t := &Tok{
		src:     "_NL_",
		tokType: TOK_BREAK,
		lnPos:   lnMarker,
		stPos:   stPos,
	}
	nl.Add(t)

	return nil
}

func readLines(fileNm string) ([]string, error) {
	file, err := os.Open(fileNm)
	if err != nil {
		return nil, fmt.Errorf("open file %q: %w", fileNm, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan file %q: %w", fileNm, err)
	}
	return lines, nil
}

func buildTok(cur []rune, curI int, lnMarker int, position int, stPos int, nl *NodeList) {
	runeSl := cur[0:curI]
	tokStr := string(runeSl)
	if debugTokenize {
		log.Printf("TOK ln %v, %v: %v r# %v", lnMarker, position, string(cur[0:curI]), curI)
	}
	tokType := tokenType(tokStr)
	t := &Tok{
		src:     tokStr,
		tokType: tokType,
		lnPos:   lnMarker,
		stPos:   stPos,
	}
	nl.Add(t)
}

func genErr(msg string) error {
	return &ParseError{msg}

}

func charClass(r rune) int {
	switch {
	case unicode.IsLetter(r):
		return CLS_LETTER
	case unicode.IsDigit(r):
		return CLS_DIGIT
	case unicode.IsSpace(r):
		return CLS_WS
	case unicode.IsPunct(r):
		return CLS_PUNCT
	case unicode.IsSymbol(r):
		return CLS_SYMBOL
	default:
		return CLS_UNKN
	}
}

func tokenType(tok string) int {
	switch tok[0] {
	case '(':
		return TOK_META
	case ';':
		return TOK_COMMENT
	case 'A':
		return TOK_A
	case 'B':
		return TOK_B
	case 'C':
		return TOK_C
	case 'F':
		return TOK_F
	case 'G':
		return TOK_G
	case 'H':
		return TOK_H
	case 'I':
		return TOK_I
	case 'J':
		return TOK_J
	case 'K':
		return TOK_K
	case 'M':
		return TOK_M
	case 'N':
		return TOK_N
	case 'O':
		return TOK_O
	case 'S':
		return TOK_S
	case 'T':
		return TOK_T
	case 'X':
		return TOK_X
	case 'Y':
		return TOK_Y
	case 'Z':
		return TOK_Z
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return TOK_NUMBER
	}
	return TOK_UNKN
}
