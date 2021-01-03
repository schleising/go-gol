package main

import (
	"math/rand"
)

// BufferedBoard is the structure holding the two boards
type BufferedBoard struct {
	rows         int
	cols         int
	currentBoard []bool
	nextBoard    []bool
}

// BufferedBoardIF is the interface to the Board structure
type BufferedBoardIF interface {
	Iterate()
}

// Initialise creates a new random Board with the given dimensions
func Initialise(rows int, cols int) *BufferedBoard {
	var boards BufferedBoard
	boards.rows = rows
	boards.cols = cols
	boards.currentBoard = make([]bool, rows*cols)
	boards.nextBoard = make([]bool, rows*cols)

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if rand.Intn(2) == 1 {
				boards.setAlive(row, col)
			} else {
				boards.setDead(row, col)
			}
		}
	}

	return &boards
}

// Iterate calculates the state of the next frame
func (b *BufferedBoard) Iterate(ch chan bool) {
	for row := 0; row < b.rows; row++ {
		for col := 0; col < b.cols; col++ {
			b.calcNextState(row, col)
		}
	}

	b.swap()

	ch <- true
}

func (b *BufferedBoard) calcNextState(row, col int) {
	livingNeighbourCount := 0
	cellCurrentlyAlive := b.GetState(row, col)
	cellStillAlive := false

	for testRowDelta := -1; testRowDelta <= 1; testRowDelta++ {
		for testColDelta := -1; testColDelta <= 1; testColDelta++ {
			if testRowDelta != 0 && testColDelta != 0 {
				testRow := b.boundRows(row + testRowDelta)
				testCol := b.boundCols(col + testColDelta)

				if b.GetState(testRow, testCol) {
					livingNeighbourCount++
				}
			}
		}
	}

	if cellCurrentlyAlive {
		if livingNeighbourCount == 2 || livingNeighbourCount == 3 {
			cellStillAlive = true
		}
	} else {
		if livingNeighbourCount == 3 {
			cellStillAlive = true
		}
	}

	b.nextBoard[b.calculateIndex(row, col)] = cellStillAlive
}

func (b *BufferedBoard) boundRows(row int) int {
	for row < 0 {
		row += b.rows
	}

	for row >= b.rows {
		row -= b.rows
	}

	return row
}

func (b *BufferedBoard) boundCols(col int) int {
	for col < 0 {
		col += b.cols
	}

	for col >= b.cols {
		col -= b.cols
	}

	return col
}

func (b *BufferedBoard) calculateIndex(row, col int) int {
	index := row*b.cols + col

	return index
}

func (b *BufferedBoard) setAlive(row, col int) {
	b.currentBoard[b.calculateIndex(row, col)] = true
}

func (b *BufferedBoard) setDead(row, col int) {
	b.currentBoard[b.calculateIndex(row, col)] = false
}

// GetState returns the current dead or alive state as a bool
func (b *BufferedBoard) GetState(row, col int) bool {
	return b.currentBoard[b.calculateIndex(row, col)]
}

// Swap swaps the current and next boards
func (b *BufferedBoard) swap() {
	b.currentBoard, b.nextBoard = b.nextBoard, b.currentBoard
}
