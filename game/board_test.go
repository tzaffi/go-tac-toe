package game

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	require.Equal(t,
		newBoard([3][3]byte{{' ', ' ', ' '}, {' ', ' ', ' '}, {' ', ' ', ' '}}),
		NewBoard())
}

func TestPrint(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		board    [3][3]byte
		toString string
	}{
		{
			name:  "empty board",
			board: [3][3]byte{{' ', ' ', ' '}, {' ', ' ', ' '}, {' ', ' ', ' '}},
			toString: ` | | 
-----
 | | 
-----
 | | `},
		{
			name:  "full board",
			board: [3][3]byte{{'X', 'O', 'X'}, {'O', 'X', 'O'}, {'X', 'O', 'X'}},
			toString: `X|O|X
-----
O|X|O
-----
X|O|X`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			board := newBoard(tc.board)
			require.Equal(t, tc.toString, board.String())
			fmt.Println(board.String())
		})
	}
}

func TestMove(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name       string
		moves      [][]int
		errors     []bool
		finalState string
	}{
		{
			name:   "center move only",
			moves:  [][]int{{1, 1}},
			errors: []bool{false},
			finalState: ` | | 
-----
 |X| 
-----
 | | `,
		},
		{
			name:   "center move, upper right, lower right",
			moves:  [][]int{{1, 1}, {0, 2}, {2, 2}},
			errors: []bool{false, false, false},
			finalState: ` | |O
-----
 |X| 
-----
 | |X`,
		},
		{
			name:   "center move, upper right, upper right, lower right",
			moves:  [][]int{{1, 1}, {0, 2}, {0, 2}, {2, 2}},
			errors: []bool{false, false, true, false},
			finalState: ` | |O
-----
 |X| 
-----
 | |X`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			board := NewBoard()
			for i, move := range tc.moves {
				err := board.Move(move[0], move[1])
				if tc.errors[i] {
					require.NotNil(t, err)
				} else {
					require.Nil(t, err)
				}
			}
			require.Equal(t, tc.finalState, board.String())
		})
	}
}

func TestWinner(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		moves  [][]int
		winner string
	}{
		{
			name:   "center move only",
			moves:  [][]int{{1, 1}},
			winner: "",
		},
		{
			name:   "center move, upper right, lower right",
			moves:  [][]int{{1, 1}, {0, 2}, {2, 2}},
			winner: "",
		},
		{
			name:   "upper left to lower left diag X",
			moves:  [][]int{{1, 1}, {0, 2}, {0, 0}, {0, 1}, {2, 2}},
			winner: "X",
		},
		{
			name:   "lowest row O",
			moves:  [][]int{{1, 1}, {2, 0}, {0, 0}, {2, 2}, {0, 1}, {2, 1}},
			winner: "O",
		},
		{
			name:   "tie",
			moves:  [][]int{{0, 0}, {0, 1}, {1, 1}, {0, 2}, {2, 0}, {1, 0}, {1, 2}, {2, 2}, {2, 1}},
			winner: "TIE",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			board := NewBoard()
			for _, move := range tc.moves {
				err := board.Move(move[0], move[1])
				require.Nil(t, err)
			}
			fmt.Println(board)
			require.Equal(t, tc.winner, board.Winner())
		})
	}
}

func TestAI(t *testing.T) {
	t.Parallel()

	numTests := 100_000

	type stat struct {
		winner string
		moves  uint8
	}

	results := make(chan stat, numTests)
	wg := sync.WaitGroup{}
	wg.Add(numTests)

	for i := range numTests {
		t.Run(fmt.Sprintf("trial %d", i), func(t *testing.T) {
			defer wg.Done()

			board := NewBoard()
			for j := range 9 {
				err := board.AI()
				require.NoError(t, err, fmt.Sprintf("problem @ trial%d:%d", i, j))
				if winner := board.Winner(); winner != "" {
					results <- stat{winner, uint8(j) + 1}
					break
				}
			}
		})
	}

	wg.Wait()
	close(results)

	t.Run("printStats", func(t *testing.T) {
		winners := map[string]uint{}
		var totalMoves uint
		for result := range results {
			winners[result.winner]++
			totalMoves += uint(result.moves)
		}
		fmt.Printf("winners: %+v\navg. moves: %.2f\n", winners, float64(totalMoves)/float64(numTests))
		require.Equal(t, uint(numTests), winners[TIE]+winners[string(X)]+winners[string(O)])
	})
}
