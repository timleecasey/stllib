package gcode

import (
	"bufio"
	"fmt"
	"github.com/timleecasey/stllib/lib/aid3/sim/reality"
	"log"
	"os"
	"strings"
	"unicode"
)

const (
	TOK_WS = iota
	TOK_UNKN
	TOK_NUMBER
	TOK_COMMENT
	TOK_PERCENT_SCOPE
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

var debugTokenize = false
var debugGcode = true

type Cmd struct {
	t      *Tok
	sibs   *Tok
	action func(a *reality.Affine)
}

type Settings struct {
}

type ParseTree struct {
	settings *Settings
	nodes    *NodeList
	stk      *Stk
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
	}

	if err := Tokenize(tree, srcFileNm); err != nil {
		return nil, err
	}
	log.Printf("Found %v tokens\n", tree.nodes.size)

	if err := MakeGcodeCommands(tree); err != nil {
		return nil, err
	}

	return tree, nil
}

func MakeGcodeCommands(t *ParseTree) error {
	nl := t.nodes
	if err := nl.Traverse(HandleToken); err != nil {
		return err
	}
	return nil
}

func HandleToken(n *Node) error {
	t := n.t
	//
	// As an example, G* may expect some amount of codes to follow, but not another G*
	//
	// If at G*, or other main type (M*), go forward to find the next non-arg token
	// then take the token as a command and the rest as possbily empty siblings
	// Also define an affine for the various pieces
	// Build a slot grammar, run the affine transform as arguments from the slots.
	//
	switch t.tokType {
	case TOK_N:
		// This is the Nth part of the line.
		break
	case TOK_M:
		switch t.src {
		case "M0", "M00": // Program stop
		case "M1", "M01": // Optional program stop
		case "M2", "M02": // end of program
		case "M3", "M03": // Spindle on clockwise
		case "M4", "M04": // Spindle on counterclockwise
		case "M5", "M05": // Spindle off
		case "M6", "M06": // Tool change
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
		switch t.src {
		case "G17":
		case "G18":
		case "G21":
		case "G00", "G0": // Rapid Positioning of Machine Tool
		case "G01", "G1": // Linear Interpolation
		case "G02", "G2": // Clockwise Arc Interpolation
		case "G03", "G3": // Counter-clockwise Interpolation
		case "G90": // Use absolute coordinates
		case "G08", "G8": // Increment Speed
		case "G09", "G9": // Decrement Speed
		//
		// Speed
		//
		case "G93": // Linear Feed Units
		case "G94": // Linear Feed Units
		case "G95": // Linear Feed Units
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
		case "G40", "G41", "G42", "G43", "G44": // Tool Offset Values
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
	case TOK_F:
	case TOK_H:
	case TOK_I:
	case TOK_J:
	case TOK_K:
	case TOK_X:
	case TOK_Y:
	case TOK_Z:
		break
	case TOK_T:
		switch t.src {
		case "T1", "T01":
			break
		case "T2", "T02":
			break
		default:
			return genErr(fmt.Sprintf("Unknown T type %v @ %v", t.src, t.lnPos))
		}
		break
	case TOK_S:
		break
	default:
		return genErr(fmt.Sprintf("Unknown token type %v @ %v", t.src, t.lnPos))
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
	comment := false
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
		//commentCh := "-"
		//if comment {
		//	commentCh = "T"
		//}
		//rStr := string(r)
		//log.Printf("CLS %v rune '%v' #%v", cls, rStr, commentCh)
		if comment || lineComment {
			cur[curI] = r
			curI++
			if r == ')' {
				comment = false
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
				//runeSl := cur[0:curI]
				//tokStr := string(runeSl)
				//log.Printf("TOK ln %v, %v: %v r# %v", lnMarker, position, string(cur[0:curI]), curI)
				//tokType := tokenType(tokStr)
				//t := &Tok{
				//	src:     tokStr,
				//	tokType: tokType,
				//	lnPos:   lnMarker,
				//	stPos:   stPos,
				//}
				//nl.Add(t)
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
				comment = true
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
