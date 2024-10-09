package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type board struct {
	cells [][]rune
	rows  int
	cols  int
	alive int
}

type update_op struct {
	x        int
	y        int
	newValue rune
}

const Alive = '𐚰'
const Dead = '𐄙'

func IF[T any](cond bool, true_val, false_val T) T {
	if cond {
		return true_val
	}
	return false_val
}

func init_board(rows int, cols int) board {
	cells := make([][]rune, rows)

	for i := range cells {
		cells[i] = make([]rune, cols)

		for k := range cells[i] {
			cells[i][k] = Dead
		}
	}

	return board{cells, rows, cols, 0}
}

func print_board(board board) {
	for _, row := range board.cells {
		for _, cell := range row {
			fmt.Printf(IF(cell == Alive, "\033[33m%-2s\033[0m", "%-2s"), string(cell))
		}
		fmt.Println()
	}
}

func apply_conways_rule(board board) {
	changes := make([]update_op, 0)

	for i, row := range board.cells {
		for j, cell := range row {
			alive_neighbors := 0

			// TOP ROW
			if i > 0 {
				alive_neighbors += IF(board.cells[i-1][j] == Alive, 1, 0)
				if j > 0 {
					alive_neighbors += IF(board.cells[i-1][j-1] == Alive, 1, 0)
				}
				if j < board.cols-1 {
					alive_neighbors += IF(board.cells[i-1][j+1] == Alive, 1, 0)
				}
			}

			//BOTTOM ROW
			if i < board.rows-1 {
				alive_neighbors += IF(board.cells[i+1][j] == Alive, 1, 0)
				if j > 0 {
					alive_neighbors += IF(board.cells[i+1][j-1] == Alive, 1, 0)
				}
				if j < board.cols-1 {
					alive_neighbors += IF(board.cells[i+1][j+1] == Alive, 1, 0)
				}
			}

			// SIDE CELLS
			if j > 0 {
				alive_neighbors += IF(board.cells[i][j-1] == Alive, 1, 0)
			}
			if j < board.cols-1 {
				alive_neighbors += IF(board.cells[i][j+1] == Alive, 1, 0)
			}

			// Death condition or Living condition
			if cell == Alive && (alive_neighbors < 2 || alive_neighbors > 3) {
				changes = append(changes, update_op{i, j, Dead})
			} else if cell == Dead && alive_neighbors == 3 {
				changes = append(changes, update_op{i, j, Alive})
			}
		}
	}

	// Apply changes
	for _, op := range changes {
		board.cells[op.x][op.y] = op.newValue
		board.alive += IF(op.newValue == Alive, 1, -1)
	}

}

func seed_board(board board) {
	board.cells[10][25] = Alive
	board.cells[10][24] = Alive
	board.cells[9][25] = Alive
	board.cells[11][25] = Alive
	board.cells[11][26] = Alive

	board.alive = 5
}

func (board board) Init() tea.Cmd {
	return nil
}

func (board board) Update(event tea.Msg) (board, tea.Cmd) {

	switch event := event.(type) {
	case tea.KeyMsg:
		switch event.String() {
		case "q":
			return board, tea.Quit
		}
	}

	return board, nil

}
func main() {
	const SPEED_FACTOR = 200
	const GENERATION = 1_0000_000

	new_board := init_board(40, 80)
	seed_board(new_board)

	for i := 0; i < GENERATION; i++ {
		// clear screen
		fmt.Print("\033[H\033[2J")

		// print stats
		fmt.Printf("Speed(ms): %d\n", SPEED_FACTOR)
		fmt.Printf("Generation: %d\n", i)
		fmt.Printf("Alive: %d\n", new_board.alive)

		print_board(new_board)
		apply_conways_rule(new_board)
		time.Sleep(time.Millisecond * SPEED_FACTOR)

	}
}
