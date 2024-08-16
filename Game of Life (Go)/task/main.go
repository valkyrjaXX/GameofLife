package main

import (
	"fmt"
	"math/rand"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	var size int
	fmt.Scan(&size)

	u := NewRandomUniverse(size)

	var generation int
	nextUniverse := u

	for {
		fmt.Println(fmt.Sprintf("Generation #%d", generation+1))
		fmt.Println(fmt.Sprintf("Alive:%d", nextUniverse.Alive()))
		nextUniverse.Print()

		nextUniverse = nextUniverse.NextGeneration()
		generation++

		time.Sleep(500 * time.Millisecond)
		clearConsole()
	}
}

func clearConsole() {
	fmt.Print("\033[H\033[2J")
}

type Universe struct {
	matrix [][]bool
}

func (u *Universe) Generation(generation int) *Universe {
	if generation == 0 {
		return u
	}

	localUniverse := u
	for g := 0; g < generation; g++ {
		localUniverse = localUniverse.NextGeneration()
	}

	return localUniverse
}

func (u *Universe) Alive() int {
	var alive int
	for row := 0; row < len(u.matrix); row++ {
		for cell := 0; cell < len(u.matrix[row]); cell++ {
			if u.matrix[row][cell] {
				alive++
			}
		}
	}

	return alive
}

func (u *Universe) NextGeneration() *Universe {
	result := make([][]bool, u.Size())
	for row := 0; row < len(u.matrix); row++ {
		result[row] = make([]bool, u.Size())
		for cell := 0; cell < len(u.matrix[row]); cell++ {
			result[row][cell] = u.willSurvive(row, cell)
		}
	}

	return &Universe{result}
}

func (u *Universe) willSurvive(row, cell int) bool {
	alive := u.matrix[row][cell]
	neighbors := u.calculateNeighbors(row, cell)
	var aliveNeighbors int
	for _, neighbor := range neighbors {
		if u.matrix[neighbor[0]][neighbor[1]] {
			aliveNeighbors++
		}
	}

	if aliveNeighbors == 3 && !alive {
		return true
	}

	if !alive {
		return false
	}

	if aliveNeighbors == 2 || aliveNeighbors == 3 {
		return true
	}

	return false
}

func (u *Universe) Print() {
	for row := 0; row < len(u.matrix); row++ {
		for cell := 0; cell < len(u.matrix[row]); cell++ {
			if u.matrix[row][cell] {
				fmt.Print("O")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func (u *Universe) Size() int {
	return len(u.matrix)
}

func (u *Universe) calculateNeighbors(row, cell int) [][]int {
	size := u.Size()
	neighboursUpRow := mod(row-1, size)
	neighboursDownRow := mod(row+1, size)
	neighbourCellLeft := mod(cell-1, size)
	neighbourCellRight := mod(cell+1, size)

	return [][]int{
		{neighboursUpRow, neighbourCellLeft},
		{neighboursUpRow, cell},
		{neighboursUpRow, neighbourCellRight},
		{row, neighbourCellLeft},
		{row, neighbourCellRight},
		{neighboursDownRow, neighbourCellLeft},
		{neighboursDownRow, cell},
		{neighboursDownRow, neighbourCellRight},
	}
}

func NewRandomUniverse(size int) *Universe {
	matrix := make([][]bool, size)
	for i := 0; i < size; i++ {
		matrix[i] = make([]bool, size)
		for j := 0; j < size; j++ {
			matrix[i][j] = r.Intn(2) == 1
		}
	}

	return &Universe{matrix: matrix}
}

func mod(a, b int) int {
	return (a%b + b) % b
}
