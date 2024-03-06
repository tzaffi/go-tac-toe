package game

import (
	"fmt"
	"math/rand"
	"strings"
)

const EMPTY = ' '
const X = byte('X')
const O = byte('O')
const TIE = "TIE"

type Board struct {
	board  [3][3]byte
	oMoves bool
}

// NewBoard creates a new empty Board
func NewBoard() *Board {
	return newBoard([3][3]byte{
		{' ', ' ', ' '},
		{' ', ' ', ' '},
		{' ', ' ', ' '},
	})
}

func newBoard(board [3][3]byte) *Board {
	return &Board{board: board}
}

func (b *Board) String() string {
	sep := "-----"

	rowStr := func(row [3]byte) string {
		return fmt.Sprintf("%c|%c|%c", row[0], row[1], row[2])
	}

	return strings.Join([]string{
		rowStr(b.board[0]),
		sep,
		rowStr(b.board[1]),
		sep,
		rowStr(b.board[2]),
	},
		"\n")
}

func (b *Board) Move(row, col int) error {
	if row < 0 || col < 0 || row > 2 || col > 2 {
		return fmt.Errorf("bad input: (%d, %d)", row, col)
	}
	if b.board[row][col] != EMPTY {
		return fmt.Errorf("cannot overwrite mark at: (%d, %d)", row, col)
	}
	mark := X
	if b.oMoves {
		mark = O
	}
	b.board[row][col] = mark
	b.oMoves = !b.oMoves
	return nil
}

func (b *Board) AI() error {
	available := [][]int{}
	for i, row := range b.board {
		for j, val := range row {
			if val == EMPTY {
				available = append(available, []int{i, j})
			}
		}
	}
	n := len(available)
	if n == 0 {
		return fmt.Errorf("should never be calling AI when a tie")
	}
	randCoord := available[rand.Intn(n)]
	return b.Move(randCoord[0], randCoord[1])
}

func (b *Board) isTie() bool {
	for _, row := range b.board {
		for _, val := range row {
			if val == EMPTY {
				return false
			}
		}
	}
	return true
}

func (b *Board) Winner() string {
	if b.isTie() {
		return TIE
	}

	winner := func(group ...byte) byte {
		if group[0] == group[1] && group[1] == group[2] {
			return group[0]
		}
		return EMPTY
	}
	for _, row := range b.board {
		if w := winner(row[0], row[1], row[2]); w != EMPTY {
			return string(w)
		}
	}
	for i := range 3 {
		if w := winner(b.board[0][i], b.board[1][i], b.board[2][i]); w != EMPTY {
			return string(w)
		}
	}
	if w := winner(b.board[0][0], b.board[1][1], b.board[2][2]); w != EMPTY {
		return string(w)
	}
	if w := winner(b.board[0][2], b.board[1][1], b.board[2][0]); w != EMPTY {
		return string(w)
	}
	return ""
}
