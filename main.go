package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tzaffi/go-tac-toe/game"
)

func parseTrimmed(input string) (int, int, error) {
	parts := strings.Split(input, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("bad input: %s", input)
	}
	row, err := strconv.ParseUint(parts[0], 10, 2)
	if err != nil {
		return 0, 0, fmt.Errorf("bad input: %s", input)
	}
	if row > 2 {
		return 0, 0, fmt.Errorf("bad input: %s", input)
	}
	col, err := strconv.ParseUint(parts[1], 10, 2)
	if err != nil {
		return 0, 0, fmt.Errorf("bad input: %s", input)
	}
	if col > 2 {
		return 0, 0, fmt.Errorf("bad input: %s", input)
	}
	return int(row), int(col), nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to Tic Tac Toe!")
	fmt.Println("Type 'exit' to quit the game.")

	board := game.NewBoard()

game:
	for {
		fmt.Println(board)
		if winner := board.Winner(); winner != "" {
			if winner == game.TIE {
				fmt.Println("Tied game.")
			} else {
				fmt.Printf("WE HAVE A WINNER: %s\n", winner)
			}

			fmt.Println("---GAME OVER---")
			break game
		}

		fmt.Print("> ")

		// Read the user input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "An error occurred while reading input. Please try again", err)
			continue
		}

		// Trim the input to handle new lines or spaces
		input = strings.TrimSpace(input)

		// Process the input
		switch input {
		case "AI":
			if err := board.AI(); err != nil {
				fmt.Printf("AI should NEVER ERROR but: %v\n", err)
			}
		case "exit":
			fmt.Println("Thank you for playing!")
			break game
		default:
			row, col, err := parseTrimmed(input)
			if err != nil {
				fmt.Printf("Try again; error: %s\n", err)
				continue
			}
			if err := board.Move(row, col); err != nil {
				fmt.Printf("Try again; error: %s\n", err)
			}
		}
	}
}
