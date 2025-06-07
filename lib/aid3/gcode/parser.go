package gcode

import (
	"bufio"
	"fmt"
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

type Cmd struct {
}

type Settings struct {
}

type ParseTree struct {
	settings *Settings
	cmds     []*Cmd
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
		cmds:     make([]*Cmd, 0),
	}

	lines, err := readLines(srcFileNm)
	if err != nil {
		log.Printf("Could not open %v: %v\n", srcFileNm, err)
	}

	for i := range lines {
		ln := lines[i]
		err := parseLine(tree, ln, i+1)
		if err != nil {
			return nil, err
		}
	}
	return tree, nil

}

func parseLine(tree *ParseTree, ln string, lnMarker int) error {

	ln = strings.ToUpper(ln)

	position := 0
	stPos := 0
	comment := false

	//
	// Used for comments
	//
	cur := make([]rune, 10000)
	curI := 0

	stk := Stk{}

	for _, r := range ln {
		position++
		cls := charClass(r)
		//commentCh := "-"
		//if comment {
		//	commentCh = "T"
		//}
		//rStr := string(r)
		//log.Printf("CLS %v rune '%v' #%v", cls, rStr, commentCh)
		if comment {
			cur[curI] = r
			curI++
			if r == ')' {
				comment = false
				buildTok(cur, curI, lnMarker, position, stPos, &stk)
				curI = 0
				stPos = 0
			}
			continue
		}
		switch cls {
		case CLS_WS:
			if curI > 0 {
				buildTok(cur, curI, lnMarker, position, stPos, &stk)
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
				//stk.Push(t)
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
			case '-':
				cur[curI] = r
				curI++
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
				//	stk.Push(t)
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
		buildTok(cur, curI, lnMarker, position, stPos, &stk)
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

func buildTok(cur []rune, curI int, lnMarker int, position int, stPos int, stk *Stk) {
	runeSl := cur[0:curI]
	tokStr := string(runeSl)
	log.Printf("TOK ln %v, %v: %v r# %v", lnMarker, position, string(cur[0:curI]), curI)
	tokType := tokenType(tokStr)
	t := &Tok{
		src:     tokStr,
		tokType: tokType,
		lnPos:   lnMarker,
		stPos:   stPos,
	}
	stk.Push(t)
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
